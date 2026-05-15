use std::{io, process::Command};

use log::debug;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum GitError {
    #[error("failed to execute git command")]
    Execute(#[from] io::Error),
    #[error("git command failed: {0}")]
    Failed(String),
    #[error("failed to get output")]
    Output(#[from] std::string::FromUtf8Error),
    #[error("failed to get current branch")]
    CurrentBranch(#[source] Box<GitError>),
    #[error("failed to checkout branch")]
    CheckoutBranch(#[source] Box<GitError>),
    #[error("failed to check if branch exists")]
    CheckBranchIfExists(#[source] Box<GitError>),
    #[error("failed to push branch to remote")]
    Push(#[source] Box<GitError>),
    #[error("current branch has uncommitted changes")]
    BranchNotClean,
    #[error("failed to get latest commit meesage")]
    LatestCommitMessage(#[source] Box<GitError>),
    #[error("failed to get remote branch")]
    GetRemoteBranch(#[source] Box<GitError>),
}

type Result<T> = std::result::Result<T, GitError>;

pub fn checkout(branch_name: &str) -> Result<String> {
    let mut args = vec!["checkout"];

    if !branch_exists(branch_name)? {
        debug!("Branch {branch_name} is new, adding '-b' flag");
        args.push("-b");
    }

    args.push(branch_name);
    execute(&args).map_err(|err| GitError::CheckoutBranch(Box::new(err)))
}

pub fn push_to_remote(branch_name: &str, remote: &str) -> Result<String> {
    execute(&["push", "-u", remote, branch_name]).map_err(|err| GitError::Push(Box::new(err)))
}

pub fn current_branch() -> Result<String> {
    execute(&["branch", "--show-current"])
        .map(|branch_name| branch_name.trim().to_string())
        .map_err(|err| GitError::CurrentBranch(Box::new(err)))
}

pub fn branch_exists(branch_name: &str) -> Result<bool> {
    match execute(&["show-ref", "--exists", &format!("refs/heads/{branch_name}")]) {
        Ok(_) => Ok(true),
        Err(GitError::Failed(ref msg)) if msg.contains("reference does not exist") => Ok(false),
        Err(GitError::Failed(ref msg)) if msg.contains("unknown option `exists'") => {
            branch_exists_fallback(branch_name)
        }
        Err(e) => Err(GitError::CheckBranchIfExists(Box::new(e))),
    }
}

/// `branch_exists_fallback` is a fallback for --exists parameter for git version < v2.44.0
fn branch_exists_fallback(branch_name: &str) -> Result<bool> {
    match execute(&["rev-parse", &format!("refs/heads/{branch_name}")]) {
        Ok(_) => Ok(true),
        Err(GitError::Failed(ref msg))
            if msg.contains("unknown revision or path not in the working tree") =>
        {
            Ok(false)
        }
        Err(e) => Err(GitError::CheckBranchIfExists(Box::new(e))),
    }
}

pub fn fetch_prune() -> Result<String> {
    execute(&["fetch", "-p"])
}

pub fn ensure_clean_status() -> Result<()> {
    let output = execute(&["status", "-s"])?;
    if !output.is_empty() {
        return Err(GitError::BranchNotClean);
    }
    Ok(())
}

pub fn get_local_branches() -> Result<String> {
    execute(&["branch"])
}

pub fn delete_local_branch(branch_name: &str) -> Result<String> {
    execute(&["branch", "-D", branch_name])
}

pub fn delete_remote_branch(remote: &str, branch_name: &str) -> Result<String> {
    execute(&["push", "-d", remote, branch_name])
}

pub fn latest_commit_msg() -> Result<String> {
    execute(&["log", "-1", "--pretty=%B"])
        .map(|commit_msg| commit_msg.trim().to_string())
        .map_err(|err| GitError::LatestCommitMessage(Box::new(err)))
}

pub fn get_origin_url() -> Result<String> {
    execute(&["remote", "get-url", "origin"])
}

pub fn get_remote_branch(local_branch_name: &str) -> Result<Option<String>> {
    match execute(&[
        "rev-parse",
        "--abbrev-ref",
        &format!("{local_branch_name}@{{upstream}}"),
    ]) {
        Ok(remote_branch) => Ok(Some(remote_branch)),
        Err(GitError::Failed(ref msg)) if msg.contains("no upstream configured for branch") => Ok(None),
        Err(e) => Err(GitError::GetRemoteBranch(Box::new(e))),
    }
}

pub fn execute(args: &[&str]) -> Result<String> {
    let output = Command::new("git").args(args).output()?;

    let stdout = String::from_utf8(output.stdout)?;
    let stderr = String::from_utf8(output.stderr)?;

    if output.status.success() {
        Ok(if stdout.is_empty() { stderr } else { stdout })
    } else {
        let error_msg = if !stderr.is_empty() {
            stderr
        } else if !stdout.is_empty() {
            stdout
        } else {
            format!("exit code: {:?}", output.status.code())
        };

        Err(GitError::Failed(error_msg))
    }
}
