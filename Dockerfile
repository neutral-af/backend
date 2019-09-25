FROM golang:1.12 as builder

WORKDIR /go/src/github.com/neutral-af/backend
COPY go.mod go.sum ./
RUN GO111MODULE=on go get -v

COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o main -v cmd/carbon-offsets/main.go

FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/neutral-af/backend/main /main

# Run the web service on container startup.
CMD ["/main"]