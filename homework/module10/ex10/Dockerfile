FROM golang:1.17 AS build 
WORKDIR /app/
COPY main.go go.mod go.sum ./
COPY metrics/metrics.go ./metrics/
ENV GO111MODULE=on \
    CGO_ENABLED=0 \	
    GOOS=linux \	
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct
RUN go build -o httpserver .

FROM busybox
COPY --from=build /app/httpserver .
EXPOSE 80
ENTRYPOINT ["/httpserver"]
