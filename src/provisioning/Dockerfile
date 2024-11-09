# Stage 1: Build Application
FROM golang:1.23-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY *.go ./
RUN go build -o /provisioner

# Stage 2: Deploy to a minimal scratch container
FROM scratch
WORKDIR /
COPY --from=build /provisioner /provisioner
EXPOSE 42069

CMD ["./provisioner"]