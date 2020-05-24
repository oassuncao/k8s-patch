package metadata

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"testing"
)

func TestGeneratePatch(t *testing.T) {
	data := &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "None"}}
	SetMetadata(data)

	type args struct {
		current runtime.Object
		obj     runtime.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *PatchData
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				current: data,
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "None"}},
			},
			want: &PatchData{
				PatchType: types.StrategicMergePatchType,
				empty:     true,
			},
		},
		{
			name: "change",
			args: args{
				current: data,
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
			},
			want: &PatchData{
				PatchType: types.StrategicMergePatchType,
				empty:     false,
			},
		},
		{
			name: "changedEqual",
			args: args{
				current: &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
				obj:     &v1.Service{Spec: v1.ServiceSpec{ClusterIP: "Other"}},
			},
			want: &PatchData{
				PatchType: types.StrategicMergePatchType,
				empty:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePatch(tt.args.current, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.PatchType, tt.want.PatchType) {
				t.Errorf("GeneratePatch() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got.empty, tt.want.empty) {
				t.Errorf("GeneratePatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchData_IsEmpty(t *testing.T) {
	type fields struct {
		PatchType types.PatchType
		Data      []byte
		empty     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "false",
			fields: fields{
				empty: false,
			},
			want: false,
		},
		{
			name: "true",
			fields: fields{
				empty: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PatchData{
				PatchType: tt.fields.PatchType,
				Data:      tt.fields.Data,
				empty:     tt.fields.empty,
			}
			if got := p.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newPatchData(t *testing.T) {
	type args struct {
		patchType types.PatchType
		data      []byte
	}
	tests := []struct {
		name string
		args args
		want *PatchData
	}{
		{
			name: "nil",
			args: args{
				data: nil,
			},
			want: &PatchData{},
		},
		{
			name: "empty",
			args: args{
				data:      []byte("{}"),
				patchType: "teste",
			},
			want: &PatchData{
				PatchType: "teste",
				Data:      []byte("{}"),
				empty:     true,
			},
		},
		{
			name: "notEmpty",
			args: args{
				data:      []byte("{ asdada }"),
				patchType: "teste",
			},
			want: &PatchData{
				PatchType: "teste",
				Data:      []byte("{ asdada }"),
				empty:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPatchData(tt.args.patchType, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPatchData() = %v, want %v", got, tt.want)
			}
		})
	}
}
