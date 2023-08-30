package controllers

import (
	"context"
	mydomainv1 "github.com/costa92/webserver-operator/api/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
	"time"
)

func TestReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = mydomainv1.AddToScheme(scheme) // replace mydomainv1 with your actual package name
	_ = appsv1.AddToScheme(scheme)     // add Deployment to scheme
	_ = corev1.AddToScheme(scheme)     // add Service to scheme

	// Create a WebServer object
	webServer := &mydomainv1.WebServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-webserver",
			Namespace: "default",
		},
		Spec: mydomainv1.WebServerSpec{
			Replicas: 3,
			Image:    "nginx",
			Port:     80,
			NodePort: 30080,
		},
	}

	// Create a fake client with the WebServer object
	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(webServer).Build()

	// Create a ReconcileWebServerReconciler object with the scheme and fake client
	r := &WebServerReconciler{
		Client: client,
		Scheme: scheme,
	}

	// Mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-webserver",
			Namespace: "default",
		},
	}

	_, err := r.Reconcile(context.Background(), req)
	assert.NoError(t, err)
	// Check if deployment has been created and has the correct name
	dep := &appsv1.Deployment{}
	err = client.Get(context.Background(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}
	time.Sleep(200 * time.Second)

	// Check if service has been created and has the correct name
	svc := &corev1.Service{}
	err = client.Get(context.Background(), types.NamespacedName{Name: req.Name + "-service", Namespace: req.Namespace}, svc)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}
}

func TestDeploymentForWebserver(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = mydomainv1.AddToScheme(scheme) // replace mydomainv1 with your actual package name
	// Create a WebServer object
	webServer := &mydomainv1.WebServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-webserver",
			Namespace: "default",
		},
		Spec: mydomainv1.WebServerSpec{
			Replicas: 3,
			Image:    "nginx",
			Port:     80,
			NodePort: 30080,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(webServer).Build()

	// Create a ReconcileWebServerReconciler object with the scheme and fake client
	r := &WebServerReconciler{
		Client: client,
		Scheme: scheme,
	}
	dep := r.deploymentForWebserver(webServer)
	if dep.Name != "test-webserver" {
		t.Fatalf("deployment name: (%v)", dep.Name)
	}
}

func TestServiceForWebserver(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = mydomainv1.AddToScheme(scheme) // replace mydomainv1 with your actual package name
	// Create a WebServer object
	webServer := &mydomainv1.WebServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-webserver",
			Namespace: "default",
		},
		Spec: mydomainv1.WebServerSpec{
			Replicas: 3,
			Image:    "nginx",
			Port:     80,
			NodePort: 30080,
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(webServer).Build()

	// Create a ReconcileWebServerReconciler object with the scheme and fake client
	r := &WebServerReconciler{
		Client: client,
		Scheme: scheme,
	}
	svr := r.serviceForWebserver(webServer)
	if svr.Name != "test-webserver-service" {
		t.Fatalf("service name: (%v)", svr.Name)
	}
}
