package play

import (
	"testing"

	"github.com/feloy/app/spec/library/parse"
)

func Test_play(t *testing.T) {
	type args struct {
		devfile     string
		commandName string
		mode        string
	}
	tests := []struct {
		name    string
		args    args
		want    []commandNode
		wantErr bool
	}{
		{
			name: "run",
			args: args{
				devfile:     "../parse/tests/devfile-correct.yaml",
				commandName: "run",
			},
		},
		{
			name: "run angular-wasm",
			args: args{
				devfile:     "../parse/tests/builder.yaml",
				commandName: "run",
			},
		},
		{
			name: "deploy angular-wasm",
			args: args{
				devfile:     "../parse/tests/builder.yaml",
				commandName: "deploy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devfile, err := parse.Parse(tt.args.devfile)
			if err != nil {
				t.Errorf(("unable to parse input devfile"))
			}

			_, err = play(devfile, tt.args.commandName, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("play() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("play() = %v, want %v", got, tt.want)
			//}
		})
	}
}
