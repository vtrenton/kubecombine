package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func main() {

	// validate correct amount of Args
	if len(os.Args) != 2 {
		fmt.Errorf("please specify 2 kubconfigs to combine")
	}

	// test given paths
	if _, err := os.Stat(os.Args[1]); errors.Is(err, os.ErrNotExist) {
  		fmt.Errorf("%s does not exist", os.Args[1])
	}

	if _, err := os.Stat(os.Args[2]); errors.Is(err, os.ErrNotExist) {
  		fmt.Errorf("%s does not exist", os.Args[2])
	}


	config1, err := loadKubeconfig(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return
	}
	
	config2, err := loadKubeconfig(os.Args[2])
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return
	}

	// Access some information from the Config struct
	fmt.Printf("Loaded Kubeconfig for Context: %s\n", config1.CurrentContext)
	for name, cluster := range config1.Clusters {
		fmt.Printf("Cluster Name: %s, Server: %s\n", name, cluster.Server)
	}

	// Access some information from the Config struct
	fmt.Printf("Loaded Kubeconfig for Context: %s\n", config2.CurrentContext)
	for name, cluster := range config2.Clusters {
		fmt.Printf("Cluster Name: %s, Server: %s\n", name, cluster.Server)
	}

}

// loadKubeconfig reads a Kubeconfig file and returns the API Config struct
func loadKubeconfig(path string) (*api.Config, error) {
	config, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	return config, nil
}

