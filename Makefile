default: build

setup:
	go get github.com/google/uuid
	go get github.com/gorilla/mux
	go get github.com/tehcyx/gochain/blockchain

build: test cover
	go build -i -o bin/app

docker:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/appdocker
	docker build -t blockchain .

run:
	docker run --rm -p 8080:8080 blockchain

test:
	go test ./...

cover:
	go test ./... -cover

scenario:
	docker network create gochain || true
	docker run --rm -p 8080:8080 -d --net gochain --name blockchain1 --link blockchain2 blockchain
	docker run --rm -p 8081:8080 -d --net gochain --name blockchain2 --link blockchain1 blockchain

clean:
	rm -rf bin