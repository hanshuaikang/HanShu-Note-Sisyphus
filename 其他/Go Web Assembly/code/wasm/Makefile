.PHONY: action, build, copy, serve

GO := "" # 替换成你的go sdk实际位置

# 用于编译 Go 代码为 WebAssembly 格式，并将结果输出到 assets/main.wasm 文件中
build:
	GOOS=js GOARCH=wasm $(GO) build -o assets/main.wasm main.go

# 用于将 wasm_exec.js 复制到 assets/ 目录中，以便在浏览器中运行 Go WebAssembly 应用
copy:
	cp -a `$(GO) env GOROOT`/misc/wasm/wasm_exec.js assets/

action: build copy

