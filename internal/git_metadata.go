package internal

import (
	"os/exec"
	"strings"
)

func RepoAddress() (string, error) {
	gitRepoCmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := gitRepoCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func Branch() (string, error) {
	gitBranchCmd := exec.Command("git", "branch", "--show-current")
	out, err := gitBranchCmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}
