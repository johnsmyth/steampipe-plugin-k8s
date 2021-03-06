package k8s

import (
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

var objectMetadataPrimaryColumns = []*plugin.Column{
	{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},

	{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the object.  Name must be unique within a namespace."},
	{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
	{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object."},
	{Name: "cluster_name", Type: proto.ColumnType_STRING, Description: "The name of the cluster which the object belongs to."},
}

var objectMetadataSecondaryColumns = []*plugin.Column{

	{Name: "generate_name", Type: proto.ColumnType_STRING, Description: "GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided."},
	{Name: "self_link", Type: proto.ColumnType_STRING, Description: "SelfLink is a URL representing this object."},
	{Name: "resource_version", Type: proto.ColumnType_STRING, Description: "An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed."},
	{Name: "generation", Type: proto.ColumnType_INT, Description: "A sequence number representing a specific generation of the desired state."},

	{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "CreationTimestamp is a timestamp representing the server time when this object was created."},
	{Name: "deletion_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted."},

	{Name: "deletion_grace_period_seconds", Type: proto.ColumnType_INT, Description: "Number of seconds allowed for this object to gracefully terminate before it will be removed from the system.  Only set when deletionTimestamp is also set."},

	{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
	{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."},

	{Name: "owner_references", Type: proto.ColumnType_JSON, Description: "List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller."},
	{Name: "finalizers", Type: proto.ColumnType_JSON, Description: "Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed."},
	{Name: "managed_fields", Type: proto.ColumnType_JSON, Description: "ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field."},
}

// hhmmmmm... why arent these working??
// I dont see the data in the raw response either...
// var typeMetaColumns = []*plugin.Column{
// 	{Name: "kind", Type: proto.ColumnType_STRING, Description: "Kind is a string value representing the REST resource this object represents."},
// 	{Name: "api_version", Type: proto.ColumnType_STRING, Description: "APIVersion defines the versioned schema of this representation of an object."},
// }

// these are resource type specific, and we should "flatten" them...
// SHould we include them raw as well??
var specStatusColumns = []*plugin.Column{
	{Name: "spec", Type: proto.ColumnType_JSON, Description: "Specification of the desired behavior of the namespace. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status."},
	{Name: "status", Type: proto.ColumnType_JSON, Description: "Most recently observed status of the namespace. This data may not be up to date.  Populated by the system.  More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status."},
}

// append the common aws columns for REGIONAL resources onto the column list
func k8sCommonColumns(columns []*plugin.Column) []*plugin.Column {
	allColumns := objectMetadataPrimaryColumns
	allColumns = append(allColumns, columns...)
	//allColumns = append(allColumns, typeMetaColumns...)
	allColumns = append(allColumns, specStatusColumns...)
	allColumns = append(allColumns, objectMetadataSecondaryColumns...)

	return allColumns
}
