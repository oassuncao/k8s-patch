package metadata

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestGetAnnotationName(t *testing.T) {
	SetAnnotationName("test")
	tests := []struct {
		name string
		want string
	}{
		{
			name: "default",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAnnotationName(); got != tt.want {
				t.Errorf("GetAnnotationName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMetadata(t *testing.T) {
	type args struct {
		obj runtime.Object
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{obj: &v1.Service{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetMetadata(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("SetMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}

			service := tt.args.obj.(*v1.Service)
			if s := service.Annotations[GetAnnotationName()]; s == "" {
				t.Errorf("The annotation is not in the object")
			}
		})
	}
}

func Test_getData(t *testing.T) {
	service := &v1.Service{}
	SetMetadata(service)

	type args struct {
		obj runtime.Object
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{
		{
			name:      "default",
			args:      args{obj: service},
			wantEmpty: false,
		},
		{
			name:      "empty",
			args:      args{obj: &v1.Service{}},
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetData(tt.args.obj)
			if len(got) != 0 && tt.wantEmpty {
				t.Errorf("The result is not empty")
			}
		})
	}
}
