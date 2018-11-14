test:
	@echo "[TEST]"
	GOOS=js GOARCH=wasm go test -exec="${HOME}/devel/go/misc/wasm/go_js_wasm_exec"
build:
	@echo "[BUILD]"
	GOOS=js GOARCH=wasm go build