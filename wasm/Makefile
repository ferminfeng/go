
wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o main.wasm

copyWasmExecJs:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./

