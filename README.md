# Discog API

This is a web application that fetches and serves data from Discogs.

## Prerequisites

Make sure you have the following installed on your system:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## How to Run

1. Clone this repository:
```bash
git clone https://github.com/prince950408/Tempest
cd Tempest
```
2. Build and run the application using Docker Compose:
```
docker-compose up --build
```
3. Wait a few minutes for the application to fetch data from Discogs.
4. If it says `Server is running on port 8080 `, you can access it at:
[http://localhost:3000](http://localhost:3000)