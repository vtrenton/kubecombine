package main

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd/api"
)

func TestBuildKubeconfig(t *testing.T) {

}

func TestConverttoYaml(t *testing.T) {
	t.Run("Bad JSON conversion of input kc", func(t *testing.T) {
		kubeconfig := &api.Config{
			APIVersion:     "v1",
			Kind:           "Config",
			Clusters:       "{}",
			AuthInfos:      "{}",
			Contexts:       "{}",
			CurrentContext: "{}",
		}
		got := convertToYAML()
	})
}
