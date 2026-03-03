package commands_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/skaramicke/git-release/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testRepo creates a temporary git repo, sets the working directory to it,
// and returns a cleanup function that restores the previous directory.
func testRepo(t *testing.T) (dir string, restore func()) {
	t.Helper()
	dir = t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		out, err := c.CombinedOutput()
		require.NoError(t, err, string(out))
	}

	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	return dir, func() { os.Chdir(orig) }
}

func tag(t *testing.T, dir, name string) {
	t.Helper()
	c := exec.Command("git", "tag", name)
	c.Dir = dir
	out, err := c.CombinedOutput()
	require.NoError(t, err, string(out))
}

// runCmd executes a root command with the given args and captures output.
func runCmd(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	root := commands.Root()
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs(args)
	err = root.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestStatus_Empty(t *testing.T) {
	_, restore := testRepo(t)
	defer restore()

	out, _, err := runCmd(t, "status")
	require.NoError(t, err)
	assert.Contains(t, out, "none")
	assert.Contains(t, out, "v0.0.1-rc")
}

func TestStatus_WithTags(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()

	tag(t, dir, "v1.2.2")
	tag(t, dir, "v1.3.0-rc.2")

	out, _, err := runCmd(t, "status")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.2.2")
	assert.Contains(t, out, "v1.3.0-rc.2")
	assert.Contains(t, out, "v1.3.0-rc.3")
	assert.Contains(t, out, "v1.3.0")
}

func TestLs_Empty(t *testing.T) {
	_, restore := testRepo(t)
	defer restore()

	out, _, err := runCmd(t, "ls")
	require.NoError(t, err)
	assert.Contains(t, out, "no release tags found")
}

func TestLs_WithTags(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()

	tag(t, dir, "v1.2.2")
	tag(t, dir, "v1.3.0-rc")
	tag(t, dir, "v1.3.0-rc.2")

	out, _, err := runCmd(t, "ls")
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	// Should appear in descending order
	var tagLines []string
	for _, l := range lines {
		if strings.Contains(l, "v1.") {
			tagLines = append(tagLines, l)
		}
	}
	require.GreaterOrEqual(t, len(tagLines), 3)
	assert.Contains(t, tagLines[0], "v1.3.0-rc.2")
	assert.Contains(t, tagLines[1], "v1.3.0-rc")
	assert.Contains(t, tagLines[2], "v1.2.2")
}

func TestStageDryRun_NoRC(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")

	out, _, err := runCmd(t, "stage", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.2.3-rc")
	assert.Contains(t, out, "dry-run")
}

func TestStageDryRun_ExistingRC_Increments(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")
	tag(t, dir, "v1.2.3-rc")

	out, _, err := runCmd(t, "stage", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.2.3-rc.2")
}

func TestStageDryRun_Minor(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")

	out, _, err := runCmd(t, "stage", "minor", "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.3.0-rc")
}

func TestStageDryRun_HigherScopeError(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")
	tag(t, dir, "v2.0.0-rc")

	_, stderr, err := runCmd(t, "stage", "minor", "--dry-run")
	assert.Error(t, err)
	assert.Contains(t, stderr, "higher-scope")
}

func TestReleaseDryRun_NoRC(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")

	out, _, err := runCmd(t, "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.2.3")
	assert.Contains(t, out, "dry-run")
}

func TestReleaseDryRun_PromotesRC(t *testing.T) {
	dir, restore := testRepo(t)
	defer restore()
	tag(t, dir, "v1.2.2")
	tag(t, dir, "v1.3.0-rc.2")

	out, _, err := runCmd(t, "--dry-run")
	require.NoError(t, err)
	assert.Contains(t, out, "v1.3.0")
}

func TestInvalidScope(t *testing.T) {
	_, restore := testRepo(t)
	defer restore()

	_, _, err := runCmd(t, "stage", "weekly")
	assert.Error(t, err)
}
