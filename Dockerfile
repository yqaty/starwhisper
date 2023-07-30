from golang as builder
COPY . /src
ENV GOPROXY "https://goproxy.cn"
RUN go build -o /build/starwhisper main.go

FROM ubuntu as prod
COPY --from=builder /build/starwhisper /usr/bin/starwhisper
WORKDIR /app
ENTRYPOINT [ "starwhisper" ]
