package usecase

import (
	"context"
	"ez-snapshot/internal/repository/storage"
	"fmt"
	"os/exec"
	"strings"
)

// DependencyChecker checks required CLI tools are available.
type DependencyChecker struct {
	Dependencies []string
	Storage      storage.Repository
}

// NewDependencyChecker returns a checker for mysql + rclone.
func NewDependencyChecker(
	s storage.Repository,
) *DependencyChecker {
	return &DependencyChecker{
		Dependencies: []string{
			"mysql",     // MySQL client
			"mysqldump", // for backup
			"rclone",    // for remote storage
		},
		Storage: s,
	}
}

// Check runs "command --version" for each dependency.
func (dc *DependencyChecker) Check() error {
	fmt.Println("Checking dependencies...")

	for _, dep := range dc.Dependencies {
		if err := checkCommand(dep); err != nil {
			return fmt.Errorf("dependency check failed: %s: %w", dep, err)
		}
	}

	fmt.Println("Checking rclone connectivity...")

	_, err := dc.Storage.List(context.Background())
	if err != nil {
		return fmt.Errorf("rclone rc API is not running, run it first using rclone rcd --rc-no-auth --rc-addr=:5572\n")
	}

	fmt.Println("✅ rclone server is running ")

	return nil
}

// checkCommand ensures binary exists and is executable.
func checkCommand(name string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("not found in PATH")
	}

	// run "<name> --version"
	cmd := exec.Command(path, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("❌ failed to run --version: %w", err)
	}

	if !strings.Contains(strings.ToLower(string(output)), strings.ToLower(name)) {
		return fmt.Errorf("❌ unexpected output: %s", strings.TrimSpace(string(output)))
	}

	fmt.Printf("✅ %s are installed \n", name)

	return nil
}
