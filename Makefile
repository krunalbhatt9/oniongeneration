help:
	echo "please run make build, make run or make clean"

build:
	go build -o bin/onionclient client.go
	go build -o bin/onionserver server.go
run:
	./bin/onionserver -router 0 &
	./bin/onionserver -router 1 &
	./bin/onionserver -router 2 &
	./bin/onionclient
clean:
	rm -f ./bin/onionserver
	rm -f ./bin/onionclient
	pkill "onionserver"