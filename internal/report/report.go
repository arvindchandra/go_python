package report

import (
    "encoding/json"
    "os"

    "github.com/example/precheckin-validator/internal/runner"
)

func WriteSummary(path string, result runner.Result) error {
    b, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, b, 0o644)
}
