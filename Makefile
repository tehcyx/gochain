default: builddocker

setup:
	go get github.com/google/uuid
	go get github.com/gorilla/mux
	go get github.com/tehcyx/gochain/blockchain

buildgo: setup
	go build -i

builddocker: setup
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o gochaindocker
	docker build -t blockchain .

run:
	docker run --rm -p 8080:8080 blockchain

scenario:
	docker network create gochain || true
	docker run --rm -p 8080:8080 -d --net gochain --name blockchain1 --link blockchain2 blockchain
	docker run --rm -p 8081:8080 -d --net gochain --name blockchain2 --link blockchain1 blockchain