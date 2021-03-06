FROM golang:1.15 as build
RUN go get github.com/go-delve/delve/cmd/dlv
ARG SERVICE_PATH

WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o bin/app ./${SERVICE_PATH}

FROM alpine
WORKDIR /opt/app
COPY --from=build /opt/app/bin/app ./app
COPY --from=build /go/bin/dlv ./dlv
#RUN apk add --no-cache tzdata
RUN apk add --no-cache ca-certificates libc6-compat tzdata
#CMD ./app
CMD ["./dlv", "--listen=:40000", "--headless", "--continue", "--api-version=2", "--accept-multiclient", "exec", "./app"]
