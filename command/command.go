package command

import (
    "fmt"
    "github.com/fatih/color"
    "os/exec"
)

func Checkout(branchName string) (string, error) {
    out, err := exec.Command("git", "checkout", "-b", branchName).Output()

    if err != nil {
        return "", fmt.Errorf("git checkout %w", err)
    }

    coloredOutput := color.GreenString(string(out))
    return coloredOutput, err
}
