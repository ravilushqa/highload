FROM golang:1.15 as build-env
ENV NAME "app"
WORKDIR /opt/${NAME}
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env AS build
ENV NAME "app"
WORKDIR /opt/${NAME}
COPY . .
RUN CGO_ENABLED=0 go build -o bin/${NAME}

FROM alpine
ENV NAME "app"
WORKDIR /opt/${NAME}
COPY --from=build /opt/${NAME}/bin/${NAME} ./${NAME}
COPY --from=build /opt/${NAME}/resources ./resources
COPY --from=build /opt/${NAME}/public ./public
RUN apk add --no-cache tzdata
EXPOSE 8080
CMD ./${NAME}


#FROM golang:1.13-alpine
#
#WORKDIR "/app"
#
#RUN apk update \
#    && apk add git \
#    && go get github.com/cosmtrek/air \
#    && go get github.com/go-delve/delve/cmd/dlv
#CMD ["air"]