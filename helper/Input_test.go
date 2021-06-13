package helper

import (
	"reflect"
	"testing"
)

func TestInput_getPresetInputValue(t *testing.T) {
	type fields struct {
		Inputs []string
	}
	tests := []struct {
		name       string
		fields     fields
		want       string
		wantPanic  bool
		wantInputs []string
	}{
		{
			name:       "Inputs is unset",
			want:       "",
			wantPanic:  true,
			wantInputs: nil,
		},
		{
			name:       "Inputs is an empty slice",
			fields:     fields{Inputs: []string{}},
			want:       "",
			wantPanic:  true,
			wantInputs: nil,
		},
		{
			name:       "Inputs contains two items",
			fields:     fields{Inputs: []string{"foo", "bar"}},
			want:       "foo",
			wantPanic:  false,
			wantInputs: []string{"bar"},
		},
		{
			name:       "Inputs contains one item",
			fields:     fields{Inputs: []string{"foo"}},
			want:       "foo",
			wantPanic:  false,
			wantInputs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				Inputs: tt.fields.Inputs,
			}
			defer func() {
				r := recover()
				if r != nil && !tt.wantPanic {
					t.Error("Got a panic but wanted none")
				} else if r == nil && tt.wantPanic {
					t.Error("Got no panic but wanted one")
				}
			}()
			if got := i.getPresetInputValue(); got != tt.want {
				t.Fatalf("getPresetInputValue() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.wantInputs, i.Inputs) {
				t.Errorf("Inputs = %v, want %v", i.Inputs, tt.wantInputs)
			}
		})
	}
}

func TestInput_hasPresetInputValues(t *testing.T) {
	type fields struct {
		Inputs []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Inputs is unset",
			want: false,
		},
		{
			name:   "Inputs is an empty slice",
			fields: fields{Inputs: []string{}},
			want:   false,
		},
		{
			name:   "Inputs is a filled slice",
			fields: fields{Inputs: []string{"foo", "bar"}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				Inputs: tt.fields.Inputs,
			}
			if got := i.hasPresetInputValues(); got != tt.want {
				t.Errorf("hasPresetInputValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
