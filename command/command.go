package command

import (
    "brcha/log"
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
        return fmt.Errorf("read env: obtaining home directory: %w", err)
    }
    log.Debug().Printf("read env: home dir = %s", homeDir)

    envPath := filepath.Join(homeDir, ".config", "brcha", ".env")
    log.Debug().Printf("read env: env path = %s", envPath)

    err = godotenv.Load(envPath)
    if err != nil {
        return fmt.Errorf("read env: load: %w", err)
    }

    return nil
}

func HasBranch(branchName string) bool {
    err := exec.Command("git", "branch", "--contains", branchName).Run()

    doesExist := err == nil
    log.Debug().Printf("has branch: %s exists = %t", branchName, doesExist)

    return doesExist
}

func Checkout(branchName string, hasBranch bool) (string, error) {
    args := []string{"checkout", branchName}
    log.Debug().Printf("checkout: args: %s", args)

    if !hasBranch {
        log.Debug().Printf("checkout: %s is a new branch, adding -b flag", branchName)
        args = slices.Insert(args, 1, "-b")
    }

    out, err := exec.Command("git", args...).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("git checkout: %s%w", string(out), err)
    }

    return string(out), nil
}
