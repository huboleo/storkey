package internal

import "fmt"

func Account(repo string, branch string) string {
	return fmt.Sprintf("storkey:%s:%s", repo, branch)
}

func Service(path string) string {
	return fmt.Sprintf("storkey:%s", path)
}

func Index() string {
	return "storkey:index:paths"
}
