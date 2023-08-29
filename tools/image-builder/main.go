package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/feloy/app/spec/library/parse"
	"github.com/feloy/app/spec/library/play"

	"github.com/feloy/app/tools/image-builder/pkg/image"
)

func usage() {
	fmt.Printf("usage: %s <run|debug|deploy> [--file app.yaml]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var devfilePath string
	flag.StringVar(&devfilePath, "file", "app.yaml", "app file name")
	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
	}

	command := flag.Args()[0]
	if command != "run" && command != "debug" && command != "deploy" {
		usage()
	}

	devfile, err := parse.Parse(devfilePath)
	if err != nil {
		panic(err)
	}

	var images [][]play.CommandDetails

	switch command {
	case "run":
		images, err = play.Run(devfile)
	case "debug":
		images, err = play.Debug(devfile)
	case "deploy":
		images, err = play.Deploy(devfile)
	}
	if err != nil {
		panic(err)
	}
	for _, img := range images {
		content, err := image.Build(img, command)
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("%s-gen.Dockerfile", img[len(img)-1].ComponentName)
		os.WriteFile(filename, content, 0644)
	}
}
