package config

import (
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Version        int                          `yaml:"version"`
    Project        string                       `yaml:"project"`
    Container      ContainerConfig              `yaml:"container"`
    QuarantineFile string                       `yaml:"quarantine_file"`
    CriticalPaths  map[string]CriticalPathEntry `yaml:"critical_paths"`
    Services       map[string]ServiceConfig     `yaml:"services"`
}

type ContainerConfig struct {
    ComposeFile string            `yaml:"compose_file"`
    Service     string            `yaml:"service"`
    Workdir     string            `yaml:"workdir"`
    Env         map[string]string `yaml:"env"`
}

type CriticalPathEntry struct {
    Services []string `yaml:"services"`
    Tests    []string `yaml:"tests"`
}

type ServiceConfig struct {
    ImageBuild  bool     `yaml:"image_build"`
    TestCommand []string `yaml:"test_command"`
}

func Load(path string) (Config, error) {
    var cfg Config
    b, err := os.ReadFile(path)
    if err != nil {
        return cfg, err
    }
    err = yaml.Unmarshal(b, &cfg)
    return cfg, err
}
