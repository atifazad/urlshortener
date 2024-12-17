# URL Shortener

A simple URL shortener application written in Go. This application provides endpoints to shorten URLs and redirect to the original URLs.

## Features

- Shorten URLs
- Redirect to original URLs
- Persistent storage using JSON file

## Requirements

- Go 1.20 or later

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/urlshortener.git
   cd urlshortener
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

1. Run the application:
   ```sh
   go run *.go
   ```

The server will start on http://localhost:8080.

2.  Shorten a URL:

        ```sh
        curl -X POST -H "Content-Type: application/json" -d '{"url": "https://www.example.com"}' http://localhost:8080/shorten
        ```

    The response will contain the shortened URL.

3.  Redirect to the original URL:
    `sh
curl -L http://localhost:8080/<short_url>
`
    Replace <short_url> with the shortened URL path.

## Testing

1. Run the tests:
   `sh
 go test -v
 `
   This will execute the unit tests and display the results.
