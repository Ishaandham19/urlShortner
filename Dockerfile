FROM golang as builder
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build -o app

FROM alpine as production
# Add certificates
RUN apk add --no-cache ca-certificates
# Copy built binary from builder
COPY --from=builder app .
# Expose port
EXPOSE 8080
# Exec built binary
CMD ./app
