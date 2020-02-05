ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine AS dev

# Setting build variables
ENV GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux

ENV APP_PATH="/gokv" \
    APP_NAME="main"

# Install git
RUN apk add --update git

# Change to app directory
WORKDIR ${APP_PATH}

# Cache Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy rest of the code
COPY . .

# Build binary files
RUN go build -ldflags="-s -w" -o  ${APP_NAME} ${APP_PATH}/examples/${APP_NAME}
RUN chmod +x ${APP_NAME}

# Final build image
FROM alpine AS prod

ENV APP_PATH="/gokv" \
    APP_NAME="main" \
    REST_PORT=8888 \
    GRPC_PORT=9999

# Change to app directory
WORKDIR ${APP_PATH}

# Copy binary from previous stage
COPY --from=dev ${APP_PATH}/${APP_NAME} ${APP_NAME}

# Expose Rest and gRPC port
EXPOSE ${REST_PORT}
EXPOSE ${GRPC_PORT}

ENTRYPOINT ["/gokv/main"]
