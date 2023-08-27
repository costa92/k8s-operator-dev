package controllers

import (
	"context"
	mydomainv1 "github.com/costa92/webserver-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebServerReconciler_Reconcile(t *testing.T) {
	// Create a fake client and scheme
	scheme := runtime.NewScheme()
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	// Create a WebServer instance for testing
	webServer := &mydomainv1.WebServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-webserver",
			Namespace: "default",
		},
		Spec: mydomainv1.WebServerSpec{
			Replicas: 2,
			Image:    "nginx:latest",
			Port:     80,
			NodePort: 30080,
		},
	}

	// Create a reconciler instance
	r := &WebServerReconciler{
		Client: client,
		Scheme: scheme,
	}

	// Create a reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      webServer.Name,
			Namespace: webServer.Namespace,
		},
	}

	// Call the Reconcile function
	result, err := r.Reconcile(context.Background(), req)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the result indicates no requeue
	assert.False(t, result.Requeue)

	// You can add more assertions here to verify the behavior of the Reconcile function
}
