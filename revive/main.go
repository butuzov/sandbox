package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheynewallace/tabby"
	"golang.org/x/sync/errgroup"
)

// perform simple integration tests for revive linter

type Stat struct {
	stdOut   bytes.Buffer
	stdErr   bytes.Buffer
	duration time.Duration
}

func main() {
	cwd, _ := os.Getwd()
	Args := []string{"./..."}
	dirs, _ := os.ReadDir(cwd)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		g, _ := errgroup.WithContext(context.Background())

		var Old, New Stat

		g.Go(func() error {
			Bin := filepath.Join(cwd, "revive.old")
			Cwd := filepath.Join(cwd, dir.Name())
			return Exec(&Old, Cwd, Bin, Args...)
		})

		g.Go(func() error {
			Bin := filepath.Join(cwd, "revive.new")
			Cwd := filepath.Join(cwd, dir.Name())
			return Exec(&New, Cwd, Bin, Args...)
		})

		if err := g.Wait(); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		fmt.Println()

		fmt.Printf("= %s %s\n", dir.Name(), strings.Repeat("=", 100))
		t := tabby.New()
		t.AddHeader("Revive", "Errors Found", "Exec Time")

		fmt.Println(Old.stdOut.String())
		// var deltaExecTimeNew =

		t.AddLine("Old", strings.Count(Old.stdOut.String(), "\n"), Old.duration)
		t.AddLine("New", strings.Count(New.stdOut.String(), "\n"), New.duration)
		t.Print()

		fmt.Println()
		fmt.Println()
	}
}

func Exec(stat *Stat, dir, bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Stdout = &stat.stdOut
	cmd.Stderr = &stat.stdOut

	now := time.Now()
	defer func() { stat.duration = time.Since(now) }()

	return cmd.Run()
}
