FROM golang:1.17 AS build 
WORKDIR /app/
COPY main.go go.mod go.sum ./
ENV GO111MODULE=on \
    CGO_ENABLED=0 \	
    GOOS=linux \	
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct
RUN go build -o service1 .

FROM busybox
COPY --from=build /app/service1 .
EXPOSE 80
ENTRYPOINT ["/service1"]
