# Stage 1: Build Application
FROM golang:1.23-alpine AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o /provisioner ./cmd/provisioner

# Stage 2: Deploy to a minimal scratch container
FROM scratch
WORKDIR /
COPY --from=build /provisioner /provisioner
EXPOSE 42069

CMD ["./provisioner"]