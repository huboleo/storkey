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

// func Pull(repo string) error {
// 	filesMap := make(map[string]string)

// }

func Add() error {
	repo, err := RepoAddress()
	if err != nil {
		return err
	}
	branch, err := Branch()
	if err != nil {
		return err
	}

	account := Account(repo, branch)
	indexFile := Index()
	var paths []string

	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".env") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			b64 := base64.StdEncoding.EncodeToString(b)
			rel, err := filepath.Rel(".", path)
			if err != nil {
				return err
			}
			service := Service(rel)
			paths = append(paths, service)
			cmdString := fmt.Sprintf("add-generic-password -a %q -s %q -w %q -U\n", account, service, b64)
			cmd := exec.Command("security", "-i")
			cmd.Stdin = strings.NewReader(cmdString)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to save env file %q: %w (output: %s)", rel, err, strings.TrimSpace(string(out)))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	idxJSON, err := json.Marshal(paths)
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

	return nil
}
