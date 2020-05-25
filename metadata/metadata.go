package metadata

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var metadataAccessor = meta.NewAccessor()
var dataAnnotation = "kubernetes-patch"

func SetAnnotationName(name string) {
	dataAnnotation = name
}

func GetAnnotationName() string {
	return dataAnnotation
}

func SetMetadata(obj runtime.Object) error {
	metadata, err := getMetadata(obj)
	if err != nil {
		return err
	}

	annotations, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return err
	}

	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[dataAnnotation] = string(metadata)
	err = metadataAccessor.SetAnnotations(obj, annotations)
	if err != nil {
		return err
	}

	return nil
}

func getMetadata(obj runtime.Object) ([]byte, error) {
	data, err := getMetadataMap(obj)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func getMetadataMap(obj runtime.Object) (map[string]interface{}, error) {
	data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, obj)
	if err != nil {
		return nil, err
	}

	dataMap := map[string]interface{}{}
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, err
	}

	removeDefaultValues(dataMap)
	return dataMap, nil
}

func getData(obj runtime.Object) ([]byte, error) {
	annotations, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return nil, err
	}

	if annotations != nil {
		if value, found := annotations[dataAnnotation]; found {
			return []byte(value), nil
		}
	}

	return nil, nil
}

func removeDefaultValuesType(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return removeDefaultValues(value.(map[string]interface{}))
	case []interface{}:
		return removeDefaultValuesSlice(value.([]interface{}))
	default:
		return true
	}
}

func removeDefaultValues(data map[string]interface{}) bool {
	if data == nil || len(data) == 0 {
		return false
	}

	for key, value := range data {
		if isDefault(value) || !removeDefaultValuesType(value) {
			delete(data, key)
		}
	}
	return true
}

func removeDefaultValuesSlice(data []interface{}) bool {
	if data == nil || len(data) == 0 {
		return false
	}

	var indexesToRemove []int
	for i := 0; i < len(data); i++ {
		value := data[i]
		if isDefault(value) || !removeDefaultValuesType(value) {
			indexesToRemove = append(indexesToRemove, i)
		}
	}

	if len(indexesToRemove) != 0 {
		for index := range indexesToRemove {
			removeIndex(data, index)
		}
	}
	return true
}

func removeIndex(data []interface{}, index int) []interface{} {
	return append(data[:index], data[index+1:]...)
}
