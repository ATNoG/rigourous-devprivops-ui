package tool

import (
	"os/exec"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
)

var Username string
var Password string
var DBIP string
var DBPort int
var Dataset string

func Analyse(
	/*
		username string,
		password string,
		databaseIp string,
		databasePort int,
		dataset string,
	*/
	/*
		globalDir string,
		localDir string,
	*/
	reportEndpoint string,
) (string, error) {
	command := []string{
		"analyse",
		Username, Password, DBIP, strconv.Itoa(DBPort), Dataset,
	}

	if fs.GlobalDir != "" {
		command = append(command, "--global-dir", fs.GlobalDir)
	}

	if fs.LocalDir != "" {
		command = append(command, "--local-dir", fs.LocalDir)
	}

	if reportEndpoint != "" {
		command = append(command, "--report-endpoint", reportEndpoint)
	}

	out, err := exec.Command("./devprivops", command...).Output()

	return string(out), err
}

func Test(
/*
	username string,
	password string,
	databaseIp string,
	databasePort int,
	dataset string,
*/
/*
	globalDir string,
	localDir string,
*/

) (string, error) {
	command := []string{
		"test",
		Username, Password, DBIP, strconv.Itoa(DBPort), Dataset,
	}

	if fs.GlobalDir != "" {
		command = append(command, "--global-dir", fs.GlobalDir)
	}

	if fs.LocalDir != "" {
		command = append(command, "--local-dir", fs.LocalDir)
	}

	out, err := exec.Command("./devprivops", command...).Output()

	return string(out), err
}
