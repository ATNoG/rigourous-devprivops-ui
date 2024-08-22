package fs

import (
	"fmt"
	"os"
	"os/exec"
)

// var Repo string

func Clone(newRepoPath string) (string, error) {
	destPath := fmt.Sprintf("%s/%s", LocalDir, newRepoPath)
	exists, err := exists(destPath)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	if !exists {
		fmt.Printf("Cloning '%s' into '%s'\n", fmt.Sprintf("file://%s/master", LocalDir), destPath)
		os.Mkdir(destPath, 0777)

		out, err := exec.Command(
			// "/usr/bin/git", "clone", fmt.Sprintf("file://%s", Repo),
			"/usr/bin/git",
			"-C", destPath,
			"clone", fmt.Sprintf("file://%s/master", LocalDir), destPath,
			// "--single-branch", "--branch", branch,
			// "--depth", "1", destination,
		).Output()

		fmt.Println(string(out))
		fmt.Println(err)
		return string(out), err
	}
	return "repo already exists", nil
}

func AddAll(dir string, user string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "add", ".").Output()

	return string(out), err
}

func Commit(dir string, user string, message string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "commit", "-m", message).Output()

	return string(out), err
}

func Push(dir string, user string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "push", "origin", "master").Output()

	return string(out), err
}

func SetupRepo(repo string, user string, email string) (string, error) {
	res, err := Clone(repo)
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	repoPath := fmt.Sprintf("%s/%s", LocalDir, repo)
	fmt.Printf("On repo '%s'\n", repoPath)

	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "config", "--local", "user.name", user).Output()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
		return string(out), err
	}
	fmt.Println(string(out))

	out, err = exec.Command("/usr/bin/git", "-C", repoPath, "config", "--local", "user.email", email).Output()
	fmt.Println(string(out))

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
	}
	return string(out), err
}
