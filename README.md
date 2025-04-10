# About The Project

This branch is dedicated for an example project. The theme for the project is a library system.

Overall System:

* Basic CRUD functions for Books, Users, Borrowing Books, and other related entities
* Integrated File Handling (Integrate with MinIO)
* Role Based Access Control
* Serve RESTful API
* Serve gRPC API
* Customizable .env file
* Provide API Documentation via Swagger Documentation
* Provide API Documentation using Exported Postman API Collection (in JSON Format)
* Dockerfile to build the project as a container
* Docker Compose for easier management of related containers
* Unit Testing
* Integration Testing
* ORM Integration (To help with query builder & Database Migration)

Future Implementation:

* OAuth2 Integration
* Prometheus and Grafana Integration for Observability and Monitoring
* Log Aggregation (probably using Loki)
* CI / CD Implementation using Github Workflows paired with Docker Hub
* Integrate Message Queueing (probably using Kafka) to simulate real-time data
* Integrate Websocket for Chatting Feature

## Tech Stack

This project will be build using :

* Main Programming Language -> Go / GoLang
* Main Database -> MariaDB
* Cache Database -> Redis
* Reverse Proxy -> NGINX

## Getting Started

There are a few software / tools that you need to install in order to run this project.

### Docker

[Docker Installation](https://www.docker.com/products/docker-desktop/)

PS : Docker Desktop is optional but very handy for beginner (including myself :D)

### Go

For straightforward way [Go](https://go.dev/doc/install)
or if you want to use a specific version [Specific Go](https://go.dev/dl/)

#### Others

If you are using window based host, WSL2 is highly recommended.

[Installing WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)

Though it is not mandatory, it is highly recommended to use WSL2 for development.

However, if you are planning to develop / use gRPC you will need a linux based host in order to compile the proto files. Hence why i added WSL2 as a recommendation.

I have not found a way to compile the proto files using Windows host. If you have any idea, please let me know.

### Installation

1. Clone the repo

   ```sh
   git clone --branch library-project https://github.com/voxtmault/mentoring.git
   ```

2. Installing Go Dependencies

   ```sh
   go get
   ```

3. Fill in the value .env file

   ```sh
   cp .env.example .env
   ```

   Then open the `.env` file and fill in the values. The default value should be enough for development purposes.

4. Change git remote url to avoid accidental pushes to base project

   ```sh
   git remote set-url origin github_username/repo_name
   git remote -v # confirm the changes
   ```

## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)
