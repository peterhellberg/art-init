NAME=art-canvas
ART_URL=https://${HOSTNAME}/${NAME}
BIN_PATH=zig-out/bin
SERVER_PATH=~/public_html/${NAME}
HOSTNAME=localhost

.PHONY: all
all: build

.PHONY: build
build:
	zig build

.PHONY: spy
spy:
	spy --exc zig-cache --inc "**/*.zig" -q clear-zig build

.PHONY: run
run: build
	npx live-server . --wait=0ms --watch=zig-out/bin

.PHONY: clean
clean:
	rm -rf .zig-cache
	rm -rf zig-out

.PHONY: deploy
deploy: build
	@ssh ${HOSTNAME} 'mkdir -p ${SERVER_PATH}/${BIN_PATH} ${SERVER_PATH}/src'
	@scp -q src/art.zig ${HOSTNAME}:${SERVER_PATH}/src/art.zig
	@scp -q src/canvas.zig ${HOSTNAME}:${SERVER_PATH}/src/canvas.zig
	@scp -q index.html ${HOSTNAME}:${SERVER_PATH}/index.html
	@scp -q script.js ${HOSTNAME}:${SERVER_PATH}/script.js
	@scp -q ${BIN_PATH}/${NAME}.wasm ${HOSTNAME}:${SERVER_PATH}/${BIN_PATH}/${NAME}.wasm
	@echo "✔ Updated ${NAME} on ${ART_URL}"
