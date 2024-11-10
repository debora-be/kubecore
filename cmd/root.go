package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig string
	clientset  *kubernetes.Clientset
)

var rootCmd = &cobra.Command{
	Use:   "kubecore",
	Short: "kubecore - Kubernetes Service Manager",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("error building kubeconfig: %v", err)
		}

		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %v", err)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	if home := homedir.HomeDir(); home != "" {
		rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	}
}
