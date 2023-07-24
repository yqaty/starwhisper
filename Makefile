NAME = main

dbp = ./internal/mydb/

DB = $(dbp)ConnectionPool.go $(dbp)user.go $(dbp)post.go $(dbp)comment.go $(dbp)redis.go $(dbp)anonymity.go $(dbp)report.go

wbp = ./internal/myweb/

WB = $(wbp)user.go $(wbp)post.go $(wbp)comment.go $(wbp)email.go $(wbp)report.go $(wbp)middle.go $(wbp)api.go

all:clean check build run

mydb : $(DB)
	@go build -o mydb $(DB)

myweb : $(WB)
	@go build -o myweb $(WB)

build : main.go
	@go build -o $(NAME) main.go

.PHONY=clean
clean:
	@if [ -f mydb ]; then rm mydb; fi
	@if [ -f myweb ]; then rm myweb; fi
	@if [ -f main ]; then rm main; fi

image : build
	docker build -t web:1 .

check :
	@go fmt $(dbp) $(wbp)
	@go vet $(dbp) $(wbp)

run : 
	docker-compose up
