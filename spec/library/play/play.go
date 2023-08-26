package play

import (
	"github.com/feloy/app/spec/api"
	"k8s.io/utils/pointer"
)

type Env struct {
	Name  string
	Value string
}

type Source struct {
	Origin *api.CommandDependsOn
	Path   string
}

func NewLocalSource(path string) Source {
	return Source{
		Path: path,
	}
}

func NewSource(component, command string, path string) Source {
	return Source{
		Origin: &api.CommandDependsOn{
			Component: &component,
			Command:   command,
		},
		Path: path,
	}
}

type CommandDetails struct {
	ComponentName string
	Context       string
	Toolkit       api.ComponentToolkit
	CommandName   string
	CommandLine   []string
	Sources       []Source
	Artifacts     []string
	Expose        []api.CommandExpose
	Resources     *api.CommandResources
	Env           []Env
}

func Run(devfile *api.App) ([][]CommandDetails, error) {
	return play(devfile, runCommandName, "nodebug")
}

func Debug(devfile *api.App) ([][]CommandDetails, error) {
	return play(devfile, debugCommandName, "debug")
}

func Deploy(devfile *api.App) ([][]CommandDetails, error) {
	return play(devfile, deployCommandName, "nodebug")
}

type commandNode struct {
	component api.Component
	command   api.Command
	dependsOn []commandNode
}

type componentMap map[string]api.Component
type commandMap map[string]map[string]api.Command

func getComponentMap(devfile *api.App) componentMap {
	result := make(componentMap)
	for _, component := range devfile.Components {
		result[component.Name] = component
	}
	return result
}

func getCommandMap(devfile *api.App) commandMap {
	result := make(commandMap)
	for _, component := range devfile.Components {
		result[component.Name] = make(map[string]api.Command)
		for _, command := range component.Commands {
			result[component.Name][command.Name] = command
		}
	}
	return result
}

func play(devfile *api.App, commandName string, mode string) ([][]CommandDetails, error) {

	components := getComponentMap(devfile)
	commands := getCommandMap(devfile)

	var result [][]CommandDetails

	for _, component := range devfile.Components {
		for _, command := range component.Commands {
			if command.Name == commandName {
				newCommand := commandNode{
					command:   command,
					component: component,
					dependsOn: getDependsOn(components, commands, command.DependsOn),
				}
				result = append(result, getSerialized(newCommand, mode))
			}
		}
	}
	return result, nil
}

func getDependsOn(components componentMap, commands commandMap, dependsOn []api.CommandDependsOn) []commandNode {
	var result []commandNode
	for _, depend := range dependsOn {
		command := commands[*depend.Component][depend.Command]
		newCommand := commandNode{
			component: components[*depend.Component],
			command:   command,
			dependsOn: getDependsOn(components, commands, command.DependsOn),
		}
		result = append(result, newCommand)
	}
	return result
}

func getSerialized(commandsTree commandNode, mode string) []CommandDetails {
	var result []CommandDetails
	for _, depend := range commandsTree.dependsOn {
		result = append(result, getSerialized(depend, mode)...)
	}
	result = append(result, toCommandDetails(commandsTree, mode))
	return result
}

func toCommandDetails(node commandNode, mode string) CommandDetails {
	var sources []Source
	for _, source := range node.command.Sources {
		sources = append(sources, NewLocalSource(source))
	}
	for _, source := range node.dependsOn {
		for _, artifact := range source.command.Artifacts {
			sources = append(sources, NewSource(source.component.Name, source.command.Name, artifact))
		}
	}
	return CommandDetails{
		ComponentName: node.component.Name,
		Context:       pointer.StringDeref(node.component.Context, "."),
		Toolkit:       *node.component.Toolkit,
		CommandName:   node.command.Name,
		CommandLine:   node.command.CommandLine,
		Sources:       sources,
		Artifacts:     node.command.Artifacts,
		Expose:        node.command.Expose,
		Resources:     node.command.Resources,
		Env:           toEnv(node.command.Env, mode),
	}
}

func toEnv(choiceEnvs []api.CommandEnv, mode string) []Env {
	result := make([]Env, 0, len(choiceEnvs))
	for _, env := range choiceEnvs {
		result = append(result, Env{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	return result
}
