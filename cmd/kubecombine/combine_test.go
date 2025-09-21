package main

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd/api"
)

// define 2 global kubeconfig objects...
var (
	testConfig1 = &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			"cluster1": {
				Server:                   "https://cluster1:6443",
				CertificateAuthorityData: []byte("b3JpZ2luYWxfZGF0YQ=="),
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"user1": {
				ClientCertificateData: []byte("b3JpZ2luYWxfZGF0YQ=="),
				ClientKeyData:         []byte("b3JpZ2luYWxfZGF0YQ=="),
			},
		},
		Contexts: map[string]*api.Context{
			"cluster1": {
				Cluster:  "cluster1",
				AuthInfo: "user1",
			},
		},
		CurrentContext: "cluster1",
	}

	testConfig2 = &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			"cluster2": {
				Server:                   "https://cluster2:6443",
				CertificateAuthorityData: []byte("bmV3X2NsdXN0ZXJfZGF0YQ=="),
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"user2": {
				ClientCertificateData: []byte("bmV3X2NsdXN0ZXJfZGF0YQ=="),
				ClientKeyData:         []byte("bmV3X2NsdXN0ZXJfZGF0YQ=="),
			},
		},
		Contexts: map[string]*api.Context{
			"cluster2": {
				Cluster:  "cluster2",
				AuthInfo: "user2",
			},
		},
		CurrentContext: "cluster2",
	}
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
