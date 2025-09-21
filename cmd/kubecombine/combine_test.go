package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestValidatePaths(t *testing.T) {
	// Create temp files just for path validation
	tmpDir := t.TempDir()
	validFile1 := filepath.Join(tmpDir, "valid1.yaml")
	validFile2 := filepath.Join(tmpDir, "valid2.yaml")
	os.WriteFile(validFile1, []byte("test"), 0644)
	os.WriteFile(validFile2, []byte("test"), 0644)

	args := []string{"kubecombine", validFile1, validFile2}
	paths, err := validatePaths(args)

	assert.NoError(t, err)
	assert.Equal(t, []string{validFile1, validFile2}, paths)
}

func TestLoadConfigFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create actual valid kubeconfig YAML content
	config1YAML := `
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: b3JpZ2luYWxfZGF0YQ==
    server: https://cluster1:6443
  name: cluster1
contexts:
- context:
    cluster: cluster1
    user: user1
  name: cluster1
current-context: cluster1
users:
- name: user1
  user:
    client-certificate-data: b3JpZ2luYWxfZGF0YQ==
    client-key-data: b3JpZ2luYWxfZGF0YQ==
`

	config2YAML := `
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: bmV3X2NsdXN0ZXJfZGF0YQ==
    server: https://cluster2:6443
  name: cluster2
contexts:
- context:
    cluster: cluster2
    user: user2
  name: cluster2
current-context: cluster2
users:
- name: user2
  user:
    client-certificate-data: bmV3X2NsdXN0ZXJfZGF0YQ==
    client-key-data: bmV3X2NsdXN0ZXJfZGF0YQ==
`

	file1 := filepath.Join(tmpDir, "config1.yaml")
	file2 := filepath.Join(tmpDir, "config2.yaml")

	os.WriteFile(file1, []byte(config1YAML), 0644)
	os.WriteFile(file2, []byte(config2YAML), 0644)

	configs, err := loadConfigFromFile([]string{file1, file2})

	assert.NoError(t, err)
	assert.Len(t, configs, 2)
	assert.Equal(t, "cluster1", configs[0].CurrentContext)
	assert.Equal(t, "cluster2", configs[1].CurrentContext)
}

func TestBuildKubeconfig(t *testing.T) {
	// Use your existing global testConfig1 and testConfig2
	configs := []*api.Config{testConfig1, testConfig2}

	result := buildKubeconfig(configs)

	// Verify combined result has both clusters
	assert.Len(t, result.Clusters, 2)
	assert.Contains(t, result.Clusters, "cluster1")
	assert.Contains(t, result.Clusters, "cluster2")

	// Verify current context comes from first config
	assert.Equal(t, "cluster1", result.CurrentContext)
}

func TestConvertToYAML(t *testing.T) {
	result, err := convertToYAML(testConfig1)

	assert.NoError(t, err)
	assert.Contains(t, result, "cluster1")
	assert.Contains(t, result, "https://cluster1:6443")
}
