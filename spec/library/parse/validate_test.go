package parse

import (
	_ "embed"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestValidateSchema(t *testing.T) {
	type args struct {
		devfile string
	}
	tests := []struct {
		name    string
		args    args
		want    func(t *testing.T, errors []gojsonschema.ResultError)
		wantErr bool
	}{
		{
			name: "validation ok",
			args: args{
				devfile: "tests/devfile-correct.yaml",
			},
			wantErr: false,
			want: func(t *testing.T, errors []gojsonschema.ResultError) {
				if len(errors) != 0 {
					t.Errorf("Expected 0 error, got %d", len(errors))
				}
			},
		},
		{
			name: "validation error",
			args: args{
				devfile: "tests/devfile-validate-errors.yaml",
			},
			wantErr: false,
			want: func(t *testing.T, errors []gojsonschema.ResultError) {
				if len(errors) != 1 {
					t.Errorf("Expected 1 error, got %d", len(errors))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSchema(tt.args.devfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}
