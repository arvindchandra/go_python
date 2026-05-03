package runner

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/example/precheckin-validator/internal/config"
    "github.com/example/precheckin-validator/internal/selector"
)

type Result struct {
    Status   string   `json:"status"`
    Services []string `json:"services"`
    Command  []string `json:"command"`
    Output   string   `json:"output"`
}

func Execute(cfg config.Config, plan selector.Plan, outputDir string) (Result, error) {
    var result Result
    for serviceName, svc := range plan.Services {
        result.Services = append(result.Services, serviceName)
        junitHostDir, err := filepath.Abs(outputDir)
        if err != nil {
            return result, err
        }
        env := os.Environ()
        for k, v := range cfg.Container.Env {
            env = append(env, fmt.Sprintf("%s=%s", k, v))
        }
        cmdArgs := append([]string{"compose", "-f", cfg.Container.ComposeFile, "run", "--rm", cfg.Container.Service}, svc.TestCommand...)
        cmd := exec.Command("docker", cmdArgs...)
        cmd.Env = append(env, fmt.Sprintf("ARTIFACT_DIR=%s", junitHostDir))
        var buf bytes.Buffer
        cmd.Stdout = &buf
        cmd.Stderr = &buf
        result.Command = append([]string{"docker"}, cmdArgs...)
        err = cmd.Run()
        result.Output = buf.String()
        if err != nil {
            result.Status = "failed"
            return result, nil
        }
    }
    result.Status = "passed"
    return result, nil
}
