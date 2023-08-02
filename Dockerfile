from golang as builder
COPY . /cmd
WORKDIR /cmd
ENV GOPATH /
ENV GOPROXY "https://goproxy.cn"
RUN go build -o /build/starwhisper main.go
  
FROM ubuntu as prod
COPY --from=builder /build/starwhisper /usr/bin/starwhisper
COPY --from=builder /cmd/config.json /app/config.json
WORKDIR /app
ENTRYPOINT [ "starwhisper" ]
