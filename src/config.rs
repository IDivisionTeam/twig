use std::{
    env,
    fs::{self, File},
    io::{self, Write},
    path::Path,
};

use figment::{
    providers::{Format, Toml},
    Figment,
};
use serde::{Deserialize, Serialize};

use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigError {
    #[error("failed to extract config")]
    Extract(#[from] Box<figment::Error>),
    #[error("failed to parse config")]
    Parse(#[from] toml::ser::Error),
    #[error("failed to create file")]
    File(#[from] io::Error),
}

impl From<figment::Error> for ConfigError {
    fn from(err: figment::Error) -> Self {
        ConfigError::Extract(Box::new(err))
    }
}

#[derive(Debug, Default, PartialEq, Serialize, Deserialize)]
pub struct Config {
    pub credentials: Credentials,
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
        fs::create_dir_all(prefix)?
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
fn get_default_config_path() -> String {
    env::var("XDG_CONFIG_HOME")
        .unwrap_or_else(|_| -> String { env::var("HOME").unwrap() + "/.config" })
}

#[cfg(target_os = "macos")]
fn get_default_config_path() -> String {
    env::var("XDG_DATA_HOME").unwrap_or_else(|_| -> String { "~/Library/".to_string() })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_read_config_local() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;
            jail.create_file(
                get_config_local_path(),
                r#"
                    [credentials]
                    host = "test_host"
                    email = "test_email"
                    auth = "basic"
                    token = "super_secret"
                "#,
            )?;

            let config: Config = read_config().unwrap();
            assert_eq!(
                config,
                Config {
                    credentials: Credentials {
                        host: "test_host".to_string(),
                        email: "test_email".to_string(),
                        auth: "basic".to_string(),
                        token: "super_secret".to_string()
                    }
                }
            );

            Ok(())
        });
    }

    #[test]
    fn test_read_config_global() {
        figment::Jail::expect_with(|jail| {
            let current_dir = jail.directory().display().to_string();
            jail.set_env("XDG_CONFIG_HOME", &current_dir);

            jail.create_dir(current_dir + "/twig")?;
            jail.create_file(
                get_config_global_path(),
                r#"
                    [credentials]
                    host = "test_host"
                    email = "test_email"
                    auth = "basic"
                    token = "super_secret"
                "#,
            )?;

            let config: Config = read_config().unwrap();
            assert_eq!(
                config,
                Config {
                    credentials: Credentials {
                        host: "test_host".to_string(),
                        email: "test_email".to_string(),
                        auth: "basic".to_string(),
                        token: "super_secret".to_string()
                    }
                }
            );

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

            let current_dir = jail.directory().display().to_string();
            jail.set_env("XDG_CONFIG_HOME", &current_dir);

            jail.create_dir(current_dir + "/twig")?;
            jail.create_file(
                get_config_global_path(),
                r#"
                    [credentials]
                    host = "test_host"
                    email = "test_email"
                    auth = "basic"
                    token = "super_secret"
                "#,
            )?;

            let config: Config = read_config().unwrap();
            assert_eq!(
                config,
                Config {
                    credentials: Credentials {
                        host: "host_from_local".to_string(),
                        email: "test_email".to_string(),
                        auth: "basic".to_string(),
                        token: "super_secret".to_string()
                    }
                }
            );

            Ok(())
        });
    }

    #[test]
    fn test_create_config_if_not_exists() {
        figment::Jail::expect_with(|jail| {
            jail.create_dir(".twig/config")?;

            create_config_if_not_exist(&get_config_local_path()).unwrap();
            assert_eq!(Path::new(&get_config_local_path()).exists(), true);
            Ok(())
        });
    }
}
