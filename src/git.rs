use std::process::Command;

use anyhow::{Context, Result, bail};
use log::debug;

pub fn checkout(branch_name: &str) -> Result<String> {
    let mut args = vec!["checkout"];

    if !branch_exists(branch_name)? {
        debug!("Branch {branch_name} is new, adding '-b' flag");
        args.push("-b");
    }

    args.push(branch_name);
    execute("git", args).context("failed to checkout branch")
}

pub fn push_to_remote(branch_name: &str, remote: &str) -> Result<String> {
    execute("git", vec!["push", "-u", remote, branch_name]).context("failed to push")
}

pub fn branch_exists(branch_name: &str) -> Result<bool> {
    let output = execute(
        "git",
        vec!["show-ref", &format!("refs/heads/{branch_name}")],
    )?;
    Ok(!output.is_empty())
}

pub fn fetch_prune() -> Result<String> {
    execute("git", vec!["fetch", "-p"])
}

pub fn branch_status() -> Result<()> {
    let output = execute("git", vec!["status", "-s"])?;
    if !output.is_empty() {
        bail!("current branch has uncommitted changes");
    }
    Ok(())
}

pub fn get_local_branches() -> Result<String> {
    execute("git", vec!["branch"])
}

pub fn delete_local_branch(branch_name: &str) -> Result<String> {
    execute("git", vec!["branch", "-D", branch_name])
}

pub fn delete_remote_branch(remote: &str, branch_name: &str) -> Result<String> {
    execute("git", vec!["push", "-d", remote, branch_name])
}

fn execute(cmd: &str, args: Vec<&str>) -> Result<String> {
    let output = Command::new(cmd)
        .args(args.as_slice())
        .output()
        .context("failed to execute command")?;

    if output.stdout.is_empty() {
        String::from_utf8(output.stderr).context("failed to get output")
    } else {
        String::from_utf8(output.stdout).context("failed to get output")
    }
}
