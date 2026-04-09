package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// file, err := os.ReadFile(".env")
	// if err != nil {
	// 	panic("shit")
	// }

	// base64file := base64.StdEncoding.EncodeToString(file)

	// commandString := fmt.Sprintf("add-generic-password -a %q -s %q -w %q -U\n", "storkey", "githubrepo", base64file)

	// cmd := exec.Command("security", "-i")

	// cmd.Stdin = strings.NewReader(commandString)

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	panic("shit2")
	// }
	// fmt.Printf("Everything went well, this the command i run: %s", output)
	cmd := exec.Command("security", "find-generic-password", "-a", "storkey", "-s", "githubrepo", "-w")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("shit")
	}
	fmt.Printf("output is: %s", out)
}
