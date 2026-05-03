package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/example/precheckin-validator/internal/config"
    "github.com/example/precheckin-validator/internal/gitdiff"
    "github.com/example/precheckin-validator/internal/report"
    "github.com/example/precheckin-validator/internal/runner"
    "github.com/example/precheckin-validator/internal/selector"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: validator run --config validator.yaml [--changed-files a,b] [--output .artifacts]")
        os.Exit(2)
    }
    switch os.Args[1] {
    case "run":
        runCmd(os.Args[2:])
    default:
        fmt.Printf("unknown command: %s
", os.Args[1])
        os.Exit(2)
    }
}

func runCmd(args []string) {
    fs := flag.NewFlagSet("run", flag.ExitOnError)
    configPath := fs.String("config", "validator.yaml", "Path to validator config")
    changedFilesArg := fs.String("changed-files", "", "Comma-separated changed files")
    outputDir := fs.String("output", ".artifacts", "Output directory")
    fs.Parse(args)

    cfg, err := config.Load(*configPath)
    if err != nil {
        fmt.Printf("config error: %v
", err)
        os.Exit(1)
    }

    var changedFiles []string
    if strings.TrimSpace(*changedFilesArg) != "" {
        changedFiles = strings.Split(*changedFilesArg, ",")
    } else {
        changedFiles, err = gitdiff.DefaultChangedFiles()
        if err != nil {
            fmt.Printf("git diff error: %v
", err)
            os.Exit(1)
        }
    }

    plan := selector.BuildPlan(cfg, changedFiles)
    if len(plan.Services) == 0 {
        fmt.Println("no impacted services found; nothing to run")
        os.Exit(0)
    }

    if err := os.MkdirAll(*outputDir, 0o755); err != nil {
        fmt.Printf("output dir error: %v
", err)
        os.Exit(1)
    }

    result, err := runner.Execute(cfg, plan, *outputDir)
    if err != nil {
        fmt.Printf("run error: %v
", err)
        os.Exit(1)
    }

    if err := report.WriteSummary(filepath.Join(*outputDir, "summary.json"), result); err != nil {
        fmt.Printf("summary write error: %v
", err)
        os.Exit(1)
    }

    fmt.Printf("status=%s
", result.Status)
    if result.Status != "passed" {
        os.Exit(1)
    }
}
