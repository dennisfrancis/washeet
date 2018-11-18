default: build check

check:
	@echo "[CHECK] START"
	@GOOS=js GOARCH=wasm go test -exec="${HOME}/devel/go/misc/wasm/go_js_wasm_exec" \
	 && (echo "[CHECK] SUCCESS") || (echo "[CHECK] FAILED" && false) 
build:
	@echo -n "[BUILD] "
	@GOOS=js GOARCH=wasm go build && (echo "SUCCESS") || (echo "[BUILD] FAILED" && false)