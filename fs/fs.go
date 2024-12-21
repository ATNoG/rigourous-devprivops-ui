// Package to abstract file system accesses,
// namely by handling lookup from both the global and local directories
//
// By default, the local path is `.devprivops/` and the global path is `/etc/devprivops/`.
// Files in the local path override those in the global directory.
//
// The unexported functions are independent of the local and global directories and are made to increase
// testability. These are the ones that should be targeted in unit tests and thus are exported in `export_test.go`.
//
// Higher level actions use git repositories to ensure concurrency.
// There is an innaccessible local repository, the `master` repository, with which users cannot interact.
// Each user is given a repository with their name, cloned from the `master`.
// When actions are completed, users can merge their work onto the master repository, making it accessible to the other users.
//
// This package only supports UNIX paths
package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"

	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/templates"
	"github.com/Joao-Felisberto/devprivops-ui/util"
	"gopkg.in/yaml.v3"
)

/*
	lookup order:
	1. /etc/appName
	2. .appName/
*/

var (
	LocalDir       = fmt.Sprintf("./.%s", util.AppName)   // The local directory
	GlobalDir      = fmt.Sprintf("/etc/%s", util.AppName) // The global directory
	m              sync.Mutex
	SessionManager *sessionmanament.SessionManager = sessionmanament.GetSessionManager()
)

// Finds the directory where the user's repository is located
//
// `sessionKey`: The user name
//
// returns: The expected path to the repository and a boolean denoting whether it exists
func GetBranch(sessionKey string) (string, bool) {
	res := fmt.Sprintf("%s/%s", LocalDir, sessionKey)
	_, err := os.Stat(res)
	if err != nil {
		return "", false
	}

	return sessionKey, true
}

// Returns the path of a file relative to the local or global root using the pre-determined paths to the local and global directories
//
// `relativePath`: the path relative to either root
//
// returns: the path to the provided file relative to the root it is in, or an error if reading any of the directories fails.
func GetFile(relativePath string, user string) (string, error) {
	branch, ok := GetBranch(user)
	fmt.Println("Branch ", branch)

	if !ok {
		return "", fmt.Errorf("could not find %s's branch", user)
	}

	fmt.Printf("rel: '%s' local: '%s', global: '%s'\n", relativePath, fmt.Sprintf("%s/%s", LocalDir, branch), GlobalDir)
	return getFile(
		relativePath,
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Returns the path of a file relative to the local or global root using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the path to the provided file relative to the root it is in, or an error if reading any of the directories fails.
func getFile(relativePath string, localRoot string, globalRoot string) (string, error) {
	localPath := fmt.Sprintf("%s/%s", localRoot, relativePath)
	if _, err := os.Stat(localPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Does not exist", localPath)
		defaultPath := fmt.Sprintf("%s/%s", globalRoot, relativePath)
		if _, err := os.Stat(defaultPath); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Does not exist", defaultPath)
			return "", err
		}
		return defaultPath, nil
	}
	return localPath, nil
}

// Returns the paths of the system descriptions relative to their respective root using the default paths to the local and global directories
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions, or an error if reading any of the directories fails.
func GetDescriptions(descriptionRoot string, user string) ([]string, error) {
	branch, ok := GetBranch(user)

	if !ok {
		return []string{}, fmt.Errorf("could not find %s's branch", user)
	}

	return getDescriptions(
		descriptionRoot,
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Returns the paths of the system descriptions relative to their respective root using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions, or an error if reading any of the directories fails.
func getDescriptions(descriptionRoot string, localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/%s/", localRoot, descriptionRoot)
	globalPath := fmt.Sprintf("%s/%s/", globalRoot, descriptionRoot)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error reading local directory: %s", err)
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("%s/%s", descriptionRoot, e.Name()))
	}

	entries, err = os.ReadDir(globalPath)
	if err != nil {
		return files, nil
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("%s/%s", descriptionRoot, e.Name()))
	}

	return files, nil
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func GetRegulations(user string) ([]string, error) {
	branch, ok := GetBranch(user)

	if !ok {
		return []string{}, fmt.Errorf("could not find %s's branch", user)
	}

	return getRegulations(
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func getRegulations(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/regulations/", localRoot)
	defaultPath := fmt.Sprintf("%s/regulations/", globalRoot)

	files := []string{}

	localRegulations, err := getDirsInDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	defaultRegulations, err := getDirsInDir(defaultPath)
	if err != nil {
		files = append(files, localRegulations...)

		return files, nil
	}

	files = append(files, localRegulations...)
	files = append(files, defaultRegulations...)

	return files, nil
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func GetTests(user string) ([]string, error) {
	branch, ok := GetBranch(user)

	if !ok {
		return []string{}, fmt.Errorf("could not find %s's branch", user)
	}

	return getTests(
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func getTests(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/tests/", localRoot)
	defaultPath := fmt.Sprintf("%s/tests/", globalRoot)

	files := []string{}

	localRegulations, err := getDirsInDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	defaultRegulations, err := getDirsInDir(defaultPath)
	if err != nil {
		files = append(files, localRegulations...)

		return files, nil
	}

	files = append(files, localRegulations...)
	files = append(files, defaultRegulations...)

	return files, nil
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func GetTestScenarios(scenario string, user string) ([]string, error) {
	fmt.Printf("Tests in %s for %s\n", scenario, user)
	branch, ok := GetBranch(user)

	if !ok {
		return []string{}, fmt.Errorf("could not find %s's branch", user)
	}

	return getTestScenarios(
		scenario,
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func getTestScenarios(scenario string, localRoot string, globalRoot string) ([]string, error) {
	relativePath := fmt.Sprintf("tests/%s/", scenario)
	localPath := fmt.Sprintf("%s/%s/", localRoot, relativePath)
	defaultPath := fmt.Sprintf("%s/%s/", globalRoot, relativePath)

	files := []string{}

	localRegulations, err := os.ReadDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	defaultRegulations, err := os.ReadDir(defaultPath)
	if err != nil {
		files = append(files, util.Map(localRegulations, func(d fs.DirEntry) string {
			return fmt.Sprintf("%s%s", relativePath, d.Name())
		})...)

		return files, nil
	}

	files = append(files, util.Map(localRegulations, func(d fs.DirEntry) string {
		return fmt.Sprintf("%s%s", relativePath, d.Name())
	})...)
	files = append(files, util.Map(defaultRegulations, func(d fs.DirEntry) string {
		return fmt.Sprintf("%s%s", relativePath, d.Name())
	})...)

	return files, nil
}

// Find all top level directories inside a directory
//
// `path`: The parent directory of which we want to know the subdirectories
//
// returns: The list of subdirectories, or an error if reading any of the directories fails.
func getDirsInDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return util.Map(
		util.Filter(entries, func(de fs.DirEntry) bool { return de.IsDir() }),
		func(de fs.DirEntry) string { return de.Name() },
	), nil
}

// Find the relative directories of each configuration file.
// The returned directories contain the root
//
// returns: The list of configuration files, or an error if reading any of the directories fails.
func GetConfigs(user string) ([]string, error) {
	branch, ok := GetBranch(user)

	if !ok {
		return []string{}, fmt.Errorf("could not find %s's branch", user)
	}

	return getConfigs(
		fmt.Sprintf("%s/%s", LocalDir, branch),
		GlobalDir,
	)
}

// Find the relative directories of each configuration file.
// The returned directories contain the root
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: The list of configuration files, or an error if reading any of the directories fails.
func getConfigs(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/config/", localRoot)
	globalPath := fmt.Sprintf("%s/config/", globalRoot)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error reading local directory: %s", err)
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("config/%s", e.Name()))
	}

	entries, err = os.ReadDir(globalPath)
	if err != nil {
		return files, nil
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("config/%s", e.Name()))
	}

	return files, nil
}

// Sets the description of an attack node in an attack tree provided its query file.
//
// `node`: The root node of the attack tree
//
// `queryFile`: The file containing he query of the node whose description we want to change, acting as its unique identifier
//
// `newDescription`: The new description to replace the old one with
//
// returns: `true` if there was a description change, `false` otherwise
func ChangeTreeDescription(node *templates.TreeNode, queryFile string, newDescription string) bool {
	fmt.Printf("Comp '%s' <=> '%s'\n", node.Query, queryFile)
	if node.Query == queryFile {
		fmt.Println("FOUND")
		node.Description = newDescription
		return true
	} else {
		for _, child := range node.Children {
			if ChangeTreeDescription(child, queryFile, newDescription) {
				return true
			}
		}
		return false
	}
}

// Saves an attack tree's metadata to a file
//
// `tree`: The attack tree's root node
//
// `file`: The file path where to store it
//
// returns: An error if marshaling to YAML or writing the file fails
func SaveTreeDescription(tree *templates.TreeNode, file string) error {
	data, err := yaml.Marshal(tree)
	if err != nil {
		return err
	}

	err = WriteFileSync(file, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// Write files synchronously to avoid corruption
//
// `file`: The path tot he file to write
//
// `data`: The contents of the file
//
// `permissions`: The permissions for the file
//
// returns: An error if writing the file fails
func WriteFileSync(file string, data []byte, permissions fs.FileMode) error {

	m.Lock()
	err := os.WriteFile(file, data, 0666)
	m.Unlock()

	return err
}

// Check if a file exists handling errors
//
// `path`: The path to check
//
// returns: true if the file exists and an error if there was any
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
