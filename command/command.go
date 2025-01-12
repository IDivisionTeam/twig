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

    log.Info().Println("environment loaded")
    return nil
}

func HasBranch(branchName string) bool {
    err := exec.Command("git", "branch", "--contains", branchName).Run()

    doesExist := err == nil
    log.Debug().Printf("%s exists locally = %t", branchName, doesExist)

    return doesExist
}

func Checkout(branchName string, hasBranch bool) (string, error) {
    args := []string{"checkout", branchName}
    log.Info().Println("executing 'git checkout'")
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

func BranchStatus() error {
    log.Info().Println("executing 'git status'")

    out, err := exec.Command("git", "status", "-s").CombinedOutput()
    if err != nil {
        return err
    }

    if outputSize := len(string(out)); outputSize > 0 {
        return fmt.Errorf("git status: current branch has uncommitted changes")
    }

    return nil
}

func GetLocalBranches() (string, error) {
    log.Info().Println("executing 'git branch'")

    out, err := exec.Command("git", "branch").CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func ExecuteFetchPrune() (string, error) {
    log.Info().Println("executing 'git fetch + prune'")

    out, err := exec.Command("git", "fetch", "-p").CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func DeleteLocalBranch(branchName string) (string, error) {
    log.Info().Printf("executing 'git branch local delete' %s", branchName)

    out, err := exec.Command("git", "branch", "-D", branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func DeleteRemoteBranch(remote string, branchName string) (string, error) {
    log.Info().Printf("executing 'git branch remote delete' %s/%s", remote, branchName)

    out, err := exec.Command("git", "push", "-d", remote, branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}
