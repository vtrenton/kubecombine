package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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
	config1, err := clientcmd.LoadFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v\n", err)
	}

	config2, err := clientcmd.LoadFromFile(os.Args[2])
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v\n", err)
	}

	// Build out the new combined kubeconfig
	outkc := buildKubeconfig(config1, config2)
	kubeconfig, err := convertToYAML(outkc)
	if err != nil {
		log.Fatalf("Unable to convert yaml: %s", err)
	}

	fmt.Println(kubeconfig)
}

func buildKubeconfig(config1, config2 *api.Config) *api.Config {
	clusters := make(map[string]*api.Cluster)
	authinfos := make(map[string]*api.AuthInfo)
	contexts := make(map[string]*api.Context)

	for _, config := range []*api.Config{config1, config2} {
		for name, cluster := range config.Clusters {
			clusters[name] = cluster
		}

		for name, authinfo := range config.AuthInfos {
			authinfos[name] = authinfo
		}

		for name, context := range config.Contexts {
			contexts[name] = context
		}
	}

	kubeconfig := &api.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		Clusters:       clusters,
		AuthInfos:      authinfos,
		Contexts:       contexts,
		CurrentContext: config1.CurrentContext,
	}
	return kubeconfig
}

func convertToYAML(kc *api.Config) (string, error) {
	// We need to encode it to json first to utilize the structs omitempty
	jsonData, err := json.Marshal(kc)
	if err != nil {
		return "", err
	}

	// need to unmarshal json to convert to yaml
	var yamlData map[string]any
	if err := json.Unmarshal(jsonData, &yamlData); err != nil {
		return "", err
	}

	// Convert to yaml
	kcyaml, err := yaml.Marshal(yamlData)
	if err != nil {
		return "", err
	}
	return string(kcyaml), nil
}
