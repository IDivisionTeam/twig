use std::process::Command;

use anyhow::{Context, Result};

pub fn checkout(branch_name: &str) -> Result<String> {
    if !branch_exists(branch_name)?.is_empty() {
        return Ok("Branch already exists".to_string());
    }
    execute("git", vec!["checkout", "-b", branch_name]).context("failed to checkout branch")
}

pub fn push_to_remote(branch_name: &str, remote: &str) -> Result<String> {
    execute("git", vec!["push", "-u", remote, branch_name]).context("failed to push")
}

pub fn branch_exists(branch_name: &str) -> Result<String> {
    execute(
        "git",
        vec!["show-ref", &format!("refs/heads/{branch_name}")],
    )
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
