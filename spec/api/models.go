// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package api

import "encoding/json"
import "fmt"

// A description of an Application, to be used for building, testing, running,
// deploying, etc the application
type App struct {
	// List of Components of the Application
	Components []Component `json:"components" yaml:"components" mapstructure:"components"`

	// Version of the Schema used for this document
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion" mapstructure:"schemaVersion"`
}

type Command struct {
	// The files created by a short-running command
	Artifacts []string `json:"artifacts,omitempty" yaml:"artifacts,omitempty" mapstructure:"artifacts,omitempty"`

	// Command to execute
	CommandLine []string `json:"commandLine" yaml:"commandLine" mapstructure:"commandLine"`

	// List of commands to be executed before this command
	DependsOn []CommandDependsOn `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty" mapstructure:"dependsOn,omitempty"`

	// Environment variables to define when running the command
	Env []CommandEnv `json:"env,omitempty" yaml:"env,omitempty" mapstructure:"env,omitempty"`

	// The list of ports exposed by a long-running command
	Expose []CommandExpose `json:"expose,omitempty" yaml:"expose,omitempty" mapstructure:"expose,omitempty"`

	// Name of the command. Either a free name, or one of the pre-defined values
	// 'debug', 'run' or 'deploy'
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Resources (CPU and Memory) necessary to run this command
	Resources *CommandResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// List of source files and directories necessary to execute the command.
	// Directories must terminate with '/'
	Sources []string `json:"sources,omitempty" yaml:"sources,omitempty" mapstructure:"sources,omitempty"`
}

type CommandDependsOn struct {
	// Name of a command
	Command string `json:"command" yaml:"command" mapstructure:"command"`

	// Name of the component defining the command. By default, the current component
	Component *string `json:"component,omitempty" yaml:"component,omitempty" mapstructure:"component,omitempty"`
}

type CommandEnv struct {
	// Name of the environment variable
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Value of the environment variable
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

type CommandExpose struct {
	// If the port is used by the debugger, false by default
	Debug *bool `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`

	// Host on which the port is bound. Default value is the 'localhost' value
	Host *CommandExposeHost `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host,omitempty"`

	// Name of the port
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Number of the Port
	Port CommandExposePort `json:"port" yaml:"port" mapstructure:"port"`

	// If the port is accessible by the end-users, false by default
	Public *bool `json:"public,omitempty" yaml:"public,omitempty" mapstructure:"public,omitempty"`
}

// Definition of a Host. The host can be obtained from, in this order, a
// configuration file, an environment variable or a constant value
type CommandExposeHost struct {
	// The host name, if none of fromFile and fromEnv are defined
	Default *string `json:"default,omitempty" yaml:"default,omitempty" mapstructure:"default,omitempty"`

	// The environment variable containing the host name
	FromEnv *string `json:"fromEnv,omitempty" yaml:"fromEnv,omitempty" mapstructure:"fromEnv,omitempty"`

	// The configuration file containing the host name
	FromFile *string `json:"fromFile,omitempty" yaml:"fromFile,omitempty" mapstructure:"fromFile,omitempty"`
}

// Definition of a Port number. The number can be obtained from, in this order, a
// configuration file, an environment variable or a constant value
type CommandExposePort struct {
	// The port number, if none of fromFile and fromEnv are defined
	Default *int `json:"default,omitempty" yaml:"default,omitempty" mapstructure:"default,omitempty"`

	// The environment variable containing the port number
	FromEnv *string `json:"fromEnv,omitempty" yaml:"fromEnv,omitempty" mapstructure:"fromEnv,omitempty"`

	// The configuration file containing the port number
	FromFile *string `json:"fromFile,omitempty" yaml:"fromFile,omitempty" mapstructure:"fromFile,omitempty"`
}

type CommandResources struct {
	// CPU necessary to run the command
	Cpu *string `json:"cpu,omitempty" yaml:"cpu,omitempty" mapstructure:"cpu,omitempty"`

	// Memory necessary to run the command
	Memory *string `json:"memory,omitempty" yaml:"memory,omitempty" mapstructure:"memory,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CommandExpose) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in CommandExpose: required")
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port in CommandExpose: required")
	}
	type Plain CommandExpose
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CommandExpose(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CommandEnv) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in CommandEnv: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value in CommandEnv: required")
	}
	type Plain CommandEnv
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CommandEnv(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Command) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["commandLine"]; !ok || v == nil {
		return fmt.Errorf("field commandLine in Command: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in Command: required")
	}
	type Plain Command
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Command(plain)
	return nil
}

// Description of the Toolkit to be used to build, test, run the Component
type ComponentToolkit struct {
	// Name of the Toolkit
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Version of the toolkit
	Version *string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComponentToolkit) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in ComponentToolkit: required")
	}
	type Plain ComponentToolkit
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ComponentToolkit(plain)
	return nil
}

// A Component is a single part of the Application. It can be a micro-service, a
// module, an executable, etc
type Component struct {
	// Commands to be used to build, test, run, deploy, etc the Component
	Commands []Command `json:"commands,omitempty" yaml:"commands,omitempty" mapstructure:"commands,omitempty"`

	// Directory containing the sources of the component, relative to the directory
	// containing the App Description. Current directory by default
	Context *string `json:"context,omitempty" yaml:"context,omitempty" mapstructure:"context,omitempty"`

	// Name of the Component
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Description of the Toolkit to be used to build, test, run the Component
	Toolkit *ComponentToolkit `json:"toolkit,omitempty" yaml:"toolkit,omitempty" mapstructure:"toolkit,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Component) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in Component: required")
	}
	type Plain Component
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Component(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CommandDependsOn) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["command"]; !ok || v == nil {
		return fmt.Errorf("field command in CommandDependsOn: required")
	}
	type Plain CommandDependsOn
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CommandDependsOn(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *App) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["components"]; !ok || v == nil {
		return fmt.Errorf("field components in App: required")
	}
	if v, ok := raw["schemaVersion"]; !ok || v == nil {
		return fmt.Errorf("field schemaVersion in App: required")
	}
	type Plain App
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = App(plain)
	return nil
}
