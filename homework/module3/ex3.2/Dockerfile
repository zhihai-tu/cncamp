FROM golang:1.17 AS build 
WORKDIR /app/
COPY main.go .
ENV GO111MODULE=off \
    CGO_ENABLED=0 \	
    GOOS=linux \	
    GOARCH=amd64
RUN go build -o httpserver .

FROM scratch
COPY --from=build /app/httpserver /
EXPOSE 80
ENTRYPOINT ["/httpserver"]
