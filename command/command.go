package command

import (
    "fmt"
    "github.com/fatih/color"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/joho/godotenv"
)

func ReadEnvVariables() error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
       return fmt.Errorf("getting home directory: %w", err)
    }

    envPath := filepath.Join(homeDir, ".config", "brcha", ".env")

    err = godotenv.Load(envPath)
    if err != nil {
        return fmt.Errorf("read .env %w", err)
    }

    return nil
}

func Checkout(branchName string) (string, error) {
    out, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("git checkout %w", err)
    }

    coloredOutput := color.GreenString(string(out))
    return coloredOutput, nil
}
