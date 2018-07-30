FROM golang:alpine as build_go_env
RUN mkdir -p /go/src/github.com/fasibio/subscriber && mkdir -p /src/bin &&  apk update && apk add git
ADD . /go/src/github.com/fasibio/subscriber

RUN cd /go/src/github.com/fasibio/subscriber && go get && cd /go/src/github.com/fasibio/subscriber && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/subscriber

FROM alpine:3.5
ARG version
ARG buildNumber
RUN mkdir /app && apk update && apk add tzdata
COPY --from=build_go_env /src/bin/subscriber /app/
ENV VERSION=${version}
ENV BUILD_NUMBER=${buildNumber}
WORKDIR /app 
EXPOSE 8001
CMD ["/app/subscriber"]