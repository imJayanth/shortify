Shortify - URL Shortener Service
================================

Description
-----------

Shortify is a URL shortening service written in Golang. It provides two endpoints, POST and GET, for creating and retrieving shortened URLs, respectively. The original URL is base62 decoded and shortened with the application's domain, which is valid for 24 hours. The original URL, shortened URL, creation time, and expiry time are stored in both MySQL database and Redis. The POST endpoint checks if the URL already exists in Redis, and if not, it generates and saves it. The GET endpoint checks if the URL exists in Redis, and if not, it retrieves it from the database and does a 301 redirect to the original URL.

Installation
------------

To run Shortify, you need to have Docker and Docker Compose installed on your system.

1.  Clone the repository to your local machine.
    `git clone https://github.com/<your-username>/shortify.git`

2.  Navigate to the project directory.
    `cd shortify`

3.  Build the Docker image.
    `docker build -t shortify .`

4.  Run the Docker containers using Docker Compose.
    `docker-compose up -d`

Shortify should now be up and running at `http://localhost:8081`.

Usage
-----

### POST Endpoint

The POST endpoint is used to create a shortened URL.

URL: `/`
Method: `POST`
Headers:
`Content-Type: application/json`
Body:
`{
    "original_url": "https://example.com/very-long-url-that-needs-to-be-shortened"
}`

Response:
`{
    "data": {
        "url": "localhost:8081/013"
    },
    "error": null
}`

### GET Endpoint

The GET endpoint is used to retrieve the original URL from a shortened URL.

URL: `data.url` from previous `POST` request

Method: `GET`

Response: A 301 redirect to the original URL.

### Postman Collection

A Postman collection is included in the repository (`shortify.postman_collection.json`) that contains sample requests for the POST and GET endpoints.
