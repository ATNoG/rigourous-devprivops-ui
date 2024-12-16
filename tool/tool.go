// Package to interact with PrivGuide vis the operating system's shell
package tool

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
)

var Username string // Fuseki username
var Password string // Fuseki password
var DBIP string     // Fuseki IP
var DBPort int      // Fuseki port
var Dataset string  // Fuseki dataset to use

// Run analysis on the user's repository
//
// `reportEndpoint`: The endpoint of the report visualizer, or "" to not use any
//
// `user`: The user whose repository is to be analysed
//
// returns: The PrivGuide output and an error if command execution failed.
func Analyse(reportEndpoint string, user string) (string, error) {
	command := []string{
		"analyse",
		Username, Password, DBIP, strconv.Itoa(DBPort), Dataset,
	}

	if fs.GlobalDir != "" {
		command = append(command, "--global-dir", fs.GlobalDir)
	}

	if fs.LocalDir != "" {
		localDir := fmt.Sprintf("%s/%s", fs.LocalDir, user)
		command = append(command, "--local-dir", localDir)
	}

	if reportEndpoint != "" {
		command = append(command, "--report-endpoint", reportEndpoint)
	}

	fmt.Println(command)
	out, err := exec.Command("./devprivops", command...).Output()

	return string(out), err
}

// Run the tests on the user's repository
//
// `reportEndpoint`: The endpoint of the report visualizer, or "" to not use any
//
// `user`: The user whose repository is to be analysed
//
// returns: The PrivGuide output and an error if command execution failed.
func Test(user string) (string, error) {
	command := []string{
		"test",
		Username, Password, DBIP, strconv.Itoa(DBPort), Dataset,
	}

	if fs.GlobalDir != "" {
		command = append(command, "--global-dir", fs.GlobalDir)
	}

	if fs.LocalDir != "" {
		localDir := fmt.Sprintf("%s/%s", fs.LocalDir, user)
		command = append(command, "--local-dir", localDir)
	}

	fmt.Println(command)
	out, err := exec.Command("./devprivops", command...).Output()

	return string(out), err
}
