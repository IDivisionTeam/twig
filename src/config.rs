use std::{
    env,
    fs::{self, File},
    io::{self, Write},
    path::Path,
};

use figment::{
    Figment,
    providers::{Format, Toml},
};
use serde::{Deserialize, Serialize};

use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigError {
    #[error("failed to extract config")]
    ExtractError(#[from] figment::Error),
    #[error("failed to parse config")]
    ParseError(#[from] toml::ser::Error),
    #[error("failed to create file")]
    FileError(#[from] io::Error),
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

    match config_path.parent() {
        Some(prefix) => fs::create_dir_all(prefix).unwrap(),
        None => (),
    };
    let content = toml::to_string(&Config::default())?;
    let mut file = File::create(config_path)?;
    file.write_all(content.as_bytes())?;

    Ok(())
}

pub fn get_config_local_path() -> String {
    return ".twig/config/twig.toml".to_string();
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
    env::var("XDG_CONFIG_HOME").unwrap_or_else(|_| -> String { "~/Library/Preferences/" });
}

