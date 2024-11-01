use std::io::{Read, Write};
use std::net::TcpStream;
use dotenv::dotenv;
use std::env;
use std::path::Path;

/// Load .env file from project root
/// 
/// Returns WOLFRAM_SHORT_RESPONSE
fn load_dotenv() -> String {
    // Establish project directory, env is in root
    let project_dir = Path::new(env!("CARGO_MANIFEST_DIR"));
    let env_path = project_dir.join("../.env");
    dotenv().ok();
    dotenv::from_path(env_path).ok();
    let app_id = env::var("WOLFRAM_SHORT_RESPONSE").expect("WOLFRAM_SHORT_RESPONSE is not set");
    return app_id;
}

/// Get query from stdin and return it
/// 
/// Returns query
fn get_query() -> String {
    print!("Enter query: ");
    std::io::stdout().flush().expect("Failed to flush stdout");
    let mut input = String::new();
    std::io::stdin().read_line(&mut input).expect("Failed to read line");
    return input;
}

/// URL encode input
/// 
/// Returns encoded input for proper http requests
fn url_encode(input: &str) -> String {
    input.chars().map(|c| {
        match c {
            // URL-safe characters remain unchanged
            'A'..='Z' | 'a'..='z' | '0'..='9' => c.to_string(),
            ' ' => "%20".to_string(),

            // Other characters are percent-encoded
            _ => format!("%{:02X}", c as u8),
        }
    }).collect()
}

/// Query WolframAlpha API
/// 
/// Returns result using api key loadeed in load_dotenv()
fn query(query: &str) -> String {
    // Create HTTP request information with API key and encoded query
    let input = url_encode(query);
    let host = "api.wolframalpha.com";
    let port = 80;
    let path = format!("/v2/query?input={}&appid={}",
        input, load_dotenv()
    );

    // Connect to server and send request
    let address = format!("{}:{}", host, port);
    let mut stream = TcpStream::connect(address).expect("Could not connect to server");
    let request = format!(
        "GET {} HTTP/1.1\r\nHost: {}\r\nConnection: close\r\n\r\n",
        path, host
    );
    stream.write_all(request.as_bytes()).expect("Failed to send request");

    // Read and parse response
    let mut response = String::new();
    stream
        .read_to_string(&mut response)
        .expect("Failed to read response");
    return response;
}

/// Parse response from WolframAlpha API
/// 
/// Returns the plaintext response from the API output
fn parse_response(response: &str) -> String {
    // Find the "Result" pod
    if let Some(result_start) = response.find("<pod title='Result'") {
        // Find the start and end of the plaintext
        if let Some(plaintext_start) = response[result_start..].find("<plaintext>") {
            if let Some(plaintext_end) = response[result_start..].find("</plaintext>") {
                // Extract and return the result from the plaintext
                let result = &response[result_start + plaintext_start + "<plaintext>".len()..
                                        result_start + plaintext_end];
                return result.trim().to_string();
            }
        }
    }
    return "No result found.".to_string();
}

fn main() {
    println!("Result: {}", parse_response(&query(&get_query())));
}