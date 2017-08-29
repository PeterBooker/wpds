package main

import (
	"fmt"
	"os/exec"
)

func startSearch(pattern string) string {

	out, err := exec.Command("bundle/pt.exe", pattern).Output()
	if err != nil {
		fmt.Println(err)
	}

	// Convert bytes to string
	results := fmt.Sprintf("%s", out)

	return results

}

func outputResults(results string, pattern string, output string) {

	switch output {
	case "file":
		fmt.Println(results)
	default:
		fmt.Println(results)
	}

}
