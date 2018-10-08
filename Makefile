all: main.wasm

main.wasm: main.go
	GOOS=js GOARCH=wasm go build -o main.wasm main.go

clean:
	rm main.wasm
