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
 - Automatic cleanup of orphaned URLs via background job
 - Status monitoring for background maintenance tasks

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

 2. Access the service at `http:localhost:8080/`

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

 - **Delete URL Mapping**
   - **URL:** `/api/url`
   - **Method:** `DELETE`
   - **Request Body:**
     ```json
     {
         "short_url": "{shortURL}",
         "user_id": "your-unique-user-id"
     }
     ```
   - **Response:** Confirms deletion of the mapping

 - **Get Cleanup Job Status**
   - **URL:** `/api/admin/cleanup-status`
   - **Method:** `GET`
   - **Response:**
     ```json
     {
         "last_run_time": "2023-04-30T12:34:56Z",
         "total_urls_removed": 42,
         "is_running": true,
         "run_interval_minutes": 1440
     }
     ```

 ## Database Schema

 The service uses SQLite with the following schema:
 - URLs table: Stores long URLs with unique IDs
 - Mappings table: Associates short URLs with URL IDs and user IDs

 ## Database Maintenance

 ByteLink includes an automatic maintenance system:

 - Background job runs periodically to clean up orphaned URLs
 - Removes URLs that no longer have any mappings pointing to them
 - Improves database efficiency and reduces storage requirements
 - Runs on a configurable schedule (default: every 24 hours)
 - Non-blocking implementation that doesn't affect user operations
 - Provides status monitoring through admin API endpoint
 - Logs job execution details for troubleshooting

 ## Monitoring

 The cleanup job can be monitored through:

 - Logs: The system prints cleanup job status to standard output
 - API endpoint: Status data available at `/api/admin/cleanup-status`
 - Status metrics include last run time and total URLs removed
 - Job failures are logged with detailed error messages

 ## Testing

 Run the tests using the following command:
 ```sh
 go test ./...
 ```

 ## Error Handling

 - Duplicate URL creation attempts for the same user are handled gracefully
 - Proper error responses for invalid requests
 - Cleanup job failures are logged but don't disrupt service operation