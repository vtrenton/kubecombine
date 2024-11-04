package main

import (
	"fmt"
	"os"
	"errors"
	"log"
	//"encoding/base64"

	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func main() {
	// validate correct amount of Args
	if len(os.Args) != 3 {
		log.Fatal("Please provide 2 kubeconfigs!")
	}

	// test given paths
	if _, err := os.Stat(os.Args[1]); errors.Is(err, os.ErrNotExist) {
  		log.Fatalf("%s does not exist", os.Args[1])
	}

	if _, err := os.Stat(os.Args[2]); errors.Is(err, os.ErrNotExist) {
  		log.Fatalf("%s does not exist", os.Args[2])
	}

	// load both kubeconfigs into a api.Config struct
	config1, err := loadKubeconfig(os.Args[1])
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v\n", err)
	}
	
	config2, err := loadKubeconfig(os.Args[2])
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v\n", err)
	}
	
	outkc := buildKubeconfig(config1, config2)
	kubeconfig, err := convertToYAML(outkc)
	if err != nil {
		log.Fatalf("Unable to convert yaml: %s", err)
	}

	fmt.Println(kubeconfig)
}

// loadKubeconfig reads a Kubeconfig file and returns the API Config struct
func loadKubeconfig(path string) (*api.Config, error) {
	config, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func buildKubeconfig(config1, config2 *api.Config) *api.Config {
	// Create
	clusters := make(map[string]*api.Cluster)
	for name, cluster := range config1.Clusters {
		clusters[name] = &api.Cluster{
			Server:                   cluster.Server,
			CertificateAuthorityData: cluster.CertificateAuthorityData,
		}
	}
	for name, cluster := range config2.Clusters {
		clusters[name] = &api.Cluster{
			Server:                   cluster.Server,
			CertificateAuthorityData: cluster.CertificateAuthorityData,
		}
	}

	authinfos := make(map[string]*api.AuthInfo)
	for name, authinfo := range config1.AuthInfos {
		authinfos[name] = authinfo
	}
	for name, authinfo := range config2.AuthInfos {
		authinfos[name] = authinfo
	}

	contexts := make(map[string]*api.Context)
	for name, context := range config1.Contexts {
		contexts[name] = context	
	}
	for name, context := range config2.Contexts {
		contexts[name] = context
	}

	kubeconfig := &api.Config {
		Clusters: clusters,
		AuthInfos: authinfos,
		Contexts: contexts,
		CurrentContext: config1.CurrentContext,
	}
	return kubeconfig
}

func convertToYAML(kc *api.Config) (string, error) {
	kcyaml, err := yaml.Marshal(kc)
	if err != nil {
		return "", err
	}
	return string(kcyaml), nil
}
