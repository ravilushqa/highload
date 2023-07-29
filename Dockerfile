FROM golang:1.19 as build
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
RUN apk add --no-cache tzdata
CMD ./app
