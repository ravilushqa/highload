FROM --platform=linux/arm64 golang:1.24.2-bullseye as build
ARG SERVICE_PATH
WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o bin/app ./${SERVICE_PATH}

FROM alpine
ARG SERVICE_PATH
WORKDIR /opt/app
COPY --from=build /opt/app/bin/app ./app
COPY --from=build /opt/app/${SERVICE_PATH}/resources ./resources
COPY --from=build /opt/app/${SERVICE_PATH}/public ./public
RUN apk add --no-cache tzdata
CMD ./app
