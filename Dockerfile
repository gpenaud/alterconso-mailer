FROM golang:alpine AS build

# Add Maintainer Info
LABEL maintainer="guillaume.penaud@gmail.com"

RUN \
  apk add --no-cache git &&\
  mkdir /application

WORKDIR /application

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Download all the dependencies
RUN go mod download

# Set build argument variables for build
ARG PROJECT
ARG RELEASE
ARG COMMIT
ARG BUILD_TIME

# Set environment variable
ENV CGO_ENABLED=0
ENV GOOS=linux

# Build the binary
RUN go build \
  -a \
  -installsuffix cgo \
  -ldflags "-s -w \
    -X main.Release=${RELEASE} \
    -X main.Commit=${COMMIT} \
    -X main.BuildTime=${BUILD_TIME}" \
    -o /alterconso-mailer \
  /application/main.go

# ---------------------------------------------------------------------------- #

FROM alpine:latest

RUN addgroup alterconso-mailer
RUN adduser --system --disabled-password --home /alterconso-mailer alterconso-mailer

WORKDIR /alterconso-mailer
USER alterconso-mailer

COPY --chown=alterconso-mailer:alterconso-mailer --from=build /alterconso-mailer .
COPY --chown=alterconso-mailer:alterconso-mailer --from=build /application/config.yaml .
COPY --chown=alterconso-mailer:alterconso-mailer --from=build /application/secrets.yaml .

EXPOSE 5000

ENTRYPOINT ["./alterconso-mailer"]
CMD ["serve"]
