/*
Package k8s implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package k8s

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-k8s"

// Plugin creates this (k8s) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		// DefaultGetConfig: &plugin.GetConfig{
		// 	ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "NoSuchEntity"}),
		// },
		TableMap: map[string]*plugin.Table{
			"k8s_deployment": tableK8sDeployment(ctx),
			"k8s_pod":        tableK8sPod(ctx),
			"k8s_namespace":  tableK8sNamespace(ctx),
			"k8s_node":       tableK8sNode(ctx),
			"k8s_replicaset": tableK8sReplicaSet(ctx),
		},
	}

	return p
}
