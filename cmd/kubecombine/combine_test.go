package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestValidatePaths(t *testing.T) {
	// Create temp files just for path validation
	tmpDir := t.TempDir()
	validFile1 := filepath.Join(tmpDir, "valid1.yaml")
	validFile2 := filepath.Join(tmpDir, "valid2.yaml")
	os.WriteFile(validFile1, []byte("test"), 0644)
	os.WriteFile(validFile2, []byte("test"), 0644)

	args := []string{"kubecombine", validFile1, validFile2}
	paths, err := validatePaths(args)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}

	if paths[0] != validFile1 || paths[1] != validFile2 {
		t.Fatalf("expected [%s, %s], got %v", validFile1, validFile2, paths)
	}
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

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(configs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(configs))
	}

	if configs[0].CurrentContext != "cluster1" {
		t.Fatalf("expected first config current context to be 'cluster1', got '%s'", configs[0].CurrentContext)
	}

	if configs[1].CurrentContext != "cluster2" {
		t.Fatalf("expected second config current context to be 'cluster2', got '%s'", configs[1].CurrentContext)
	}
}

func TestBuildKubeconfig(t *testing.T) {
	// Use your existing global testConfig1 and testConfig2
	configs := []*api.Config{testConfig1, testConfig2}

	result := buildKubeconfig(configs)

	// Verify combined result has both clusters
	if len(result.Clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(result.Clusters))
	}

	if _, exists := result.Clusters["cluster1"]; !exists {
		t.Fatal("expected cluster1 to exist in combined config")
	}

	if _, exists := result.Clusters["cluster2"]; !exists {
		t.Fatal("expected cluster2 to exist in combined config")
	}

	// Verify current context comes from first config
	if result.CurrentContext != "cluster1" {
		t.Fatalf("expected current context to be 'cluster1', got '%s'", result.CurrentContext)
	}
}

func TestConvertToYAML(t *testing.T) {
	result, err := convertToYAML(testConfig1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "cluster1") {
		t.Fatal("expected YAML output to contain 'cluster1'")
	}

	if !strings.Contains(result, "https://cluster1:6443") {
		t.Fatal("expected YAML output to contain 'https://cluster1:6443'")
	}
}
