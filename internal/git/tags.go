// Package git provides git operations for git-release.
package git

import (
	"os/exec"
	"strings"
	"time"

	"github.com/skaramicke/git-release/internal/semver"
)

// TagInfo holds a parsed version together with its commit metadata.
type TagInfo struct {
	Tag  semver.Version
	Hash string // 7-char short hash
	Date time.Time
}

// ListLocalTags returns all local tags that parse as valid release tags with
// the given prefix. Unparseable tags are silently ignored.
func ListLocalTags(dir, prefix string) ([]semver.Version, error) {
	tags, err := rawLocalTags(dir)
	if err != nil {
		return nil, err
	}
	return filterParseable(tags, prefix), nil
}

// ListLocalTagsWithInfo returns local tags with commit hash and date.
func ListLocalTagsWithInfo(dir, prefix string) ([]TagInfo, error) {
	// format: <tag>\t<short-hash>\t<unix-timestamp>
	c := exec.Command("git", "tag", "--format=%(refname:short)\t%(objectname:short)\t%(creatordate:unix)")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return nil, err
	}

	var result []TagInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}
		tag := parts[0]
		hash := parts[1]
		v, err := semver.ParseTag(tag, prefix)
		if err != nil {
			continue
		}
		var ts int64
		if _, err := strings.NewReader(parts[2]).Read(nil); err == nil {
			ts = parseUnix(parts[2])
		}
		result = append(result, TagInfo{
			Tag:  v,
			Hash: hash,
			Date: time.Unix(ts, 0),
		})
	}
	return result, nil
}

// ListRemoteTags fetches tags from the remote and returns all that parse as
// valid release tags. This is the authoritative source used before every
// mutating operation.
func ListRemoteTags(dir, remote, prefix string) ([]semver.Version, error) {
	c := exec.Command("git", "ls-remote", "--tags", remote)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		// format: <hash>\trefs/tags/<tagname>
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		ref := parts[1]
		// Skip peeled tag refs (^{})
		if strings.HasSuffix(ref, "^{}") {
			continue
		}
		name := strings.TrimPrefix(ref, "refs/tags/")
		tags = append(tags, name)
	}

	return filterParseable(tags, prefix), nil
}

// CurrentBranch returns the name of the currently checked-out branch.
func CurrentBranch(dir string) (string, error) {
	c := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsClean returns true if the working tree has no uncommitted changes.
func IsClean(dir string) (bool, error) {
	c := exec.Command("git", "status", "--porcelain")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "", nil
}

// CreateAndPushTag creates a local tag and pushes it to the remote.
// If dryRun is true, it prints what would happen without doing it.
func CreateAndPushTag(dir, remote, tagName, ref string, sign, dryRun bool) error {
	if dryRun {
		return nil
	}

	args := []string{"tag"}
	if sign {
		args = append(args, "-s")
	}
	args = append(args, tagName, ref)

	c := exec.Command("git", args...)
	c.Dir = dir
	if out, err := c.CombinedOutput(); err != nil {
		return &GitError{cmd: "git tag", output: string(out), err: err}
	}

	push := exec.Command("git", "push", "--no-verify", remote, "refs/tags/"+tagName)
	push.Dir = dir
	if out, err := push.CombinedOutput(); err != nil {
		return &GitError{cmd: "git push", output: string(out), err: err}
	}

	return nil
}

// ResolveRef returns the full commit hash for a ref (tag name, branch, etc.).
func ResolveRef(dir, ref string) (string, error) {
	c := exec.Command("git", "rev-list", "-n", "1", ref)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// --- helpers ---

func rawLocalTags(dir string) ([]string, error) {
	c := exec.Command("git", "tag", "-l")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

func filterParseable(tags []string, prefix string) []semver.Version {
	var result []semver.Version
	for _, tag := range tags {
		v, err := semver.ParseTag(tag, prefix)
		if err == nil {
			result = append(result, v)
		}
	}
	return result
}

func parseUnix(s string) int64 {
	var n int64
	for _, ch := range strings.TrimSpace(s) {
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int64(ch-'0')
	}
	return n
}
