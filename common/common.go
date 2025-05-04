package common

import (
    "errors"
    "fmt"
    "os/exec"
    "slices"
    "twig/log"
)

func HasBranch(branchName string) bool {
    err := exec.Command("git", "branch", "--contains", branchName).Run()

    doesExist := err == nil
    if doesExist {
        log.Debug().Printf("Branch %q exists locally", branchName)
    } else {
        log.Debug().Printf("Branch %q does not exist locally", branchName)
    }

    return doesExist
}

func Checkout(branchName string, hasBranch bool) (string, error) {
    log.Info().Println(fmt.Sprintf("Checkout to %q", branchName))

    args := []string{"checkout", branchName}

    if !hasBranch {
        log.Debug().Printf(fmt.Sprintf("Branch %q is new, adding '-b' flag", branchName))
        args = slices.Insert(args, 1, "-b")
    }

    out, err := exec.Command("git", args...).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func BranchStatus() error {
    log.Info().Println("Check branch status")

    out, err := exec.Command("git", "status", "-s").CombinedOutput()
    if err != nil {
        return err
    }

    if outputSize := len(string(out)); outputSize > 0 {
        return errors.New("current branch has uncommitted changes")
    }

    return nil
}

func GetLocalBranches() (string, error) {
    log.Info().Println("Get local branches")

    out, err := exec.Command("git", "branch").CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func ExecuteFetchPrune() (string, error) {
    log.Info().Println("Run fetch and prune")

    out, err := exec.Command("git", "fetch", "-p").CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func DeleteLocalBranch(branchName string) (string, error) {
    log.Info().Println(fmt.Sprintf("Delete local branch %q", branchName))

    out, err := exec.Command("git", "branch", "-D", branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func DeleteRemoteBranch(remote string, branchName string) (string, error) {
    log.Info().Println(fmt.Sprintf("Delete remote branch '%s/%s'", remote, branchName))

    out, err := exec.Command("git", "push", "-d", remote, branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func PushToRemote(branchName string, remote string) (string, error) {
    log.Info().Println(fmt.Sprintf("Push branch to remote '%s/%s'", remote, branchName))

    out, err := exec.Command("git", "push", "-u", remote, branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}
