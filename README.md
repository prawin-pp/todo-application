# Todo Application

This project contains two workspaces, one for a web application built using the Svelte framework and one for an API built using Golang.

## Installation

To get started with this project, simply navigate to the main workspaces directory and run the Makefile using the following command:

```
make
```

This will run the Docker Compose up and begin running web application, api server and PostgreSQL.


After that you need to run migration script by run the following commands:

```
make migrate
```

## Web

The web application is built using the Svelte framework, a modern and efficient frontend framework. To work on the web application, navigate to the "web" directory and run the following commands:

```
npm install
npm run dev
```

This will install all necessary dependencies and start the development server for the web application.


## API

The API is built using Golang, a powerful and efficient backend programming language. To work on the API, navigate to the "api" directory and run the following commands:

```
go get ./...
go run ./cmd/server/main.go
```

This will start the API server.

## Usage

Once the web application and API are running, you can access the web application by navigating to http://localhost:5173/ in your web browser. From there, you can interact with the application and the API in a variety of ways.

## Default Test User
```
username: tester01 , password: 1111
username: tester02 , password: 2222
username: tester03 , password: 3333
```