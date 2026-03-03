// Package config reads git-release configuration from git config keys.
package config

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// ErrNotInRepo is returned when the given directory is not inside a git repo.
var ErrNotInRepo = errors.New("not inside a git repository")

// Config holds all git-release configuration for a repository.
type Config struct {
	TagPrefix     string // default: "v"
	Remote        string // default: "origin"
	ReleaseBranch string // default: "" (any branch)
	RequireClean  bool   // default: true
	SignTags      bool   // default: false
}

// Load reads git-release config from the repo that contains dir.
func Load(dir string) (*Config, error) {
	if !isRepo(dir) {
		return nil, ErrNotInRepo
	}

	cfg := &Config{
		TagPrefix:    "v",
		Remote:       "origin",
		RequireClean: true,
	}

	cfg.TagPrefix = gitConfigStr(dir, "release.tagPrefix", "v")
	cfg.Remote = gitConfigStr(dir, "release.remote", "origin")
	cfg.ReleaseBranch = gitConfigStr(dir, "release.releaseBranch", "")
	cfg.RequireClean = gitConfigBool(dir, "release.requireClean", true)
	cfg.SignTags = gitConfigBool(dir, "release.signTags", false)

	return cfg, nil
}

func isRepo(dir string) bool {
	c := exec.Command("git", "rev-parse", "--git-dir")
	c.Dir = dir
	return c.Run() == nil
}

func gitConfigStr(dir, key, fallback string) string {
	c := exec.Command("git", "config", "--get", key)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return fallback
	}
	val := strings.TrimSpace(string(out))
	if val == "" {
		return fallback
	}
	return val
}

func gitConfigBool(dir, key string, fallback bool) bool {
	c := exec.Command("git", "config", "--get", key)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return fallback
	}
	b, err := strconv.ParseBool(strings.TrimSpace(string(out)))
	if err != nil {
		return fallback
	}
	return b
}
