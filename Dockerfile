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
RUN go build main.go

FROM scratch

COPY --from=builder build/main .
COPY --from=builder build/module/mail /module/mail/

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["/main", "serve"]