package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// execPrime runs the prime subcommand with the given args and returns stdout.
func execPrime(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := Root()
	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetErr(&strings.Builder{})
	root.SetArgs(append([]string{"prime"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestPrime_NoArgs_PrintsReference(t *testing.T) {
	out, err := execPrime(t)
	require.NoError(t, err)
	assert.Contains(t, out, "git release status")
	assert.Contains(t, out, "git release stage")
}

func TestPrime_UnknownTool(t *testing.T) {
	_, err := execPrime(t, "vim")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tool")
}

// withFakeHome runs fn with HOME set to a temp dir and returns the temp dir path.
func withFakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	return home
}

// ── Claude ────────────────────────────────────────────────────────────────────

func TestPrime_Claude_Global(t *testing.T) {
	home := withFakeHome(t)
	_, err := execPrime(t, "claude", "--global")
	require.NoError(t, err)

	skillPath := filepath.Join(home, ".claude", "skills", "git-release", "SKILL.md")
	data, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "name: git-release")
	assert.Contains(t, content, "git release status")
}

func TestPrime_Claude_Local(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := execPrime(t, "claude")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

func TestPrime_Claude_Local_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("CLAUDE.md", []byte("# existing\n"), 0644))

	_, err := execPrime(t, "claude")
	require.NoError(t, err)

	data, _ := os.ReadFile("CLAUDE.md")
	content := string(data)
	assert.Contains(t, content, "# existing")
	assert.Contains(t, content, "git release status")
}

func TestPrime_Claude_Idempotent(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	require.NoError(t, execPrimeDrop(t, "claude"))
	require.NoError(t, execPrimeDrop(t, "claude"))

	data, _ := os.ReadFile("CLAUDE.md")
	// Marker should appear exactly once (not duplicated)
	assert.Equal(t, 1, strings.Count(string(data), primeSection))
}

func execPrimeDrop(t *testing.T, args ...string) error {
	t.Helper()
	_, err := execPrime(t, args...)
	return err
}

// ── OpenCode ──────────────────────────────────────────────────────────────────

func TestPrime_OpenCode_Global(t *testing.T) {
	home := withFakeHome(t)
	_, err := execPrime(t, "opencode", "--global")
	require.NoError(t, err)

	path := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

func TestPrime_OpenCode_Local(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := execPrime(t, "opencode")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

// ── Copilot ───────────────────────────────────────────────────────────────────

func TestPrime_Copilot_Local(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := execPrime(t, "copilot")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, ".github", "copilot-instructions.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

func TestPrime_Copilot_GlobalWarns(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	root := Root()
	errBuf := &strings.Builder{}
	root.SetOut(&strings.Builder{})
	root.SetErr(errBuf)
	root.SetArgs([]string{"prime", "copilot", "--global"})
	_ = root.Execute()

	// Should still write locally but warn that global is not supported
	assert.Contains(t, errBuf.String(), "global")
}

// ── Cursor ────────────────────────────────────────────────────────────────────

func TestPrime_Cursor_Global(t *testing.T) {
	home := withFakeHome(t)
	_, err := execPrime(t, "cursor", "--global")
	require.NoError(t, err)

	path := filepath.Join(home, ".cursor", "rules", "git-release.md")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

func TestPrime_Cursor_Local(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := execPrime(t, "cursor")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, ".cursor", "rules", "git-release.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")
}

// ── Aider ─────────────────────────────────────────────────────────────────────

func TestPrime_Aider_Global(t *testing.T) {
	home := withFakeHome(t)
	_, err := execPrime(t, "aider", "--global")
	require.NoError(t, err)

	// Should create the prime file
	primePath := filepath.Join(home, ".config", "git-release", "prime.md")
	data, err := os.ReadFile(primePath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")

	// Should add a read: entry to ~/.aider.conf.yml
	conf, err := os.ReadFile(filepath.Join(home, ".aider.conf.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(conf), "prime.md")
}

func TestPrime_Aider_Local(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := execPrime(t, "aider")
	require.NoError(t, err)

	primePath := filepath.Join(dir, ".git-release-prime.md")
	data, err := os.ReadFile(primePath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "git release status")

	conf, err := os.ReadFile(filepath.Join(dir, ".aider.conf.yml"))
	require.NoError(t, err)
	assert.Contains(t, string(conf), ".git-release-prime.md")
}
