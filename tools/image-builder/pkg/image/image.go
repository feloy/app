package image

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/feloy/app/spec/library/play"
)

func Build(steps []play.CommandDetails, mode string) ([]byte, error) {
	buf := &bytes.Buffer{}
	workdirs := map[string]string{}

	for _, step := range steps {
		stepName := fmt.Sprintf("%s-%s", step.ComponentName, step.CommandName)
		fromImage := "universal"
		workdir := ""

		switch step.Toolkit.Name {
		case "go":
			if step.CommandName == "run" && mode == "nodebug" {
				// Image for running go app
				fromImage = "ubi8/ubi-micro"
			} else {
				// Image for building go
				fromImage = "registry.access.redhat.com/ubi9/go-toolset"
				if step.Toolkit.Version != nil {
					fromImage = fromImage + ":" + *step.Toolkit.Version
				}
				workdir = "/opt/app-root/src"
			}

		case "nodejs":
			fromImage = "node"
			if step.Toolkit.Version != nil {
				fromImage = fromImage + ":" + *step.Toolkit.Version
			}
			workdir = "/app"

		case "jdk":
			fromImage = "registry.access.redhat.com/ubi8/openjdk-17:1.16-1"
			workdir = "/home/jboss"
		}
		fmt.Fprintf(buf, "FROM %s AS %s-%s\n\n", fromImage, step.ComponentName, step.CommandName)

		if workdir != "" {
			workdirs[stepName] = workdir
			fmt.Fprintf(buf, "WORKDIR \"%s\"\n\n", workdir)
		}

		for _, source := range step.Sources {
			dest := source.Path
			if strings.Contains(dest, "*") {
				dest = filepath.Dir(dest) + "/"
			}
			if source.Origin == nil {
				fmt.Fprintf(buf, "COPY %s/%s %s\n", step.Context, source.Path, dest)
			} else {
				sourceName := fmt.Sprintf("%s-%s", *source.Origin.Component, source.Origin.Command)
				fmt.Fprintf(buf, "COPY --from=%s %s/%s %s\n", sourceName, workdirs[sourceName], source.Path, dest)
			}
		}
		if len(step.Sources) > 0 {
			fmt.Fprintln(buf)
		}

		if workdir != "" {
			fmt.Fprintf(buf, "ENV WORKDIR=%q\n", workdir)
		}
		for _, env := range step.Env {
			fmt.Fprintf(buf, "ENV %s=%q\n", env.Name, env.Value)
		}
		if workdir != "" || len(step.Env) > 0 {
			fmt.Fprintln(buf)
		}

		if step.CommandLine != nil {
			if step.CommandName == "run" || step.CommandName == "debug" {
				cmd := "CMD ["
				for i, cmdline := range step.CommandLine {
					if i > 0 {
						cmd += ", "
					}
					cmd += "\"" + cmdline + "\""
				}
				cmd += "]"
				fmt.Fprint(buf, cmd)
			} else {
				// deploy
				fmt.Fprintf(buf, "RUN %s\n\n", strings.Join(step.CommandLine, " "))
			}
		}
	}
	return buf.Bytes(), nil
}
