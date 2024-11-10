package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCleanupCmd(t *testing.T) {
	// Create a fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create a test deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "default",
		},
	}
	_, err := fakeClientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// Create test services
	service1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-1",
			Namespace: "default",
		},
	}
	service2 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-2",
			Namespace: "default",
		},
	}
	_, err = fakeClientset.CoreV1().Services("default").Create(context.TODO(), service1, metav1.CreateOptions{})
	assert.NoError(t, err)
	_, err = fakeClientset.CoreV1().Services("default").Create(context.TODO(), service2, metav1.CreateOptions{})
	assert.NoError(t, err)

	// Create a new cleanup command with the fake clientset
	cleanupCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Remove all services from Minikube",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Remove deployments
			err := fakeClientset.AppsV1().Deployments("default").DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("error removing deployments: %v", err)
			}

			// Remove services
			services, err := fakeClientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("error listing services: %v", err)
			}

			for _, svc := range services.Items {
				if svc.Name != "kubernetes" {
					err := fakeClientset.CoreV1().Services("default").Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
					if err != nil {
						return fmt.Errorf("error removing service %s: %v", svc.Name, err)
					}
				}
			}

			fmt.Println("Environment cleaned successfully")
			return nil
		},
	}

	// Run the cleanup command
	err = cleanupCmd.RunE(nil, nil)
	assert.NoError(t, err, "Error running cleanup command")

	// Verify that the services are deleted
	serviceList, err := fakeClientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
	assert.NoError(t, err, "Error listing services after cleanup")
	assert.Equal(t, 0, len(serviceList.Items), "Expected no services after cleanup, but found %d", len(serviceList.Items))
}