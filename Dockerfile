FROM golang:1.24 as build
RUN apt-get update && apt-get install -y \
    libssl-dev \
    pkg-config \
    build-essential \
    gcc \
    libc6-dev
RUN go install github.com/go-delve/delve/cmd/dlv@latest
ARG SERVICE_PATH

WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o bin/app ./${SERVICE_PATH}

FROM alpine
WORKDIR /opt/app
COPY --from=build /opt/app/bin/app ./app
COPY --from=build /go/bin/dlv ./dlv
#RUN apk add --no-cache tzdata
RUN apk add --no-cache ca-certificates libc6-compat tzdata
CMD ["./app"]
# CMD ["./dlv", "--listen=:40000", "--headless", "--continue", "--api-version=2", "--accept-multiclient", "exec", "./app"]
