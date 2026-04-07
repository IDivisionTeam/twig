use std::{io, process::Command};

use log::debug;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum GitError {
    #[error("failed to execute command")]
    Execute(#[from] io::Error),
    #[error("failed to get output")]
    Output(#[from] std::string::FromUtf8Error),
    #[error("failed to checkout branch")]
    CheckoutBranch,
    #[error("failed to push branch to remote")]
    Push,
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
    execute(args).map_err(|_| GitError::CheckoutBranch)
}

pub fn push_to_remote(branch_name: &str, remote: &str) -> Result<String> {
    execute(vec!["push", "-u", remote, branch_name]).map_err(|_| GitError::Push)
}

pub fn branch_exists(branch_name: &str) -> Result<bool> {
    let output = execute(vec!["show-ref", &format!("refs/heads/{branch_name}")])?;
    Ok(!output.is_empty())
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

    Ok(if output.stdout.is_empty() {
        String::from_utf8(output.stderr)
    } else {
        String::from_utf8(output.stdout)
    }?)
}
