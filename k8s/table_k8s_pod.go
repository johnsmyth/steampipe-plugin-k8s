package k8s

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableK8sPod(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8s_pod",
		Description: "Kubernetes Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sPod,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sPods,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the pod.  Name must be unique within a namespace."},
			{Name: "generate_name", Type: proto.ColumnType_STRING, Description: "GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "self_link", Type: proto.ColumnType_STRING, Description: "SelfLink is a URL representing this object."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object."},
			{Name: "resource_version", Type: proto.ColumnType_STRING, Description: "An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed."},
			{Name: "generation", Type: proto.ColumnType_INT, Description: "A sequence number representing a specific generation of the desired state."},

			{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "CreationTimestamp is a timestamp representing the server time when this object was created."},
			{Name: "deletion_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted."},

			{Name: "deletion_grace_period_seconds", Type: proto.ColumnType_INT, Description: "Number of seconds allowed for this object to gracefully terminate before it will be removed from the system.  Only set when deletionTimestamp is also set."},

			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
			{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."},

			{Name: "owner_references", Type: proto.ColumnType_JSON, Description: "List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller."},
			{Name: "finalizers", Type: proto.ColumnType_JSON, Description: "Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed."},
			{Name: "cluster_name", Type: proto.ColumnType_STRING, Description: "he name of the cluster which the object belongs to."},
			{Name: "managed_fields", Type: proto.ColumnType_JSON, Description: "ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field."},

			// hhmmmmm... why arent these working??
			{Name: "kind", Type: proto.ColumnType_STRING, Description: "Kind is a string value representing the REST resource this object represents."},
			{Name: "api_version", Type: proto.ColumnType_STRING, Description: "APIVersion defines the versioned schema of this representation of an object."},

			{Name: "spec", Type: proto.ColumnType_JSON, Description: "Specification of the desired behavior of the pod. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status."},
			{Name: "status", Type: proto.ColumnType_JSON, Description: "Most recently observed status of the pod. This data may not be up to date.  Populated by the system.  More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status."},
			//{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		},
	}
}

//// HYDRATE FUNCTIONS

func listK8sPods(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sPods")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		d.StreamListItem(ctx, pod)
	}

	return nil, nil
}

func getK8sPod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sPod")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return pod, nil
}

func isNotFoundError(err error) bool {
	if strings.HasSuffix(err.Error(), "not found") {
		return true
	}
	return false
}
