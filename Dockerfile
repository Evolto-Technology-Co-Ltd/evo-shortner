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
FROM alpine:3.16
 
WORKDIR /root/

# Copy only required data into this image
COPY --from=build-env /evo-shortner /root/.
COPY --from=build-env go/src/github.com/ariebrainware/evo-shortner/app.env /root/.

# Expose application port
EXPOSE 3110
 
# Start app
CMD ["/bin/sh","-c","/root/evo-shortner"]