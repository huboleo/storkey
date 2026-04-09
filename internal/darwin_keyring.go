package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Save(deleteFiles bool) error {
	repo, err := RepoAddress()
	if err != nil {
		return err
	}
	branch, err := Branch()
	if err != nil {
		return err
	}
	repoRoot, err := RepoRoot()
	if err != nil {
		return err
	}
	scanDir, err := os.Getwd()
	if err != nil {
		return err
	}

	account := Account(repo, branch)
	indexFile := Index()
	scope, err := filepath.Rel(repoRoot, scanDir)
	if err != nil {
		return err
	}
	var paths []string

	err = filepath.WalkDir(scanDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasPrefix(d.Name(), ".env") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b64 := base64.StdEncoding.EncodeToString(b)
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		service := Service(rel)
		paths = append(paths, rel)
		cmdString := fmt.Sprintf("add-generic-password -a %q -s %q -w %q -U\n", account, service, b64)
		cmd := exec.Command("security", "-i")
		cmd.Stdin = strings.NewReader(cmdString)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to save env file %q: %w (output: %s)", rel, err, strings.TrimSpace(string(out)))
		}
		return nil
	})
	if err != nil {
		return err
	}

	existingPaths, err := readIndexPaths(account, indexFile)
	if err != nil {
		return err
	}

	pathsToSave := mergeIndexPaths(existingPaths, paths, scope)

	idxJSON, err := json.Marshal(pathsToSave)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"security",
		"add-generic-password",
		"-a", account,
		"-s", indexFile,
		"-w", string(idxJSON),
		"-U",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to save index file: %w (output: %s)", err, strings.TrimSpace(string(out)))
	}

	if deleteFiles {
		for _, f := range paths {
			err := os.Remove(filepath.Join(repoRoot, f))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Pull() error {
	repo, err := RepoAddress()
	if err != nil {
		return err
	}
	branch, err := Branch()
	if err != nil {
		return err
	}
	repoRoot, err := RepoRoot()
	if err != nil {
		return err
	}
	scanDir, err := os.Getwd()
	if err != nil {
		return err
	}

	scope, err := filepath.Rel(repoRoot, scanDir)
	if err != nil {
		return err
	}

	account := Account(repo, branch)
	indexFile := Index()

	paths, err := readIndexPaths(account, indexFile)
	if err != nil {
		return err
	}

	for _, rel := range dedupePaths(paths) {
		if !isPathInScope(rel, scope) {
			continue
		}

		cmd := exec.Command(
			"security",
			"find-generic-password",
			"-a", account,
			"-s", Service(rel),
			"-w",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to read env file %q: %w (output: %s)", rel, err, strings.TrimSpace(string(out)))
		}

		b64 := strings.TrimSpace(string(out))
		content, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return fmt.Errorf("failed to decode env file %q: %w", rel, err)
		}

		target := filepath.Join(repoRoot, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("failed to create directories for %q: %w", rel, err)
		}
		if err := os.WriteFile(target, content, 0o600); err != nil {
			return fmt.Errorf("failed to write env file %q: %w", rel, err)
		}
	}

	return nil
}

func readIndexPaths(account string, indexFile string) ([]string, error) {
	cmd := exec.Command(
		"security",
		"find-generic-password",
		"-a", account,
		"-s", indexFile,
		"-w",
	)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		if strings.Contains(output, "could not be found") {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read index file: %w (output: %s)", err, output)
	}
	if output == "" {
		return nil, nil
	}

	var paths []string
	if err := json.Unmarshal([]byte(output), &paths); err != nil {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	return paths, nil
}

func mergeIndexPaths(existingPaths []string, scannedPaths []string, scope string) []string {
	if filepath.Clean(scope) == "." {
		return dedupePaths(scannedPaths)
	}

	merged := make([]string, 0, len(existingPaths)+len(scannedPaths))
	seen := make(map[string]bool, len(existingPaths)+len(scannedPaths))
	prefix := scope + "/"

	for _, p := range existingPaths {
		if p == scope || strings.HasPrefix(p, prefix) {
			continue
		}
		if seen[p] {
			continue
		}
		seen[p] = true
		merged = append(merged, p)
	}

	for _, p := range scannedPaths {
		if seen[p] {
			continue
		}
		seen[p] = true
		merged = append(merged, p)
	}

	return merged
}

func dedupePaths(paths []string) []string {
	out := make([]string, 0, len(paths))
	seen := make(map[string]bool, len(paths))
	for _, p := range paths {
		if seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}

func isPathInScope(path string, scope string) bool {
	if filepath.Clean(scope) == "." {
		return true
	}

	prefix := scope + "/"
	return path == scope || strings.HasPrefix(path, prefix)
}
