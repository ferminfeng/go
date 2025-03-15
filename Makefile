.PHONY: buildGame,runGame

#buildGame:
#	CGO_ENABLED=1 go build -ldflags="-w -s" -o build/testClint ./cmd/main.go

runGame:
	cd ./game/op_symbol && CGO_ENABLED=1 go run main.go
