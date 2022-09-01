package main

import (
	"context"
	"flag"
	"fmt"
	dbc "github.com/aquasecurity/trivy/pkg/db"
	"github.com/yyong-brs/trivy-mirror/trivy"
	"os"
	"runtime"
	"sync"
)

// Common variables.
var (
	description string = "Application to mirror trivy."
	gitCommit   string = "n/a"
	name        string = "trivy-mirror"
	source      string = "https://github.com/yyong-brs/trivy-mirror"
)

func main() {
	// Print version.
	if (len(os.Args) > 1) && (os.Args[1] == "version") {
		fmt.Printf("Description:    %s\n", description)
		fmt.Printf("Git Commit:     %s\n", gitCommit)
		fmt.Printf("Go Version:     %s\n", runtime.Version())
		fmt.Printf("Name:           %s\n", name)
		fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Source:         %s\n", source)
		return
	}

	// Print flags related messages to stdout instead of stderr.
	flag.CommandLine.SetOutput(os.Stdout)

	requestWg := &sync.WaitGroup{}
	dbUpdateWg := &sync.WaitGroup{}
	catchDir, _ := os.Getwd()
	worker := trivy.NewDBWorker(dbc.NewClient(catchDir, true, false))
	ctx := context.Background()
	if err := worker.Update(ctx, catchDir, dbUpdateWg, requestWg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(-1)
	}
	fmt.Println("Update DB success ......")
}
