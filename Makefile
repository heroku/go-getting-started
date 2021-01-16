GOCMD=go
GOTEST=$(GOCMD) test
GORUN=${GOCMD} run

all: test

start: 
				${GORUN} server.go
test:
				$(GOTEST) -v ./handlers/*
test-1:
				$(GOTEST) -v ./handlers/handlers.go \
./handlers/handlers_test_helpers.go \
./handlers/01_get_handler_test.go

test-2:
				$(GOTEST) -v ./handlers/handlers.go \
./handlers/handlers_test_helpers.go \
./handlers/02_post_handler_test.go

test-3:
				$(GOTEST) -v ./handlers/handlers.go \
./handlers/handlers_test_helpers.go \
./handlers/03_dispatch_handler_test.go