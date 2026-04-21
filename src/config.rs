use figment::value::Value;
use figment::{
    Figment,
    providers::{Format, Toml},
};
use serde::{Deserialize, Deserializer, Serialize};
use std::collections::BTreeMap;
#[cfg(not(test))]
use std::env;
use std::fmt::Display;
use std::fmt::Write as fmtWrite;
use std::{
    collections::HashMap,
    fmt,
    fs::{self, File},
    io::{self, Write},
    path::Path,
};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigError {
    #[error("failed to extract config")]
    Extract(#[from] Box<figment::Error>),
    #[error("failed to parse config")]
    Parse(#[from] toml::ser::Error),
    #[error("failed to create file")]
    File(#[from] io::Error),
    #[error("failed to get config file")]
    MissingConfig,
}

impl From<figment::Error> for ConfigError {
    fn from(err: figment::Error) -> Self {
        ConfigError::Extract(Box::new(err))
    }
}

#[derive(Debug, Default, PartialEq, Serialize, Deserialize)]
pub struct Config {
    pub credentials: Credentials,
    #[serde(default)]
    pub project: Project,
    pub remote: Option<Remote>,
    #[serde(default)]
    pub mapping: Mapping,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Credentials {
    pub host: String,
    pub email: String,
    pub auth: String,
    pub token: String,
}

impl Default for Credentials {
    fn default() -> Self {
        Self {
            host: "your_jira_host".to_string(),
            email: "your_jira_email".to_string(),
            auth: "basic".to_string(),
            token: "your_jira_token".to_string(),
        }
    }
}

impl Display for Credentials {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "\
            credentials.auth={auth}\n\
            credentials.email={email}\n\
            credentials.host={host}\n\
            credentials.token={token}\
            ",
            auth = self.auth,
            email = self.email,
            host = self.host,
            token = self.token,
        )
    }
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Project {
    pub branch: String,
    pub remote: String,
    pub exclude_phrases: Vec<String>,
}

impl Default for Project {
    fn default() -> Self {
        Self {
            branch: "development".to_string(),
            remote: "origin".to_string(),
            exclude_phrases: ["front", "mobile", "android", "ios", "be", "web", "spike", "eval"]
                .map(std::string::ToString::to_string)
                .to_vec(),
        }
    }
}

impl Display for Project {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "\
            project.branch={branch}\n\
            project.exclude_phrases=[{exclude_phrases}]\n\
            project.remote={remote}\
            ",
            branch = self.branch,
            exclude_phrases = self.exclude_phrases.join(", "),
            remote = self.remote,
        )
    }
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum RemoteProvider {
    #[serde(rename = "github")]
    GitHub,
    #[serde(rename = "gitlab")]
    Gitlab,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Remote {
    pub provider: RemoteProvider,
    pub token: String,
    host: Option<String>,
    #[serde(default)]
    pub labels: Vec<String>,
}

impl Remote {
    pub fn get_host(&self) -> &str {
        self.host.as_deref().unwrap_or({
            match self.provider {
                RemoteProvider::GitHub => "https://api.github.com",
                RemoteProvider::Gitlab => "https://gitlab.com",
            }
        })
    }
}

impl Display for Remote {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        writeln!(
            f,
            "\
            remote.provider={provider:?}\n\
            remote.token=hidden",
            provider = self.provider,
        )?;
        if let Some(url) = &self.host {
            writeln!(
                f,
                "\
                remote.url={url}"
            )?;
        }

        if !self.labels.is_empty() {
            writeln!(f)?;
            write!(f, "remote.labels=[{}]", self.labels.join(","))?;
        }
        Ok(())
    }
}

pub type MappingType = String;

#[derive(Debug, Default, PartialEq, Serialize, Deserialize)]
pub struct Mapping {
    #[serde(flatten, default, deserialize_with = "transpose_map")]
    pub entries: HashMap<String, MappingType>,
}

impl Display for Mapping {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        if self.entries.is_empty() {
            return write!(f, "");
        }

        let inverted_mapping = self.entries.clone().into_iter().fold(
            BTreeMap::<MappingType, Vec<String>>::new(),
            |mut acc, (k, v)| {
                acc.entry(v).or_default().push(k);
                acc
            },
        );

        let mut formatted_str: String =
            inverted_mapping
                .into_iter()
                .fold(String::new(), |mut acc, (k, v)| {
                    let values = v.join(", ");
                    let _ = writeln!(acc, "mapping.{k}=[{values}]");
                    acc
                });

        // removes trailing new line
        formatted_str.pop();

        write!(f, "{formatted_str}")
    }
}

fn transpose_map<'de, D>(deserializer: D) -> Result<HashMap<String, MappingType>, D::Error>
where
    D: Deserializer<'de>,
{
    let original: HashMap<MappingType, Vec<String>> = HashMap::deserialize(deserializer)?;
    let mut transposed: HashMap<String, MappingType> = HashMap::new();

    for (key, values) in original {
        for value in values {
            if value == "0" {
                continue;
            }
            transposed.entry(value).or_insert(key.clone());
        }
    }

    Ok(transposed)
}

pub fn read_config_value(key: &str) -> Result<Option<Value>, ConfigError> {
    let global_config_exists = Path::new(&get_config_global_path()).exists();
    let local_config_exists = Path::new(&get_config_local_path()).exists();

    if !global_config_exists && !local_config_exists {
        return Err(ConfigError::MissingConfig);
    }

    let mut f = Figment::new();
    if global_config_exists {
        f = f.merge(Toml::file(get_config_global_path()));
    }
    if local_config_exists {
        f = f.merge(Toml::file(get_config_local_path()));
    }

    let value: Option<Value> = f.find_value(key).ok();
    Ok(value)
}

pub fn read_config() -> Result<Config, ConfigError> {
    let global_config_exists = Path::new(&get_config_global_path()).exists();
    let local_config_exists = Path::new(&get_config_local_path()).exists();
    if !global_config_exists && !local_config_exists {
        return Ok(Config::default());
    }
    let mut f = Figment::new();
    if global_config_exists {
        f = f.merge(Toml::file(get_config_global_path()));
    }
    if local_config_exists {
        f = f.merge(Toml::file(get_config_local_path()));
    }
    let config: Config = f.extract()?;
    Ok(config)
}

pub fn create_config_if_not_exist(config_path: &str) -> Result<(), ConfigError> {
    let config_path = Path::new(&config_path);
    if config_path.exists() {
        return Ok(());
    }

    if let Some(prefix) = config_path.parent() {
        fs::create_dir_all(prefix)?;
    }

    let content = toml::to_string(&Config::default())?;
    let mut file = File::create(config_path)?;
    file.write_all(content.as_bytes())?;

    Ok(())
}

pub fn get_config_local_path() -> String {
    ".twig/config/twig.toml".to_string()
}

pub fn get_config_global_path() -> String {
    let user_config_dir = get_default_config_path();
    format!("{user_config_dir}/twig/twig.toml")
}

#[cfg(target_os = "linux")]
#[cfg(not(test))]
fn get_default_config_path() -> String {
    env::var("XDG_CONFIG_HOME").unwrap_or_else(|_| -> String { env::var("HOME").unwrap() + "/.config" })
}

#[cfg(target_os = "macos")]
#[cfg(not(test))]
fn get_default_config_path() -> String {
    env::var("XDG_DATA_HOME").unwrap_or_else(|_| -> String { env::var("HOME").unwrap() + "/Library" })
}

#[cfg(test)]
pub use tests::get_default_config_path;

#[cfg(test)]
#[allow(clippy::result_large_err)]
mod tests {
    use super::*;
    use std::io;
    use std::sync::LazyLock;
    use tempfile::TempDir;

    #[test]
    fn test_read_config_local() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;
            jail.create_file(get_config_local_path(), &build_test_file_content())?;

            let config: Config = read_config().unwrap();
            assert_eq!(config, build_test_config_model());

            Ok(())
        });
    }

    #[test]
    fn test_read_config_global() {
        figment::Jail::expect_with(|_| {
            create_global_config(&build_test_file_content()).unwrap();

            let config: Config = read_config().unwrap();
            assert_eq!(config, build_test_config_model());

            Ok(())
        });
    }

    #[test]
    fn test_read_config_local_over_global() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;
            jail.create_file(
                get_config_local_path(),
                r#"
                    [credentials]
                    host = "host_from_local"
                "#,
            )?;

            create_global_config(
                r#"
                    [credentials]
                    host = "test_host"
                    email = "test_email"
                    auth = "basic"
                    token = "super_secret"
                "#,
            )
            .unwrap();

            let config: Config = read_config().unwrap();
            assert_eq!(
                config,
                Config {
                    credentials: Credentials {
                        host: "host_from_local".to_string(),
                        email: "test_email".to_string(),
                        auth: "basic".to_string(),
                        token: "super_secret".to_string()
                    },
                    project: Project::default(),
                    remote: None,
                    mapping: Mapping {
                        entries: HashMap::default()
                    },
                }
            );

            Ok(())
        });
    }

    #[test]
    fn test_read_config_skip_zero_issue_type() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;
            jail.create_file(
                get_config_local_path(),
                &build_test_file_content_with_zero_issue_type(),
            )?;

            let config: Config = read_config().unwrap();
            assert_eq!(config, build_test_config_model_with_zero_issue_type());

            Ok(())
        });
    }

    #[test]
    fn test_create_config_if_not_exists() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;

            create_config_if_not_exist(&get_config_local_path()).unwrap();
            assert!(Path::new(&get_config_local_path()).exists());
            Ok(())
        });
    }

    fn create_global_config(content: &str) -> io::Result<()> {
        let mut global_config_file = fs::File::create(get_config_global_path())?;
        global_config_file.write_all(content.as_bytes())?;
        Ok(())
    }

    static TEMP_HOME_DIR: LazyLock<TempDir> = LazyLock::new(|| {
        let home_dir = tempfile::tempdir().unwrap();
        fs::create_dir(home_dir.path().join("twig")).unwrap();
        home_dir
    });

    pub fn get_default_config_path() -> String {
        TEMP_HOME_DIR
            .path()
            .to_path_buf()
            .into_os_string()
            .into_string()
            .unwrap()
    }

    fn build_test_file_content() -> String {
        r#"
            [credentials]
            host = "test_host"
            email = "test_email"
            auth = "basic"
            token = "super_secret"
            [project]
            branch = "main"
            remote = "myorigin"
            exclude_phrases = ["test", "super_test"]
            [remote]
            provider = "github"
            host = "remote-url"
            token = "some-token"
            labels = ["test-label1", "test-label2"]
            [mapping]
            build = ["1.1", "1.2"]
            chore = ["2"]
            ci = ["3"]
            docs = ["4"]
            feat = ["5"]
            fix = ["6"]
            perf = ["7"]
            refactor = ["8"]
            revert = ["9"]
            style = ["10"]
            temp = ["11"]
            test = ["12"]
        "#
        .to_string()
    }

    fn build_test_file_content_with_zero_issue_type() -> String {
        r#"
            [credentials]
            host = "test_host"
            email = "test_email"
            auth = "basic"
            token = "super_secret"
            [project]
            branch = "main"
            remote = "myorigin"
            exclude_phrases = ["test", "super_test"]
            [mapping]
            build = ["0"]
            chore = ["0"]
            ci = ["0"]
            docs = ["0"]
            feat = ["0"]
            fix = ["0"]
            perf = ["0"]
            refactor = ["0"]
            revert = ["0"]
            style = ["0"]
            temp = ["0"]
            test = ["0"]
        "#
        .to_string()
    }

    fn build_test_config_model() -> Config {
        Config {
            credentials: Credentials {
                host: "test_host".to_string(),
                email: "test_email".to_string(),
                auth: "basic".to_string(),
                token: "super_secret".to_string(),
            },
            project: Project {
                branch: "main".to_string(),
                remote: "myorigin".to_string(),
                exclude_phrases: vec!["test".to_string(), "super_test".to_string()],
            },
            remote: Some(Remote {
                provider: RemoteProvider::GitHub,
                host: Some("remote-url".to_string()),
                token: "some-token".to_string(),
                labels: vec![("test-label1".to_string()), ("test-label2".to_string())],
            }),
            mapping: Mapping {
                entries: HashMap::from([
                    ("1.1".to_string(), "build".to_string()),
                    ("1.2".to_string(), "build".to_string()),
                    ("2".to_string(), "chore".to_string()),
                    ("3".to_string(), "ci".to_string()),
                    ("4".to_string(), "docs".to_string()),
                    ("5".to_string(), "feat".to_string()),
                    ("6".to_string(), "fix".to_string()),
                    ("7".to_string(), "perf".to_string()),
                    ("8".to_string(), "refactor".to_string()),
                    ("9".to_string(), "revert".to_string()),
                    ("10".to_string(), "style".to_string()),
                    ("11".to_string(), "temp".to_string()),
                    ("12".to_string(), "test".to_string()),
                ]),
            },
        }
    }

    fn build_test_config_model_with_zero_issue_type() -> Config {
        Config {
            credentials: Credentials {
                host: "test_host".to_string(),
                email: "test_email".to_string(),
                auth: "basic".to_string(),
                token: "super_secret".to_string(),
            },
            project: Project {
                branch: "main".to_string(),
                remote: "myorigin".to_string(),
                exclude_phrases: vec!["test".to_string(), "super_test".to_string()],
            },
            remote: None,
            mapping: Mapping {
                entries: HashMap::new(),
            },
        }
    }
}
