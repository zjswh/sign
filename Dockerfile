FROM golang:alpine AS development

WORKDIR $GOPATH/src

ADD . $GOPATH/src/app

RUN export GOPROXY=https://goproxy.cn && cd $GOPATH/src/app && go build

FROM xavierror/go_alpine AS production

COPY --from=development /go/src/app /app

# 按情况修改
EXPOSE 8000

RUN chmod +x /app/sign

ENTRYPOINT ["/app/sign"]
