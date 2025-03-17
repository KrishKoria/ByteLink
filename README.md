# ByteLink

A Compact URL Shortener Using Go and Redis

## Overview

ByteLink is a URL shortening service built with Go and Redis. It allows users to create short URLs that redirect to long URLs, making it easier to share and manage links.

## Features

- Generate short URLs for long links
- Store URL mappings in Redis
- Redirect short URLs to their original long URLs
- Simple and efficient URL shortening algorithm using SHA-256 and Base58 encoding

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/KrishKoria/ByteLink.git
    cd ByteLink
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up Redis:
    - Ensure Redis is installed and running on `127.0.0.1:6379`.

## Usage

1. Start the server:
    ```sh
    go run main.go
    ```

2. Create a short URL:
    - Send a POST request to `http://localhost:8080/create` with the following JSON body:
        ```json
        {
            "long_url": "https://www.example.com",
            "user_id": "your-unique-user-id"
        }
        ```

3. Redirect to the long URL:
    - Access the short URL in your browser, e.g., `http://localhost:8080/{shortURL}`.

## API Endpoints

- **Create Short URL**
    - **URL:** `/create`
    - **Method:** `POST`
    - **Request Body:**
        ```json
        {
            "long_url": "https://www.example.com",
            "user_id": "your-unique-user-id"
        }
        ```
    - **Response:**
        ```json
        {
            "message": "Short URL created successfully",
            "short_url": "http://localhost:8080/{shortURL}"
        }
        ```

- **Redirect to Long URL**
    - **URL:** `/{shortURL}`
    - **Method:** `GET`
    - **Response:** Redirects to the original long URL.

## Testing

Run the tests using the following command:
```sh
go test ./...