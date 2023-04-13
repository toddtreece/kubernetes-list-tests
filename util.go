package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

// generate a list of ids distributed across the total number of items
func generateIds(prefix string, total_items, id_count int) []string {
	ids := []string{}

	for i := 0; i < total_items; i += (total_items / id_count) {
		ids = append(ids, fmt.Sprintf("%s%d", prefix, i))
	}

	return ids
}

func getHostAndToken() (host string, token string, err error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	if p := os.Getenv("KUBECONFIG"); p != "" {
		kubeConfigPath = p
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return
	}

	if config.BearerToken == "" {
		return
	}

	token = config.BearerToken
	host = config.Host

	return
}
