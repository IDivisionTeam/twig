use anyhow::Result;
use octocrab::Octocrab;
use tokio::runtime::Builder;

pub struct RemoteChangesParams {
    pub owner: String,
    pub repo: String,
    pub branch: String,
    pub default_branch: String,
    pub title: String,
    pub labels: Vec<String>,
}

pub trait VCSClient {
    fn create_changes_on_remote(&self, params: RemoteChangesParams) -> Result<String>;
}

pub struct DummyClient;

impl DummyClient {
    pub fn new() -> Self {
        Self
    }
}

impl VCSClient for DummyClient {
    fn create_changes_on_remote(&self, _params: RemoteChangesParams) -> Result<String> {
        Ok("remote config is not supplied".to_string())
    }
}

pub struct GithubClient {
    client: Octocrab,
    rt: tokio::runtime::Runtime,
}

impl GithubClient {
    pub fn new(token: String, host: String) -> Result<Self> {
        let rt = Builder::new_current_thread().enable_all().build()?;

        let client =
            rt.block_on(async { Octocrab::builder().base_uri(host)?.personal_token(token).build() })?;
        Ok(Self { client, rt })
    }
}

impl VCSClient for GithubClient {
    fn create_changes_on_remote(&self, params: RemoteChangesParams) -> Result<String> {
        let url = self.rt.block_on(async {
            let pr = self
                .client
                .pulls(&params.owner, &params.repo)
                .create(&params.title, &params.branch, &params.default_branch)
                .send()
                .await?;

            if !params.labels.is_empty() {
                self.client
                    .issues(&params.owner, &params.repo)
                    .add_labels(pr.number, &params.labels)
                    .await?;
            }

            Ok::<reqwest::Url, anyhow::Error>(pr.html_url.unwrap())
        })?;
        Ok(url.to_string())
    }
}
