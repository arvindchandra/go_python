package gitdiff

import (
    "bytes"
    "os/exec"
    "strings"
)

func DefaultChangedFiles() ([]string, error) {
    cmd := exec.Command("git", "diff", "--name-only", "HEAD~1..HEAD")
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    lines := strings.Split(strings.TrimSpace(out.String()), "
")
    var files []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            files = append(files, line)
        }
    }
    return files, nil
}
