FROM golang:alpine as build_go_env
RUN mkdir -p /go/src/github.com/fasibio/superEasyPubSub && mkdir -p /src/bin &&  apk update && apk add git
ADD . /go/src/github.com/fasibio/superEasyPubSub

RUN cd /go/src/github.com/fasibio/superEasyPubSub && go get && cd /go/src/github.com/fasibio/superEasyPubSub && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/pubsub

FROM alpine:3.5
ARG version
ARG buildNumber
RUN mkdir /app 
COPY --from=build_go_env /src/bin/pubsub /app/
ENV VERSION=${version}
ENV BUILD_NUMBER=${buildNumber}
WORKDIR /app 
EXPOSE 8000
CMD ["/app/pubsub"]