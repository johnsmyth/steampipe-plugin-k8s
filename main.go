package main

import (
	"github.com/turbot/steampipe-plugin-k8s/k8s"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: k8s.Plugin})
}
