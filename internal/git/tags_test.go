package git_test

import (
	"os/exec"
	"testing"

	"github.com/skaramicke/git-release/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
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
	return dir
}

func addTag(t *testing.T, dir, tag string) {
	t.Helper()
	c := exec.Command("git", "tag", tag)
	c.Dir = dir
	out, err := c.CombinedOutput()
	require.NoError(t, err, string(out))
}

func TestListLocalTags_Empty(t *testing.T) {
	dir := makeRepo(t)
	tags, err := git.ListLocalTags(dir, "v")
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestListLocalTags(t *testing.T) {
	dir := makeRepo(t)
	addTag(t, dir, "v1.2.3")
	addTag(t, dir, "v1.2.3-rc")
	addTag(t, dir, "v1.3.0-rc.2")
	addTag(t, dir, "not-a-release") // should be filtered out

	tags, err := git.ListLocalTags(dir, "v")
	require.NoError(t, err)

	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		tagStrings[i] = tag.String("v")
	}

	assert.ElementsMatch(t, []string{"v1.2.3", "v1.2.3-rc", "v1.3.0-rc.2"}, tagStrings)
}

func TestListLocalTags_CustomPrefix(t *testing.T) {
	dir := makeRepo(t)
	addTag(t, dir, "rel1.0.0")
	addTag(t, dir, "v1.0.0") // wrong prefix, should be filtered

	tags, err := git.ListLocalTags(dir, "rel")
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Equal(t, "rel1.0.0", tags[0].String("rel"))
}

func TestTagCommitInfo(t *testing.T) {
	dir := makeRepo(t)
	addTag(t, dir, "v1.0.0")

	tags, err := git.ListLocalTagsWithInfo(dir, "v")
	require.NoError(t, err)
	require.Len(t, tags, 1)

	info := tags[0]
	assert.Equal(t, "v1.0.0", info.Tag.String("v"))
	assert.NotEmpty(t, info.Hash)
	assert.Len(t, info.Hash, 7)
	assert.False(t, info.Date.IsZero())
}
