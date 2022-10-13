# Build Stage
# First pull Golang image
FROM golang:1.19.2-alpine3.16 as build-env
 
# Set environment variable
ENV APP_NAME evo-shortner
ENV CMD_PATH main.go
 
# Copy application data into image
COPY . $GOPATH/src/github.com/ariebrainware/$APP_NAME
WORKDIR $GOPATH/src/github.com/ariebrainware/$APP_NAME
 
# Budild application
RUN go mod download
RUN go mod vendor

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/github.com/ariebrainware/$APP_NAME/$CMD_PATH
 
# Run Stage
FROM alpine:3.14
 
# Set environment variable
ARG APP_NAME \
  VERSION \
  PORT \
  KEY_LENGTH \
  ROOT_HOST \
  ENVIRONMENT \
  MONGO_HOST \
  MONGO_DATABASE \
  MONGO_COLLECTION \
  MONGO_PASSWORD \
  ALLOWED_ORIGIN
 
# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
 
# Expose application port
EXPOSE 8081
 
# Start app
CMD ./$APP_NAME