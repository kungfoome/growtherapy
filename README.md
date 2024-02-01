# Wikipedia View Count API Wrapper

## Overview

This application is a Go web server that serves as a wrapper around the Wikipedia API. It retrieves view counts for a specified article within a given month and year, considering leap years for accurate date calculation. The view count is returned in JSON format.

## Requirements

- Go 1.21 or later
- Docker (optional for containerization)

## Building and Running

### Locally

To run this application locally, follow these steps:

1. Clone the repository and navigate to its directory.
2. Build the application:

    ```shell
    go build -o wikiViewCounter ./cmd/server
    ```

    or 

    ```shell
    go run ./cmd/server
    ```

3. Run the application:

    ```shell
    ./wikiViewCounter
    ```

The server will start on port 8080. Adjust the port in the main.go file if needed.

### Using Docker

To containerize and run the application using Docker, execute the following commands:

1. Build the Docker image:

    ```shell
    docker build -t wiki-view-counter .
    ```

2. Run the container:

    ```shell
    docker run -p 8080:8080 wiki-view-counter
    ```

## API Usage

To fetch the view count for a specific article, make a GET request to the `/getViews` endpoint with the `article`, `year`, and `month` query parameters:

```http
GET /getViews?article=<ARTICLE_TITLE>&year=<YYYY>&month=<MM>
```

### Example

```shell
curl "http://localhost:8080/getViews?article=Go_(programming_language)&year=2023&month=04"
```

Response:

```shell
{
    "views": 123456
}
```

# Testing

```shell
go test ./...
```