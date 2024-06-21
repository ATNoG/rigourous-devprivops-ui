package fs

import (
	"fmt"
	"os/exec"
)

var Repo string

func Clone(branch string, destination string) (string, error) {
	out, err := exec.Command(
		"/usr/bin/git", "clone", fmt.Sprintf("file://%s", Repo),
		"--single-branch", "--branch", branch,
		"--depth", "1", destination).Output()

	return string(out), err
}

func AddAll(dir string) (string, error) {
	out, err := exec.Command("/usr/bin/git", "-C", dir, "add", ".").Output()

	return string(out), err
}

func Commit(dir string, message string) (string, error) {
	out, err := exec.Command("/usr/bin/git", "-C", dir, "commit", "-m", message).Output()

	return string(out), err
}

func SetupRepo(repo string, destination string, user string, email string) (string, error) {
	res, err := Clone(repo, destination)
	if err != nil {
		return res, err
	}

	out, err := exec.Command("/usr/bin/git", "-C", destination, "config", "--local", "user.name", user).Output()
	if err != nil {
		return string(out), err
	}

	out, err = exec.Command("/usr/bin/git", "-C", destination, "config", "--local", "user.name", user).Output()

	return string(out), err
}
