FROM golang:1.17 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o ontest .

FROM alpine:3.13 AS certer

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=certer /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder build/ontest .
COPY --from=builder build/module/mail /module/mail/
COPY ontest /usr/bin/ontest

USER 1000

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
ENTRYPOINT ["ontest"]
CMD [ "serve" ]