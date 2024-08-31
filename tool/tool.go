package tool

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
)

var Username string
var Password string
var DBIP string
var DBPort int
var Dataset string

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
