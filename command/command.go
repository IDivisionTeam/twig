package command

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "slices"

    "github.com/joho/godotenv"
)

func ReadEnvVariables() error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("obtaining home directory: %w", err)
    }

    envPath := filepath.Join(homeDir, ".config", "brcha", ".env")

    err = godotenv.Load(envPath)
    if err != nil {
        return fmt.Errorf("reading .env: %w", err)
    }

    return nil
}

func HasBranch(branchName string) bool {
    err := exec.Command("git", "branch", "--contains", branchName).Run()
    return err == nil
}

func Checkout(branchName string, hasBranch bool) (string, error) {
    args := []string{"checkout", branchName}

    if !hasBranch {
        args = slices.Insert(args, 1, "-b")
    }

    out, err := exec.Command("git", args...).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("git checkout: %s%w", string(out), err)
    }

    return string(out), nil
}
