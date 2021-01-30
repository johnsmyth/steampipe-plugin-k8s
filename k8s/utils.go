package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func GetNewClientset(ctx context.Context, connectionManager *connection.Manager) (*kubernetes.Clientset, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("GetNewClientset")

	// have we already created and cached the session?
	serviceCacheKey := "k8s"

	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset Found in Cache !!!!")
		return cachedData.(*kubernetes.Clientset), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	connectionManager.Cache.Set(serviceCacheKey, clientset)
	if _, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset Found in Cache after adding !!!!")
	} else {
		logger.Warn("!!!! Clientset NOT Found in Cache after adding !!!!")
	}
	return clientset, err
}

func v1TimeToRFC3339(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case v1.Time:
		return v.ToUnstructured(), nil
	case *v1.Time:
		if v == nil {
			return nil, nil
		}
		return v.ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("Invalid time format %T!\n", v)
	}
}
