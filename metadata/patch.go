package metadata

import (
	"github.com/oassuncao/k8s-patch/schema"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

var resources schema.Resource

type PatchData struct {
	PatchType types.PatchType
	Data      []byte
	empty     bool
}

func (p *PatchData) IsEmpty() bool {
	return p.empty
}

func newPatchData(patchType types.PatchType, data []byte) *PatchData {
	patchData := &PatchData{
		PatchType: patchType,
		Data:      data,
	}

	patchData.empty = data != nil && string(data) == "{}"
	return patchData
}

func GeneratePatch(current, obj runtime.Object) (*PatchData, error) {
	initialData, err := getData(current)
	if err != nil {
		return nil, err
	}

	err = SetMetadata(obj)
	if err != nil {
		return nil, err
	}

	objMetadata, err := getMetadata(obj)
	if err != nil {
		return nil, err
	}

	currentMetadata, err := getMetadata(current)
	if err != nil {
		return nil, err
	}

	patchMetadata, err := getLookupPatchMeta(current)
	if err != nil {
		return nil, err
	}

	patch, err := strategicpatch.CreateThreeWayMergePatch(initialData, objMetadata, currentMetadata, patchMetadata, true)
	if err != nil {
		return nil, err
	}

	return newPatchData(types.StrategicMergePatchType, patch), nil
}

func getLookupPatchMeta(obj runtime.Object) (strategicpatch.LookupPatchMeta, error) {
	if resources != nil {
		resource := resources.LookupResource(obj.GetObjectKind().GroupVersionKind())
		if resource != nil {
			return strategicpatch.NewPatchMetaFromOpenAPI(resource), nil
		}
	}

	patchMetadata, err := strategicpatch.NewPatchMetaFromStruct(obj)
	if err != nil {
		return nil, err
	}

	return patchMetadata, err
}

func SetResource(r schema.Resource) {
	resources = r
}
