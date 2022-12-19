# image-resizing-server-go
HTTP API for uploading, optimizing, and serving images.
Written in Go. Application consists of 2 parts: HTTP API server for uploading and downloading images and server for optimizing and resizing images. Images are sent from HTTP server for optimizing via RabbitMq queue.
Image optimization in this realization is simple resizing into smaller size images (75%, 50% and 25% from original size).

## usage
To run locally, you need RabbitMq server running. You can execute `make run-rabbit` – that will start RabbitMq docker container on ports 5672 and 15672 on your local machine. If you'd like to start your own RabbitMq server, please make sure to change configuration files in `configs/` directory.
Then run
```
make run-api
make run-opt
```
It will start HTTP API and optimization servers respectively.

## clean code
Application was built with clean code principles in mind. Everything inside both API and optimization servers relies upon interfaces (publisher/consumer, services, repository). It is no trouble to replace some components with other using different libraries or technologies.
For example, currently application uses file storage implementation of repository, but as Repository interface is used, it can be replaced with real database easily.

## graceful shutdown
Both servers may perform difficult and time consuming operations (sending, receiving and optimizing images), so in order to prevent data loss when server is being stopped, graceful shutdown was implemented. It waits for all started proccesses to finish and only then stops application (or waits for timeout if running proccesses can't be finished fast enough)

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

