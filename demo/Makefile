all: main.wasm

main.wasm: main.go ../*.go
	GOOS=js GOARCH=wasm go build -o main.wasm main.go

clean:
	rm -f main.wasm
