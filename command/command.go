package command

import (
    "fmt"
    "github.com/fatih/color"
    "os/exec"

    "github.com/joho/godotenv"
)

func ReadEnvVariables() error {
    err := godotenv.Load(".env")
    if err != nil {
        return fmt.Errorf("read .env %w", err)
    }

    return nil
}

func Checkout(branchName string) (string, error) {
    out, err := exec.Command("git", "checkout", "-b", branchName).Output()
    if err != nil {
        return "", fmt.Errorf("git checkout %w", err)
    }

    coloredOutput := color.GreenString(string(out))
    return coloredOutput, nil
}
