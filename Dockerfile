# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23 AS build-stage

RUN useradd -u 10001 nonroot

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Run the tests
RUN CGO_ENABLED=0 go test -v ./... --cover

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /reaphur

# Deploy the application binary into a lean image
FROM scratch AS final

WORKDIR /

COPY --from=build-stage /reaphur /reaphur

EXPOSE 8080

COPY --from=0 /etc/passwd /etc/passwd
USER nonroot

ENTRYPOINT ["/reaphur"]
