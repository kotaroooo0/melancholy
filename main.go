package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"golang.org/x/xerrors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	imageName := flag.String("i", "", "image name (required)")
	dockerfilePath := flag.String("f", "Dockerfile", "Dockerfile path")
	flag.Parse()

	if *imageName == "" {
		log.Print("-i <image name> is must be required")
		flag.Usage()
		os.Exit(2)
	}

	stageNames, err := readStageNames(*dockerfilePath)
	if err != nil {
		return xerrors.Errorf("read stage names: %w", err)
	}

	fmt.Println("# ----- Build image -----")
	fmt.Printf("docker build -t %s --cache-from=%s --build-arg BUILDKIT_INLINE_CACHE=1 .\n", fmt.Sprintf("%s:%s", *imageName, stageNames[len(stageNames)-1]), strings.Join(stageNames, ","))

	fmt.Println("# ----- Attach tags -----")
	for i, stageName := range stageNames {
		if i == len(stageNames)-1 {
			break
		}
		fmt.Printf("docker build -t %s --target=%s --build-arg BUILDKIT_INLINE_CACHE=1 . &\n", fmt.Sprintf("%s:%s", *imageName, stageName), stageName)
	}
	fmt.Println("wait")

	fmt.Println("# ----- Push images -----")
	for _, stageName := range stageNames {
		fmt.Printf("docker push %s &\n", fmt.Sprintf("%s:%s", *imageName, stageName))
	}
	fmt.Println("wait")
	return nil
}

func readStageNames(dockerfilePath string) ([]string, error) {
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, xerrors.Errorf("open Dockerfile: %w", err)
	}
	defer file.Close()

	parsed, err := parser.Parse(file)
	if err != nil {
		return nil, xerrors.Errorf("parse Dockerfile: %w", err)
	}

	stages, _, err := instructions.Parse(parsed.AST)
	if err != nil {
		return nil, xerrors.Errorf("parse instructions: %w", err)
	}

	stageNames := make([]string, 0, len(stages))
	for _, stage := range stages {
		if stage.Name == "" {
			stageNames = append(stageNames, "latest")
			break
		}
		stageNames = append(stageNames, stage.Name)
	}
	return stageNames, nil
}
