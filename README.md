# app Description

An `app` description can be used to describe an application, to help deploy or run the application on different platforms.

The description of the application contains information related to the application only, independently of
the platform on which the application could be ran or deployed.

## Components

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

## Services

It is also possible to define a list of services, used by the application. Only the specifications of the service are given
(the product and its version).

A component can reference services defined in the `app` description. By referencing a service, a component can indicate which
credentials it expects, to be able work with the service.

## Before / After

Before
```
README
------
- this project uses Go version 1.17 and postgresql 14.
- set the .env file or env var as shown in the .env.example file
- fetch all dependencies using `go mod download`
- run `make run` to start the server locally

.env.example
------------
export PORT=6000
export POSTGRESQL_URL='postgres://user:password@host:port/dbname'
```

After
```
components:
  - name: api
    context: api
    toolkit:
      name: go
      version: "1.17"
    commands:
      - name: dependencies
        sources:
          - go.mod
          - go.sum
        commandLine:
          - go
          - mod
          - download
        
      - name: run
        dependsOn:
          - command: dependencies
        sources:
          - Makefile
          - go.mod
          - go.sum
          - main.go
          - conduit/
          - mock/
          - postgres/
          - server/
        commandLine:
          - make
          - run
        expose:
          - name: api
            port:
              fromEnv: PORT
            public: true
    services:
      - name: database
        connection:
          fromEnv: POSTGRESQL_URL
  
services:
  - name: database
    compatible:
      - name: postgresql
        version: "14"

```

## Focus on the tooling

The `app` description describes only how the application can be built and executed, without
any asumption on the platform on which this will be done.

Imagine the application is intended to be deployed either in a bare-metal system (Linux, Windows, Mac, etc),
or using containers (with Docker Compose or Kubernetes), or on a Platform as a Service (Google Cloud Run, etc).

This description gives the opportunity for the creation of tools specific to platforms.

For example, you can imagine a tool that will make sure that the `go` suite is installed on the local system,
a PostgreSQL v14 service is running, then run the command lines on the local system to build the component,
then install it in a dedicated place (this step is not described in the `app` description, but knowing
where to deploy a n executable is the responsability of the platform maintainer, not of the developer, right?).

Another tool could create a Dockerfile based on the instructions to build the component
([tools/image-builder](./tools/image-builder) is a simple example),
then it could run Docker or Podman (or any other tool) to build the container images,
then push the images to a registry. Then, it could deploy the container (using Docker Compose, Kubernetes, etc),
by injecting the `PORT` and `POSTGRESQL_URL` environment variables with the appropraite values: the tool
could create some Service to make the api accessible on some port, and inject this same value. the tool
could also deploy a PostgreSQL instance with some credentials (coming from different sources, depending
on if the deployment is for develpoment, staging or production) and inject the `POSTGRESQL_URL` with these credentials.
Many information needed in this process (the tool to build images, the image registry, how and where to deploy the database, 
which port numbers to use, etc) are the knowledge of the platform operator, not the developer, right?
