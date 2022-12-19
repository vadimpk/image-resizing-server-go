# image-resizing-server-go
HTTP API for uploading, optimizing, and serving images. Written in Go. 

Application consists of 2 services: an HTTP server for uploading and downloading images and an optimizing service for resizing images.
Images are sent from the HTTP server to the optimizing service via the RabbitMq queue after uploading.
\
Local file storage is used as a repository for storing images. 
Image optimization is simply resizing into smaller-sized images (75%, 50%, and 25% from original size).

## usage
To run locally, you need RabbitMq server running. You can execute `make run-rabbit` which will start RabbitMq docker container on ports 5672 and 15672 on your local machine.
If you'd like to start your own RabbitMq server, please make sure to change configuration files in `configs/` directory.
\
Then run
```
make run-api
make run-opt
```
It will start HTTP API and optimization servers respectively.

In `configs/` directory you can set the directory for optimized images to be saved into (default `images/`). 

## uploading / downloading
Open `index.html` to manually upload images to server or run:
```
go run client.go --filename <your-img-filename> --times <number-of-requests-to-send>
```
to run any amount of requests to the server with the same specified image (default `test-img.jpg`).

To download visit `/download/{img-id}?quality=100/75/50/25` endpoint of HTTP server in your browser.

## clean code
Application was built with clean code principles in mind.
Everything inside both API and optimization servers relies upon interfaces (publisher/consumer, services, repository).
It is no trouble to replace some components with others using different libraries or technologies.
\
For example, currently, application uses file storage implementation of a repository, but as `Repository` interface is used, it can be replaced with a real database easily.

## graceful shutdown
Both servers may perform difficult and time-consuming operations (sending, receiving, and optimizing images), so in order to prevent data loss when server is being stopped, graceful shutdown was implemented.
It waits for all started processes to finish and only then stops the application (or waits for timeout if running processes can't be finished fast enough).

## stack
  - [github.com/gorilla/mux](github.com/gorilla/mux)
  - [github.com/rabbitmq/amqp091-go](github.com/rabbitmq/amqp091-go)
  - [github.com/spf13/viper](github.com/spf13/viper)
  - [github.com/teris-io/shortid](github.com/teris-io/shortid)
  - [github.com/nfnt/resize](github.com/nfnt/resize)
  
## architecture
Visualization of how dependency injection is used inside both services (note: not full application architecture)

![img-resizing-server-go-5](https://user-images.githubusercontent.com/65962115/208397435-df23c8e3-ca6c-4327-87fe-e819f5b98918.jpg)

## ways to improve
  - implement Postgres instead of file storage
  - create docker images for both services (previous step is required, or at least preferable, to implement this one)
  - use multiple goroutines to consume messages asynchronously on rabbitmq side (problem: how to stop gracefully?)
  - send big images in parts to prevent data loss and optimize queue

