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

	paths, err := validatePaths(os.Args)
	if err != nil {
		//handle error
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	configs, err := loadConfigFromFile(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	// Build out the new combined kubeconfig
	outkc := buildKubeconfig(configs)
	kubeconfig, err := convertToYAML(outkc)
	if err != nil {
		log.Fatalf("Unable to convert yaml: %s", err)
	}

	fmt.Println(kubeconfig)
}

func loadConfigFromFile(paths []string) ([]*api.Config, error) {
	var configs []*api.Config
	for _, path := range paths {
		// load both kubeconfigs into a api.Config struct
		objconfig, err := clientcmd.LoadFromFile(path)
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v\n", err)
		}
		configs = append(configs, objconfig)
	}
	return configs, nil
}

func validatePaths(args []string) ([]string, error) {
	// validate correct amount of Args
	if len(args) < 3 {
		log.Fatal("Please provide at least 2 kubeconfigs!")
	}

	// Capture the paths for Args
	var paths []string
	for _, c := range args[1:] {
		// test given paths
		if _, err := os.Stat(c); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file %s does not exist", c)
		}
		paths = append(paths, c)
	}
	return paths, nil
}

func buildKubeconfig(configs []*api.Config) *api.Config {
	clusters := make(map[string]*api.Cluster)
	authinfos := make(map[string]*api.AuthInfo)
	contexts := make(map[string]*api.Context)

	for _, config := range configs {
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
		CurrentContext: configs[0].CurrentContext,
	}
	return kubeconfig
}

func convertToYAML(kc *api.Config) (string, error) {
	// We need to encode it to json first to utilize the structs omitempty
	jsonData, err := json.Marshal(kc)
	if err != nil {
		return "", fmt.Errorf("JSON Marshal Error: %s", err)
	}

	// need to unmarshal json to convert to yaml
	var yamlData map[string]any
	if err := json.Unmarshal(jsonData, &yamlData); err != nil {
		return "", fmt.Errorf("JSON Unmarshal Error: %s", err)
	}

	// Convert to yaml
	kcyaml, err := yaml.Marshal(yamlData)
	if err != nil {
		return "", fmt.Errorf("YAML Marshal Error: %s", err)
	}
	return string(kcyaml), nil
}
