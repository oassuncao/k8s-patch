package compare

import (
	"github.com/oassuncao/k8s-patch/metadata"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
)

func DeepEqualPatch(current, obj runtime.Object) (bool, error) {
	patch, err := metadata.GeneratePatch(current, obj)
	if err != nil {
		return false, err
	}

	return patch.IsEmpty(), nil
}

func DeepEqual(current, obj runtime.Object) (bool, error) {
	initialData, err := metadata.GetData(current)
	if err != nil {
		return false, err
	}

	objMetadata, err := metadata.GetMetadata(obj)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(initialData, objMetadata), nil
}
