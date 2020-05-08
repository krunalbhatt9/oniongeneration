hello:
	echo "Hello"

build:
	go build -o bin/client client.go
	go build -o bin/server server.go
run:
	nohup ./bin/server -router 0 &
	nohup ./bin/server -router 1 &
	./bin/client