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
    #[error("failed to checkout branch")]
    CheckoutBranch(#[source] Box<GitError>),
    #[error("failed to check if branch exists")]
    CheckBranchIfExists(#[source] Box<GitError>),
    #[error("failed to push branch to remote")]
    Push(#[source] Box<GitError>),
    #[error("current branch has uncommitted changes")]
    BranchNotClean,
}

type Result<T> = std::result::Result<T, GitError>;

pub fn checkout(branch_name: &str) -> Result<String> {
    let mut args = vec!["checkout"];

    if !branch_exists(branch_name)? {
        debug!("Branch {branch_name} is new, adding '-b' flag");
        args.push("-b");
    }

    args.push(branch_name);
    execute(args).map_err(|err| GitError::CheckoutBranch(Box::new(err)))
}

pub fn push_to_remote(branch_name: &str, remote: &str) -> Result<String> {
    execute(vec!["push", "-u", remote, branch_name]).map_err(|err| GitError::Push(Box::new(err)))
}

pub fn branch_exists(branch_name: &str) -> Result<bool> {
    match execute(vec![
        "show-ref",
        "--exists",
        &format!("refs/heads/{branch_name}"),
    ]) {
        Ok(_) => Ok(true),
        Err(GitError::Failed(ref msg)) if msg.contains("reference does not exist") => Ok(false),
        Err(GitError::Failed(ref msg)) if msg.contains("unknown option `exists'") => branch_exists_fallback(branch_name),
        Err(e) => Err(GitError::CheckBranchIfExists(Box::new(e))),
    }
}

/// branch_exists_fallback is a fallback for --exists parameter for git version < v2.44.0
fn branch_exists_fallback(branch_name: &str) -> Result<bool> {
    match execute(vec![
        "rev-parse",
        &format!("refs/heads/{branch_name}"),
    ]) {
        Ok(_) => Ok(true),
        Err(GitError::Failed(ref msg)) if msg.contains("unknown revision or path not in the working tree") => Ok(false),
        Err(e) => Err(GitError::CheckBranchIfExists(Box::new(e))),
    }
}


pub fn fetch_prune() -> Result<String> {
    execute(vec!["fetch", "-p"])
}

pub fn ensure_clean_status() -> Result<()> {
    let output = execute(vec!["status", "-s"])?;
    if !output.is_empty() {
        return Err(GitError::BranchNotClean);
    }
    Ok(())
}

pub fn get_local_branches() -> Result<String> {
    execute(vec!["branch"])
}

pub fn delete_local_branch(branch_name: &str) -> Result<String> {
    execute(vec!["branch", "-D", branch_name])
}

pub fn delete_remote_branch(remote: &str, branch_name: &str) -> Result<String> {
    execute(vec!["push", "-d", remote, branch_name])
}

pub fn execute(args: Vec<&str>) -> Result<String> {
    let output = Command::new("git").args(args.as_slice()).output()?;

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
