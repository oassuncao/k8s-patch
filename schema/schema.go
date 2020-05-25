package schema

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kube-openapi/pkg/util/proto"
)

type Resource interface {
	LookupResource(kind schema.GroupVersionKind) proto.Schema
}

type SchemaResource struct {
	resources map[schema.GroupVersionKind]proto.Schema
}

func (r *SchemaResource) LookupResource(kind schema.GroupVersionKind) proto.Schema {
	if s, found := r.resources[kind]; found {
		return s
	}
	return nil
}

func NewSchemaResource(models proto.Models) SchemaResource {
	resources := SchemaResource{}
	resources.resources = make(map[schema.GroupVersionKind]proto.Schema)
	for _, modelName := range models.ListModels() {
		model := models.LookupModel(modelName)
		kinds := parseGroupVersionKind(model)
		for _, kind := range kinds {
			resources.resources[kind] = model
		}
	}
	return resources
}
