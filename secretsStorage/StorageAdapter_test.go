package secretsStorage

import (
	"reflect"
	"testing"
)

func TestCreateStorageAdapter(t *testing.T) {
	type args struct {
		name   string
		config Storage
	}
	tests := []struct {
		name    string
		args    args
		want    StorageAdapter
		wantErr bool
	}{
		{
			name: "Unknown type",
			args: args{
				name: "testcase01",
				config: Storage{
					StorageType: "unknown",
					Config:      nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		// testing with actual instances like Keepass is hard as DeepEqual does
		// not work correctly with StorageAdapters
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStorageAdapter(tt.args.name, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStorageAdapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStorageAdapter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
