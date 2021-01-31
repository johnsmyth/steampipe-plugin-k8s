package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableK8sReplicaSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8s_replicaset",
		Description: "Kubernetes ReplicaSet ensures that a specified number of pod replicas are running at any given time.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sReplicaSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sReplicaSets,
		},
		Columns: k8sCommonColumns([]*plugin.Column{}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sReplicaSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sReplicaSets")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	replicaSets, err := clientset.AppsV1().ReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range replicaSets.Items {
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}

func getK8sReplicaSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sReplicaSet")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	rs, err := clientset.AppsV1().ReplicaSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return rs, nil
}
