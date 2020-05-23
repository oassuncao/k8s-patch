package compare

import (
	"github.com/oassuncao/k8s-patch/metadata"
	"k8s.io/apimachinery/pkg/runtime"
)

func DeepEqual(current, obj runtime.Object) (bool, error) {
	patch, err := metadata.GeneratePatch(current, obj)
	if err != nil {
		return false, err
	}

	return patch.IsEmpty(), nil
}
