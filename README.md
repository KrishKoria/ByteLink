 # ByteLink

 A Compact URL Shortener Using Go, SQLite, and Redis

 ## Overview

 ByteLink is a URL shortening service built with Go that allows users to create short URLs that redirect to long URLs, making it easier to share and manage links. The service uses SQLite for persistent storage and Redis for caching.

 ## Features

 - Generate short URLs for long links
 - SQLite database for persistent storage of URL mappings
 - Redis caching for improved performance
 - User-specific URL management
 - Redirect short URLs to their original destinations
 - URL shortening algorithm using SHA-256 and Base58 encoding
 - Unique constraint handling to prevent duplicate mappings

 ## Installation

 1. Clone the repository:
    ```sh
    git clone https:github.com/KrishKoria/ByteLink.git
    cd ByteLink
    ```

 2. Install dependencies:
    ```sh
    go mod tidy
    ```

 3. Set up SQLite and Redis:
    - Ensure Redis is installed and running on `127.0.0.1:6379`
    - SQLite database will be created automatically

 ## Usage

 1. Start the server:
    ```sh
    go run main.go
    ```

 2. Access the service at `http://localhost:8080/`

 ## API Endpoints

 - **Create Short URL**
   - **URL:** `/create`
   - **Method:** `POST`
   - **Request Body:**
     ```json
     {
         "long_url": "https:www.example.com",
         "user_id": "your-unique-user-id"
     }
     ```
   - **Response:**
     ```json
     {
         "message": "Short URL created successfully",
         "short_url": "http:localhost:8080/{shortURL}"
     }
     ```

 - **Redirect to Long URL**
   - **URL:** `/{shortURL}`
   - **Method:** `GET`
   - **Response:** Redirects to the original long URL

 - **Get Specific URL for User**
   - **URL:** `/api/url?short_url={shortURL}&user_id={userID}`
   - **Method:** `GET`
   - **Response:** Returns the long URL associated with the short URL for the specific user

 - **Get All User URLs**
   - **URL:** `/api/urls?user_id={userID}`
   - **Method:** `GET`
   - **Response:** Returns all URL mappings for the specified user

 ## Database Schema

 The service uses SQLite with the following schema:
 - URLs table: Stores long URLs with unique IDs
 - Mappings table: Associates short URLs with URL IDs and user IDs

 ## Testing

 Run the tests using the following command:
 ```sh
 go test ./...
 ```

 ## Error Handling

 - Duplicate URL creation attempts for the same user are handled gracefully
 - Proper error responses for invalid requests