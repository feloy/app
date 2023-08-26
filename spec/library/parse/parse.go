package parse

import (
	"os"

	"github.com/feloy/app/spec/api"
	"gopkg.in/yaml.v3"
)

func Parse(app string) (*api.App, error) {
	var result api.App
	f, err := os.Open(app)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}
	postParse(&result)
	return &result, nil
}

// postParse completes partial fields:
//   - dependsOn fields with an empty Component
//     get a value for Component with the component of current command
func postParse(appli *api.App) {
	for i := range appli.Components {
		for j := range appli.Components[i].Commands {
			for d := range appli.Components[i].Commands[j].DependsOn {
				depend := &appli.Components[i].Commands[j].DependsOn[d]
				if depend.Component == nil {
					depend.Component = &appli.Components[i].Name
				}
			}
		}
	}
}
