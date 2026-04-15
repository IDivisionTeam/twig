use anyhow::Result;

pub trait VCSClient {
    fn create_changes_on_remote(&self) -> Result<()>;
}

pub struct GithubClient {

}

impl GithubClient {
    pub fn new() -> Self {
        Self{}
    }
}

impl VCSClient for GithubClient {
    fn create_changes_on_remote(&self) -> Result<()> {
        todo!()
    }
}