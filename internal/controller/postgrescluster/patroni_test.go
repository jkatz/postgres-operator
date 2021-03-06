// +build envtest

package postgrescluster

/*
 Copyright 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/crunchydata/postgres-operator/internal/naming"
	"github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1beta1"
	"go.opentelemetry.io/otel"
	"gotest.tools/v3/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestPatroniReplicationSecret(t *testing.T) {
	// Garbage collector cleans up test resources before the test completes
	if strings.EqualFold(os.Getenv("USE_EXISTING_CLUSTER"), "true") {
		t.Skip("USE_EXISTING_CLUSTER: Test fails due to garbage collection")
	}

	// setup the test environment and ensure a clean teardown
	tEnv, tClient, cfg := setupTestEnv(t, ControllerName)
	t.Cleanup(func() { teardownTestEnv(t, tEnv) })

	r := &Reconciler{}
	ctx, cancel := setupManager(t, cfg, func(mgr manager.Manager) {
		r = &Reconciler{
			Client:   tClient,
			Recorder: mgr.GetEventRecorderFor(ControllerName),
			Tracer:   otel.Tracer(ControllerName),
			Owner:    ControllerName,
		}
	})
	t.Cleanup(func() { teardownManager(cancel, t) })

	// test postgrescluster values
	var (
		clusterName = "hippocluster"
		namespace   = "postgres-operator-test-" + rand.String(6)
		clusterUID  = types.UID("hippouid")
	)

	ns := &corev1.Namespace{}
	ns.Name = namespace
	assert.NilError(t, tClient.Create(ctx, ns))
	t.Cleanup(func() { assert.Check(t, tClient.Delete(ctx, ns)) })

	// create a PostgresCluster to test with
	postgresCluster := &v1beta1.PostgresCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: namespace,
			UID:       clusterUID,
		},
	}

	rootCA, err := r.reconcileRootCertificate(ctx, postgresCluster)
	assert.NilError(t, err)

	t.Run("reconcile", func(t *testing.T) {
		_, err = r.reconcileReplicationSecret(ctx, postgresCluster, rootCA)
		assert.NilError(t, err)
	})

	t.Run("validate", func(t *testing.T) {

		patroniReplicationSecret := &corev1.Secret{ObjectMeta: naming.ReplicationClientCertSecret(postgresCluster)}
		patroniReplicationSecret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
		err = r.Client.Get(ctx, client.ObjectKeyFromObject(patroniReplicationSecret), patroniReplicationSecret)
		assert.NilError(t, err)

		t.Run("ca.crt", func(t *testing.T) {

			clientCert, ok := patroniReplicationSecret.Data["ca.crt"]
			assert.Assert(t, ok)

			assert.Assert(t, strings.HasPrefix(string(clientCert), "-----BEGIN CERTIFICATE-----"))
			assert.Assert(t, strings.HasSuffix(string(clientCert), "-----END CERTIFICATE-----\n"))
		})

		t.Run("tls.crt", func(t *testing.T) {

			clientCert, ok := patroniReplicationSecret.Data["tls.crt"]
			assert.Assert(t, ok)

			assert.Assert(t, strings.HasPrefix(string(clientCert), "-----BEGIN CERTIFICATE-----"))
			assert.Assert(t, strings.HasSuffix(string(clientCert), "-----END CERTIFICATE-----\n"))
		})

		t.Run("tls.key", func(t *testing.T) {

			clientKey, ok := patroniReplicationSecret.Data["tls.key"]
			assert.Assert(t, ok)

			assert.Assert(t, strings.HasPrefix(string(clientKey), "-----BEGIN EC PRIVATE KEY-----"))
			assert.Assert(t, strings.HasSuffix(string(clientKey), "-----END EC PRIVATE KEY-----\n"))
		})

	})

	t.Run("check replication certificate secret projection", func(t *testing.T) {
		// example auto-generated secret projection
		testSecretProjection := &v1.SecretProjection{
			LocalObjectReference: v1.LocalObjectReference{
				Name: naming.ReplicationClientCertSecret(postgresCluster).Name,
			},
			Items: []v1.KeyToPath{
				{
					Key:  naming.ReplicationCert,
					Path: naming.ReplicationCertPath,
				},
				{
					Key:  naming.ReplicationPrivateKey,
					Path: naming.ReplicationPrivateKeyPath,
				},
				{
					Key:  naming.ReplicationCACert,
					Path: naming.ReplicationCACertPath,
				},
			},
		}

		rootCA, err := r.reconcileRootCertificate(ctx, postgresCluster)
		assert.NilError(t, err)

		testReplicationSecret, err := r.reconcileReplicationSecret(ctx, postgresCluster, rootCA)
		assert.NilError(t, err)

		t.Run("check standard secret projection", func(t *testing.T) {
			secretCertProj := replicationCertSecretProjection(testReplicationSecret)

			assert.DeepEqual(t, testSecretProjection, secretCertProj)
		})
	})

}

func TestReconcilePatroniStatus(t *testing.T) {
	ctx := context.Background()

	tEnv, tClient, cfg := setupTestEnv(t, ControllerName)
	t.Cleanup(func() { teardownTestEnv(t, tEnv) })
	r := &Reconciler{}
	ctx, cancel := setupManager(t, cfg, func(mgr manager.Manager) {
		r = &Reconciler{
			Client:   mgr.GetClient(),
			Recorder: mgr.GetEventRecorderFor(ControllerName),
			Tracer:   otel.Tracer(ControllerName),
			Owner:    ControllerName,
		}
	})
	t.Cleanup(func() { teardownManager(cancel, t) })

	namespace := "test-reconcile-patroni-status"
	systemIdentifier := "6952526174828511264"
	createResources := func(index, readyReplicas int,
		writeAnnotation bool) (*v1beta1.PostgresCluster, *observedInstances) {

		i := strconv.Itoa(index)
		clusterName := "patroni-status-" + i
		instanceName := "test-instance-" + i
		instanceSet := "set-" + i

		labels := map[string]string{
			naming.LabelCluster:     clusterName,
			naming.LabelInstanceSet: instanceSet,
			naming.LabelInstance:    instanceName,
		}

		postgresCluster := &v1beta1.PostgresCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterName,
				Namespace: namespace,
			},
		}

		runner := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      instanceName,
				Labels:    labels,
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
				},
			},
		}

		endpoints := &corev1.Endpoints{
			ObjectMeta: naming.PatroniDistributedConfiguration(postgresCluster),
		}
		if writeAnnotation {
			endpoints.ObjectMeta.Annotations = make(map[string]string)
			endpoints.ObjectMeta.Annotations["initialize"] = systemIdentifier
		}
		assert.NilError(t, tClient.Create(ctx, endpoints, &client.CreateOptions{}))

		instance := &Instance{
			Name: instanceName, Runner: runner,
		}
		for i := 0; i < readyReplicas; i++ {
			instance.Pods = append(instance.Pods, &v1.Pod{
				Status: v1.PodStatus{
					Conditions: []v1.PodCondition{{
						Type:    corev1.PodReady,
						Status:  corev1.ConditionTrue,
						Reason:  "test",
						Message: "test",
					}},
				},
			})
		}
		observedInstances := &observedInstances{}
		observedInstances.forCluster = []*Instance{instance}

		return postgresCluster, observedInstances
	}

	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	assert.NilError(t, tClient.Create(ctx, ns))
	t.Cleanup(func() { assert.Check(t, tClient.Delete(ctx, ns)) })

	testsCases := []struct {
		requeueExpected bool
		readyReplicas   int
		writeAnnotation bool
	}{
		{requeueExpected: false, readyReplicas: 1, writeAnnotation: true},
		{requeueExpected: true, readyReplicas: 1, writeAnnotation: false},
		{requeueExpected: false, readyReplicas: 0, writeAnnotation: false},
		{requeueExpected: false, readyReplicas: 0, writeAnnotation: false},
	}

	for i, tc := range testsCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			postgresCluster, observedInstances := createResources(i, tc.readyReplicas,
				tc.writeAnnotation)
			result, err := r.reconcilePatroniStatus(ctx, postgresCluster, observedInstances)
			if tc.requeueExpected {
				assert.NilError(t, err)
				assert.Assert(t, result.RequeueAfter == 1*time.Second)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, reconcile.Result{})
			}
		})
	}
}
