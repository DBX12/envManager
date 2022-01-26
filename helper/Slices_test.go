package helper

import (
	"envManager/internal"
	"testing"
)

func TestSliceStringContains(t *testing.T) {
	type args struct {
		needle   string
		haystack []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Needle exists",
			args: args{
				needle:   "foo",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: true,
		},
		{
			name: "Needle does not exist",
			args: args{
				needle:   "faz",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: false,
		},
		{
			name: "Empty haystack",
			args: args{
				needle:   "foo",
				haystack: []string{},
			},
			want: false,
		},
		{
			name: "Empty needle",
			args: args{
				needle:   "",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: false,
		},
		{
			name: "Empty needle and haystack",
			args: args{
				needle:   "",
				haystack: []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceStringContains(tt.args.needle, tt.args.haystack); got != tt.want {
				t.Errorf("SliceStringContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStringLinearSearch(t *testing.T) {
	type args struct {
		needle   string
		haystack []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Empty haystack",
			args: args{
				needle:   "foo",
				haystack: []string{},
			},
			want: -1,
		},
		{
			name: "Empty needle",
			args: args{
				needle:   "",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: -1,
		},
		{
			name: "Empty needle and haystack",
			args: args{
				needle:   "",
				haystack: []string{},
			},
			want: -1,
		},
		{
			name: "Existing needle",
			args: args{
				needle:   "foo",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: 1,
		},
		{
			name: "Two existing needles return the first instance",
			args: args{
				needle:   "foo",
				haystack: []string{"bar", "foo", "baz", "foo"},
			},
			want: 1,
		},
		{
			name: "Non-existing needle",
			args: args{
				needle:   "faz",
				haystack: []string{"bar", "foo", "baz"},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceStringLinearSearch(tt.args.needle, tt.args.haystack); got != tt.want {
				t.Errorf("SliceStringLinearSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestSliceStringEqual(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty slices",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: true,
		},
		{
			name: "Differing length",
			args: args{
				a: []string{"foo", "bar"},
				b: []string{"foo", "bar", "baz"},
			},
			want: false,
		},
		{
			name: "Equal slices, different order",
			args: args{
				a: []string{"foo", "bar"},
				b: []string{"bar", "foo"},
			},
			want: true,
		},
		{
			name: "Equal slices, equal order",
			args: args{
				a: []string{"foo", "bar"},
				b: []string{"foo", "bar"},
			},
			want: true,
		},
		{
			name: "Unequal slices",
			args: args{
				a: []string{"foo", "bar"},
				b: []string{"foo", "baz"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceStringEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("SliceStringEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
