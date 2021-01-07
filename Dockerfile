FROM golang:1.15 as build-env
WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env AS build
WORKDIR /opt/app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/app ./${SERVICE_PATH}

FROM alpine
WORKDIR /opt/app
COPY --from=build /opt/app/bin/app ./app
COPY --from=build /opt/app/resources ./resources
COPY --from=build /opt/app/public ./public
RUN apk add --no-cache tzdata
CMD ./app
