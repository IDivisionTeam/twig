use crate::network::model::{Credentials, AUTH_BASIC, AUTH_BEARER};
use base64::engine::general_purpose;
use base64::Engine;
use reqwest::{header, Client, ClientBuilder, Method, RequestBuilder};

pub struct TwigClient {
    creds: Credentials,
    client: Client,
}

impl TwigClient {
    pub fn new(credentials: Credentials) -> Self {
        let mut builder = Client::builder();

        builder = configure_headers(builder, &credentials);

        Self {
            creds: credentials,
            client: builder.build().unwrap_or(Default::default()),
        }
    }

    pub fn send() {
        // TODO: impl
    }

    fn prepare(self: &Self, method: Method, path: String) -> RequestBuilder {
        let url = format!(
            "https://{host}/rest/api/2/{path}",
            host = self.creds.host,
            path = path
        );

        self.client.request(method, url)
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
