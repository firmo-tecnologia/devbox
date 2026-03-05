package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Image        string
	NoPull       bool
	Shell        bool
	WorkspaceDir string
	GitRoot      string
	ClaudeDir    string
	ClaudeJSON   string
	DotbinsConf  string
	DotbinsCache string
}

func New(image string, noPull bool, shell bool, workspaceDir string, gitRoot string, dotbinsConf string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}

	if dotbinsConf == "" {
		dotbinsConf = filepath.Join(home, ".dotbins", "dotbins.yaml")
	}

	return &Config{
		Image:        image,
		NoPull:       noPull,
		Shell:        shell,
		WorkspaceDir: workspaceDir,
		GitRoot:      gitRoot,
		ClaudeDir:    filepath.Join(home, ".claude"),
		ClaudeJSON:   filepath.Join(home, ".claude.json"),
		DotbinsConf:  dotbinsConf,
		DotbinsCache: filepath.Join(home, ".devbox", "dotbins"),
	}, nil
}

func (c *Config) EnsureDirs() error {
	dirs := []string{c.ClaudeDir, c.DotbinsCache}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}
	if err := ensureFile(c.ClaudeJSON); err != nil {
		return fmt.Errorf("creating %s: %w", c.ClaudeJSON, err)
	}
	return nil
}

func ensureFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return f.Close()
}

func (c *Config) HasDotbinsConf() bool {
	_, err := os.Stat(c.DotbinsConf)
	return err == nil
}
