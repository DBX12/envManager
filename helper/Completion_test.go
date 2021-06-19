package helper

import (
	"reflect"
	"testing"
)

func TestCompletion(t *testing.T) {
	type args struct {
		possibleValues []string
		excludedValues []string
		withPrefix     string
	}
	defaultPossibleValues := []string{"foo", "bar", "baz", "faz"}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "No exclusions and no prefix",
			args: args{
				possibleValues: defaultPossibleValues,
				excludedValues: nil,
				withPrefix:     "",
			},
			want: defaultPossibleValues,
		},
		{
			name: "With exclusions and no prefix",
			args: args{
				possibleValues: defaultPossibleValues,
				excludedValues: []string{"foo"},
				withPrefix:     "",
			},
			want: []string{"bar", "baz", "faz"},
		},
		{
			name: "With prefix and no exclusions",
			args: args{
				possibleValues: defaultPossibleValues,
				excludedValues: nil,
				withPrefix:     "ba",
			},
			want: []string{"bar", "baz"},
		},
		{
			name: "With prefix and exclusions",
			args: args{
				possibleValues: defaultPossibleValues,
				excludedValues: []string{"bar"},
				withPrefix:     "ba",
			},
			want: []string{"baz"},
		},
		{
			name: "With empty slice for exclusions",
			args: args{
				possibleValues: defaultPossibleValues,
				excludedValues: []string{},
				withPrefix:     "",
			},
			want: defaultPossibleValues,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Completion(tt.args.possibleValues, tt.args.excludedValues, tt.args.withPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Completion() = %v, want %v", got, tt.want)
			}
		})
	}
}
