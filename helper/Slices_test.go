package helper

import (
	"envManager/internal"
	"testing"
)

func TestSliceStringRemove(t *testing.T) {
	type args struct {
		value string
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty slice",
			args: args{
				value: "foo",
				slice: []string{},
			},
			want: []string{},
		},
		{
			name: "Empty value",
			args: args{
				value: "",
				slice: []string{"bar", "foo", "baz"},
			},
			want: []string{"bar", "foo", "baz"},
		},
		{
			name: "Empty value and slice",
			args: args{
				value: "",
				slice: []string{},
			},
			want: []string{},
		},
		{
			name: "Existing value",
			args: args{
				value: "foo",
				slice: []string{"bar", "foo", "baz"},
			},
			want: []string{"bar", "baz"},
		},
		{
			name: "Two existing values, both get removed",
			args: args{
				value: "foo",
				slice: []string{"bar", "foo", "baz", "foo"},
			},
			want: []string{"bar", "baz"},
		},
		{
			name: "Non-existing value",
			args: args{
				value: "faz",
				slice: []string{"bar", "foo", "baz"},
			},
			want: []string{"bar", "foo", "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceStringRemove(tt.args.value, tt.args.slice)
			if len(got) == len(tt.want) && len(got) == 0 {
				// we got an empty slice and wanted an empty slice, all is good
				return
			}
			if !internal.AssertStringSliceEqual(t, tt.want, got) {
				t.Errorf("SliceStringRemove() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("do not motify input slice", func(t *testing.T) {
		slc := []string{"bar", "foo", "baz"}
		got := SliceStringRemove("foo", slc)
		want := []string{"bar", "baz"}

		if !internal.AssertStringSliceEqual(t, want, got) {
			t.Errorf("SliceStringRemove() = %v, want %v", got, want)
		}

		if !internal.AssertStringSliceEqual(t, []string{"bar", "foo", "baz"}, slc) {
			t.Errorf("SliceStringRemove() = %v, want %v", slc, []string{"bar", "foo", "baz"})
		}
	})

}

func TestSliceStringUnique(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Normal",
			args: args{input: []string{"val1", "val2", "val1", "val3"}},
			want: []string{"val1", "val2", "val3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceStringUnique(tt.args.input)
			if !internal.AssertStringSliceEqual(t, tt.want, got) {
				t.Errorf("SliceStringUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}
