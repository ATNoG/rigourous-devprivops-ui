package fs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// var Repo string

// Clones the default master branch into another directory.
//
// `newRepoPath`: Path to the new cloned repository
//
// returns: The output of the git commands and an error if executing the commands fails.
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

// Adds all changes to the git staging area
//
// `user`: The owner of the repository
//
// returns: The output of the git commands and an error if executing the commands fails.
func AddAll(user string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "add", ".").Output()

	return string(out), err
}

// Commits all changes to the git repository
//
// `user`: The owner of the repository
//
// returns: The output of the git commands and an error if executing the commands fails.
func Commit(user string, message string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "commit", "-m", message).Output()

	return string(out), err
}

// Commits all changes to the git repository
//
// `user`: The owner of the repository
//
// returns: The output of the git commands and an error if executing the commands fails.
func Push(user string) (string, error) {
	repoPath := fmt.Sprintf("%s/%s", LocalDir, user)
	fmt.Println("/usr/bin/git", "-C", repoPath, "push", "origin", "master")
	out, err := exec.Command("/usr/bin/git", "-C", repoPath, "push", "origin", "master").Output()

	fmt.Printf(">>> %s\n", string(out))
	return string(out), err
}

// Clones the master repository onto the desired path for a user.
//
// `repo`: The path where to clone the master repository
//
// `user`: The owner of the repository
//
// `email`: The user's email
//
// returns: The output of the git commands and an error if executing the commands fails.
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

// Finds all conflicts between the current repository and the master.
//
// `repo`: The path to the current repository
//
// returns: The output of the git commands and an error if executing the commands fails.
func GetConflicts(repo string) ([]string, error) {
	out, err := exec.Command("git", "-C", repo, "pull", "origin", "master", "--no-rebase").Output()
	if err != nil {
		fmt.Println(string(out))
		fmt.Printf("Could not pull origin: %s\n", err)
		if err.Error() != "exit status 1" {
			return []string{}, err
		}
	} else {
		fmt.Println(string(out))
	}

	out, err = exec.Command("git", "-C", repo, "status").Output()
	if err != nil {
		fmt.Println("Error reading git status output:", err)
		return []string{}, err
	}
	fmt.Println(string(out))

	conflictFiles := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "both modified:") {
			file := strings.TrimSpace(strings.TrimPrefix(line, "	both modified:   "))
			conflictFiles = append(conflictFiles, file)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading git status output:", err)
		return []string{}, err
	}

	out, err = exec.Command("git", "-C", repo, "merge", "--abort").Output()
	fmt.Println(string(out))
	fmt.Println(err)

	return conflictFiles, err
}
