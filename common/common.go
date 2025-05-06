package common

import (
    "errors"
    "fmt"
    "slices"
    "twig/git"
    "twig/log"
)

func HasGit() bool {
    return git.Command(git.Version).Run() == nil
}

func HasBranch(branchName string) bool {
    err := git.Command(git.Branch, "--contains", branchName).Run()

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

    args := []string{branchName}

    if !hasBranch {
        log.Debug().Printf(fmt.Sprintf("Branch %q is new, adding '-b' flag", branchName))
        args = slices.Insert(args, 0, "-b")
    }

    out, err := git.Command(git.Checkout, args...).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func BranchStatus() error {
    log.Info().Println("Check branch status")

    out, err := git.Command(git.Status, "-s").CombinedOutput()
    if err != nil {
        return err
    }

    if outputSize := len(string(out)); outputSize > 1 {
        return errors.New("current branch has uncommitted changes")
    }

    return nil
}

func GetLocalBranches() (string, error) {
    log.Info().Println("Get local branches")

    out, err := git.Command(git.Branch).CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func ExecuteFetchPrune() (string, error) {
    log.Info().Println("Run fetch and prune")

    out, err := git.Command(git.Fetch, "-p").CombinedOutput()
    if err != nil {
        return "", err
    }

    return string(out), nil
}

func DeleteLocalBranch(branchName string) (string, error) {
    log.Info().Println(fmt.Sprintf("Delete local branch %q", branchName))

    out, err := git.Command(git.Branch, "-D", branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func DeleteRemoteBranch(remote string, branchName string) (string, error) {
    log.Info().Println(fmt.Sprintf("Delete remote branch '%s/%s'", remote, branchName))

    out, err := git.Command(git.Push, "-d", remote, branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}

func PushToRemote(branchName string, remote string) (string, error) {
    log.Info().Println(fmt.Sprintf("Push branch to remote '%s/%s'", remote, branchName))

    out, err := git.Command(git.Push, "-u", remote, branchName).CombinedOutput()
    if err != nil {
        return string(out), err
    }

    return string(out), nil
}
