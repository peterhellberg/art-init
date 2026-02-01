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
	zig build --watch

.PHONY: run
run: build
	@command -v npx >/dev/null 2>&1 || { echo >&2 "missing npx, install npm"; exit 1; }
	npx --yes live-server . --ignore="src,.git,.zig-cache" --wait=0ms --watch=.

.PHONY: clean
clean:
	rm -rf .zig-cache
	rm -rf zig-out

.PHONY: deploy
deploy: build
	@ssh ${HOSTNAME} 'mkdir -p ${SERVER_PATH}/${BIN_PATH} ${SERVER_PATH}/src'
	@scp -q build.zig ${HOSTNAME}:${SERVER_PATH}/build.zig
	@scp -q build.zig.zon ${HOSTNAME}:${SERVER_PATH}/build.zig.zon
	@scp -q index.html ${HOSTNAME}:${SERVER_PATH}/index.html
	@scp -q script.js ${HOSTNAME}:${SERVER_PATH}/script.js
	@scp -q src/canvas.zig ${HOSTNAME}:${SERVER_PATH}/src/canvas.zig
	@scp -q ${BIN_PATH}/${NAME}.wasm ${HOSTNAME}:${SERVER_PATH}/${BIN_PATH}/${NAME}.wasm
	@echo "✔ Updated ${NAME} on ${ART_URL}"
