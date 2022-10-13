###############################################################################
#############################     BUILD SETUP     #############################
###############################################################################

FROM golang:latest as builder
LABEL owner="Ritalin Company"
LABEL maintainer="Arie Brainware"

# Golang build environment
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  CC=musl-gcc

COPY . .

# Install dependencies
RUN go mod download


# Build image
RUN go build -a -installsuffix cgo -o main.go

# production stage
FROM alpine:latest

# FIX: fix invalid timezone issue
RUN apk add --no-cache tzdata

# Expose port, follow this guideline to see existed microser
EXPOSE 3110

WORKDIR /root/

COPY --from=builder /root main
CMD ["/bin/sh","-c","source /etc/profile && ./main"]
