package selector

import (
    "strings"

    "github.com/example/precheckin-validator/internal/config"
)

type Plan struct {
    Services map[string]config.ServiceConfig `json:"services"`
    Tests    []string                        `json:"tests"`
    Changed  []string                        `json:"changed"`
}

func BuildPlan(cfg config.Config, changed []string) Plan {
    plan := Plan{Services: map[string]config.ServiceConfig{}, Changed: changed}
    testSet := map[string]struct{}{}
    for prefix, entry := range cfg.CriticalPaths {
        for _, f := range changed {
            if strings.HasPrefix(f, prefix) {
                for _, s := range entry.Services {
                    if svc, ok := cfg.Services[s]; ok {
                        plan.Services[s] = svc
                    }
                }
                for _, t := range entry.Tests {
                    testSet[t] = struct{}{}
                }
            }
        }
    }
    for t := range testSet {
        plan.Tests = append(plan.Tests, t)
    }
    return plan
}
