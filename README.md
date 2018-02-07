# gochain
educational implementation of a simple blockchain in go

## How to build/run

No matter if you want to run it from docker or on your machine with the executable, run make setup first to install dependencies.
### For development on the machine (no docker)
If you want to run only one node on your machine, use the `make buildgo` to create the executable called `gochain`, run with `./gochain` from the folder and connect to the exposed endpoints via your browser.

### For development on the machine (docker)
If you want to run only one node on your machine inside a docker container, run the `make` command and then `make run`. This will start the docker container and bind the port locally to :8080 from the exposed docker port.

### Run two nodes via docker on one machine
If you want to run two nodes on your machine (both inside docker), run the `make` command and then use `make scenario` to start both. They are now both exposed on port :8080 and :8081 and you can access them both from your browser. The command will also setup a docker network named gochain and name the containers (blockchain1 and blockchain2) so for registering the node you use one of those hostnames.

## Exposed API
Running the containers you will have access to this API via your browser
 - [GET]  [/chain](http://localhost:8080/chain) for the current chain
 - [GET]  [/mine](http://localhost:8080/mine) to create a new block and include all pending transactions
 - [POST] [/nodes/register](http://localhost:8080/nodes/register) to register a new node (example below)
 - [GET]  [/nodes/resolve](http://localhost:8080/nodes/resolve) query all registered nodes to find the longest chain
 - [POST] [/transactions/new](http://localhost:8080/transactions/new) add a new transaction (example below)


## Examples

### GET /chain

#### request
```bash
curl --request GET \
  --url http://localhost:8080/chain
```

#### response
```json
{
    "chain": [
        {
            "index": 1,
            "timestamp": 1515801024,
            "transactions": [],
            "proof": 100,
            "previous_hash": "1",
            "hash": "9344c531ee070e63db6595da8412278aea178fbea14dd61b32bfde5a95605b1b"
        },
        {
            "index": 2,
            "timestamp": 1515801050,
            "transactions": [
                {
                    "sender": "0",
                    "recipient": "4e999a37-2c09-4427-b0b7-c8dfefa62b5b",
                    "amount": 1
                }
            ],
            "proof": 29031,
            "previous_hash": "9344c531ee070e63db6595da8412278aea178fbea14dd61b32bfde5a95605b1b",
            "hash": "35cdd0bf31cf3ff33bd2b31abf22ba32b1c4cceb753873f82d88963d7c08cb32"
        }
    ],
    "length": 2
}
```

### GET /mine

#### request
```bash
curl --request GET \
  --url http://localhost:8080/mine
```

#### response
```json
{
    "block": {
        "index": 2,
        "timestamp": 1515801050,
        "transactions": [
            {
                "sender": "0",
                "recipient": "4e999a37-2c09-4427-b0b7-c8dfefa62b5b",
                "amount": 1
            }
        ],
        "proof": 29031,
        "previous_hash": "9344c531ee070e63db6595da8412278aea178fbea14dd61b32bfde5a95605b1b",
        "hash": "35cdd0bf31cf3ff33bd2b31abf22ba32b1c4cceb753873f82d88963d7c08cb32"
    },
    "message": "Success: New block mined 2",
    "status_code": 201
}
```

### POST /nodes/register

#### request
```bash
curl --request POST \
  --url http://localhost:8080/nodes/register \
  --header 'Content-Type: application/json' \
  --data '{ "nodes": [ {"address": "http://blockchain2:8080","comment": "first node"} ] }'
```

#### response
```json
{
    "message": "Success: 1 node(s) successfully added",
    "nodes": {
        "blockchain2:8080": "first node"
    },
    "status_code": 201
}
```

### GET /nodes/resolve

#### request
```bash
curl --request GET \
  --url http://localhost:8081/nodes/resolve
```

#### response
```json
{
    "chain": [
        {
            "index": 1,
            "timestamp": 1515801024,
            "transactions": [],
            "proof": 100,
            "previous_hash": "1",
            "hash": "9344c531ee070e63db6595da8412278aea178fbea14dd61b32bfde5a95605b1b"
        },
        {
            "index": 2,
            "timestamp": 1515801050,
            "transactions": [
                {
                    "sender": "0",
                    "recipient": "4e999a37-2c09-4427-b0b7-c8dfefa62b5b",
                    "amount": 1
                }
            ],
            "proof": 29031,
            "previous_hash": "9344c531ee070e63db6595da8412278aea178fbea14dd61b32bfde5a95605b1b",
            "hash": "35cdd0bf31cf3ff33bd2b31abf22ba32b1c4cceb753873f82d88963d7c08cb32"
        }
    ],
    "message": "Our chain was replaced",
    "status_code": 200
}
```

### POST /transaction/new
#### request
```bash
curl --request POST \
  --url http://localhost:8080/transactions/new \
  --header 'Content-Type: application/json' \
  --data '{\n "sender": "my address",\n "recipient": "someone else'\''s address",\n "amount": 2.0\n}'
```

#### response
```json
{
    "message": "Success: Transaction will be added to block 2",
    "status_code": 201
}
```