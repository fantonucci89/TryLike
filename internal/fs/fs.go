package fs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Folder represents a folder managed by trylike.
type Folder struct {
	Name string
	Path string
}

// ListFolders returns all immediate subdirectories in basePath.
// If basePath doesn't exist it is created.
func ListFolders(basePath string) ([]Folder, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("creating base path: %w", err)
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("reading base path: %w", err)
	}

	var folders []Folder
	for _, e := range entries {
		if e.IsDir() {
			folders = append(folders, Folder{
				Name: e.Name(),
				Path: filepath.Join(basePath, e.Name()),
			})
		}
	}
	return folders, nil
}

// CreateEmpty creates an empty folder with the given name inside basePath.
func CreateEmpty(basePath, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("folder name cannot be empty")
	}
	target := filepath.Join(basePath, name)
	return os.MkdirAll(target, 0755)
}

// CreateFromGitHub clones a GitHub repo into basePath/<name>.
// repoURL must look like a valid git URL.
func CreateFromGitHub(basePath, name, repoURL string) error {
	name = strings.TrimSpace(name)
	repoURL = strings.TrimSpace(repoURL)
	if name == "" {
		return fmt.Errorf("folder name cannot be empty")
	}
	if !isValidGitURL(repoURL) {
		return fmt.Errorf("invalid git URL: %q", repoURL)
	}
	target := filepath.Join(basePath, name)
	cmd := exec.Command("git", "clone", repoURL, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Rename renames oldName to newName inside basePath.
func Rename(basePath, oldName, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return fmt.Errorf("new folder name cannot be empty")
	}
	src := filepath.Join(basePath, oldName)
	dst := filepath.Join(basePath, newName)
	return os.Rename(src, dst)
}

// Move moves a folder to destPath (which must be within $HOME).
// If destPath doesn't exist it is created.
func Move(src, destPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// Expand ~ if present.
	if strings.HasPrefix(destPath, "~/") {
		destPath = filepath.Join(home, destPath[2:])
	}
	// Ensure destination is within $HOME.
	rel, err := filepath.Rel(home, destPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("destination must be within $HOME (%s)", home)
	}
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}
	name := filepath.Base(src)
	dst := filepath.Join(destPath, name)
	return os.Rename(src, dst)
}

// Delete removes a folder and all its contents.
func Delete(path string) error {
	return os.RemoveAll(path)
}

// isValidGitURL does a basic sanity check on a git URL.
func isValidGitURL(url string) bool {
	if url == "" {
		return false
	}
	// Accept https://, git://, git@, ssh://, or plain paths.
	prefixes := []string{
		"https://", "http://", "git://", "git@", "ssh://",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(url, p) {
			return true
		}
	}
	return false
}
