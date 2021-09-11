FROM golang:1.17

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/drvspw/go-https-server
WORKDIR /go/src/github.com/drvspw/go-https-server

RUN mkdir -p /etc/go-https-server
RUN make tools
RUN make build

CMD if [ "${APP_ENV}" = "production" ]; then \
      /go/bin/go-https-server-linux-amd64; \
    else \
      go get github.com/pilu/fresh && fresh; \
    fi

EXPOSE 8090
