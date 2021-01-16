FROM golang:1.13

ENV CGO_ENABLED 0

WORKDIR /src/app

RUN addgroup --system projects && adduser --system projects --ingroup projects

RUN chown -R projects:projects /src/app

USER projects

COPY . .

RUN go install -v ./...
