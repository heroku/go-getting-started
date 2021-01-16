# FIFA World Cup Winners

This project exposes a Web API for accessing historic data from
the FIFA World Cup championship.

## Running tests

A proper Go environment is required in order to run this project.
Once setup, tests can be run with the following command:

`go test -v ./handlers/`

## Running the server

Once all tests are passing, the server can be started with
the `go run server.go` command.

## Testing the API manually

Start the server with `go run server.go` and then
use the example commands printed to the console to
test the program.

### Running with Docker

To build the image from the Dockerfile, run:

`docker build -t project-fifa-world-cup .`

To start an interactive shell, run:

`docker run -it --rm --name run-fifa project-fifa-world-cup`

From inside the shell, run the tests with:

`go test handlers/*`
