package metadata

import "testing"

type MyStruct struct {
	name string
}

func Test_getTypeName(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{obj: MyStruct{}},
			want: "MyStruct",
		},
		{
			name: "pointer",
			args: args{obj: &MyStruct{}},
			want: "MyStruct",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTypeName(tt.args.obj); got != tt.want {
				t.Errorf("getTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDefault(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string",
			args: args{value: ""},
			want: true,
		},
		{
			name: "notEmpty",
			args: args{value: " "},
			want: false,
		},
		{
			name: "nil",
			args: args{value: nil},
			want: true,
		},
		{
			name: "bool",
			args: args{value: false},
			want: true,
		},
		{
			name: "struct",
			args: args{value: MyStruct{name: "Value"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDefault(tt.args.value); got != tt.want {
				t.Errorf("isDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
