package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestIngressCmd(t *testing.T) {
	// Create a fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create a test service name
	serviceName := "test-service"

	// Set the hostname and path flags
	hostname = "example.com"
	path = "/test"

	// Run the createIngress function with the fake clientset
	err := createIngress(serviceName, fakeClientset)
	assert.NoError(t, err)

	// Get the created ingress from the fake clientset
	ingress, err := fakeClientset.NetworkingV1().Ingresses("default").Get(context.TODO(), serviceName+"-ingress", metav1.GetOptions{})
	assert.NoError(t, err)

	// Assert the ingress properties
	assert.Equal(t, serviceName+"-ingress", ingress.Name)
	assert.Equal(t, "nginx", *ingress.Spec.IngressClassName)
	assert.Equal(t, hostname, ingress.Spec.Rules[0].Host)
	assert.Equal(t, path, ingress.Spec.Rules[0].HTTP.Paths[0].Path)
	assert.Equal(t, serviceName, ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name)
	assert.Equal(t, int32(80), ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number)
}
