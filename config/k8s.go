package baoConfig

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var openbaoNamespace string = "openbao"
var podPort int = 8200

type keySecret struct {
	Key        []string `json:"keys"`
	KeyEncoded []string `json:"keys_base64"`
}

// Get list of DNS names fro k8s pods
func (configInstance *MonitorConfig) MigratePodConfig(config *rest.Config) error {
	// create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// client for core
	coreClient := clientset.CoreV1()

	ctx := context.Background()

	// get pod list
	pods, err := coreClient.Pods(openbaoNamespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		return err
	}

	// clear existing DNS names
	configInstance.DNSnames = make(map[string]URL)

	// Use pod and its ip to fill in the "DNSNames" section
	r, _ := regexp.Compile(`stx-openbao-\d$`)
	for _, pod := range pods.Items {
		podName := pod.ObjectMeta.Name
		if r.Match([]byte(podName)) {
			podIP := pod.Status.PodIP
			podURL := fmt.Sprintf("%v.%v.pod.cluster.local", strings.ReplaceAll(podIP, ".", "-"), openbaoNamespace)
			configInstance.DNSnames[podName] = URL{podURL, podPort}
		}
	}

	// Validate input for DNSnames
	err = configInstance.validateDNS()
	if err != nil {
		return err
	}

	return nil
}

// Get root token and unseal key shards from k8s secrets
func (configInstance *MonitorConfig) MigrateSecretConfig(config *rest.Config) error {
	// create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// client for secret
	secretClient := clientset.CoreV1().Secrets(openbaoNamespace)

	ctx := context.Background()

	// get secrets list
	secrets, err := secretClient.List(ctx, metaV1.ListOptions{})
	if err != nil {
		return err
	}

	// Clear existing configs
	configInstance.Tokens = make(map[string]Token)
	configInstance.UnsealKeyShards = make(map[string]KeyShards)

	// Use secrets to fill in the "Tokens" and "UnsealKeyShards" section
	for _, secret := range secrets.Items {
		secretName := secret.ObjectMeta.Name
		if strings.HasPrefix(secretName, "cluster-key") {
			secretData := secret.Data["strdata"]
			if secretName == "cluster-key-root" {
				// secretData should be the root token
				configInstance.Tokens[secretName] = Token{Duration: 0, Key: string(secretData)}
			} else {
				// secretData should be an unseal key shard and its base 64 encoded version
				var newKey keySecret
				err := json.Unmarshal(secretData, &newKey)
				if err != nil {
					return err
				}
				configInstance.UnsealKeyShards[secretName] = KeyShards{
					Key:       newKey.Key[0],
					KeyBase64: newKey.KeyEncoded[0],
				}
			}
		}
	}

	// Validate input for Tokens
	err = configInstance.validateTokens()
	if err != nil {
		return err
	}

	// Validate input for unseal key shards
	err = configInstance.validateKeyShards()
	if err != nil {
		return err
	}

	return nil
}

// Get both configs
func (configInstance *MonitorConfig) MigrateK8sConfig(config *rest.Config) error {

	err := configInstance.MigratePodConfig(config)
	if err != nil {
		return err
	}

	err = configInstance.MigrateSecretConfig(config)
	if err != nil {
		return err
	}

	return nil
}
