# app Description

An `app` description can be used to describe an application, to help deploy or run the application on different platforms.

The description of the application contains information related to the application only, independently of
the platform on which the application could be ran or deployed.

An application is composed of one or several components. A component is a program or module 
written in a specific language.

The main information for a Component is the description of the command to use for running, debugging or deploying this application.
This command will be named `run`, `debug` or `deploy`.

This main command can depend on other commands, for example `install` and `build` commands, which must be executed before the component
can be ran or deployed.

Two types of commands can be distinguished: short-runnig commands, and long-running commands. `run` and `debug` commands
are expected to be long-running, and dependent commands are expected to be short-running.

A short-running command will generally generate some artifact (binary, etc).

A long-running command will generally expose some ports.

