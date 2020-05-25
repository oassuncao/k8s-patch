package compare

import (
	"github.com/oassuncao/k8s-patch/metadata"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestDeepEqualPatch(t *testing.T) {
	data := &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "None"}}
	metadata.SetMetadata(data)

	type args struct {
		current runtime.Object
		obj     runtime.Object
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				current: data,
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "None"}},
			},
			want: true,
		},
		{
			name: "change",
			args: args{
				current: data,
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
			},
			want: false,
		},
		{
			name: "changedEqual",
			args: args{
				current: &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeepEqualPatch(tt.args.current, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepEqual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeepEqual() got = %v, want %v", got, tt.want)
			}
		})
	}
}
