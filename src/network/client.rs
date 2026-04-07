use crate::network::api::HttpClient;
use crate::network::model::{AUTH_BASIC, AUTH_BEARER, Credentials};
use anyhow::{Context, Ok, Result};
use base64::Engine;
use base64::engine::general_purpose;
use reqwest::blocking::{Client, ClientBuilder, RequestBuilder};
use reqwest::{Method, blocking::Response, header};
use serde::Serialize;
use serde::de::DeserializeOwned;

pub trait ApiError {
    fn errors(&self) -> Vec<String>;
}

pub struct TwigClient {
    host: String,
    client: Client,
}

impl TwigClient {
    pub fn new(credentials: Credentials) -> Result<Self> {
        let mut builder = Client::builder();
        builder = configure_headers(builder, &credentials);

        let host = format!("https://{host}/rest/api/2", host = credentials.host);
        let client = builder.build().context("failed to build client")?;
        Ok(Self { host, client })
    }

    fn request(&self, method: Method, path: &str) -> RequestBuilder {
        let url = format!("{host}/{path}", host = self.host, path = path);
        self.client.request(method, url)
    }
}

impl HttpClient for TwigClient {
    fn get<T: DeserializeOwned, E: ApiError + DeserializeOwned>(
        &self,
        path: &str,
        params: Vec<(&str, &str)>,
    ) -> Result<T> {
        let response = self.request(Method::GET, path).query(&params).send()?;
        raise_for_status::<E>(response)
            .context("GET request failed")?
            .json::<T>()
            .context("failed to parse GET response")
    }

    fn post<T: DeserializeOwned, S: Serialize, E: ApiError + DeserializeOwned>(
        &self,
        path: &str,
        body: S,
    ) -> Result<T> {
        let response = self.request(Method::POST, path).json(&body).send()?;
        raise_for_status::<E>(response)
            .context("GET request failed")?
            .json::<T>()
            .context("failed to parse GET response")
    }
}

fn configure_headers(builder: ClientBuilder, credentials: &Credentials) -> ClientBuilder {
    let mut headers = header::HeaderMap::new();

    let auth_value = create_auth_header(credentials);
    let mut header = header::HeaderValue::from_str(auth_value.as_str()).unwrap();
    header.set_sensitive(true);

    headers.insert("Authorization", header);
    builder.default_headers(headers)
}

fn create_auth_header(credentials: &Credentials) -> String {
    let auth = match credentials.auth.as_str() {
        AUTH_BASIC => (
            String::from("Basic"),
            basic_auth(&credentials.email, &credentials.token),
        ),
        AUTH_BEARER => (String::from("Bearer"), String::from(&credentials.token)),
        _ => panic!("Unsupported auth type"),
    };

    format!("{auth_type} {value}", auth_type = auth.0, value = auth.1)
}

fn basic_auth(username: &str, password: &str) -> String {
    let auth = format!("{}:{}", username, password);
    general_purpose::STANDARD.encode(auth.as_bytes())
}

fn raise_for_status<E: ApiError + DeserializeOwned>(response: Response) -> Result<Response> {
    if let Err(err) = response.error_for_status_ref() {
        let mut all_errors = anyhow::Error::new(err);
        let error_response = response.json::<E>().unwrap();
        for err in error_response.errors() {
            all_errors = all_errors.context(err.clone());
        }
        return Err(all_errors);
    }
    Ok(response)
}
