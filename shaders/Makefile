NAME=art-canvas
ART_URL=https://${HOSTNAME}/shaders/${NAME}
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
	spy --exc zig-cache --inc "src/*" --exc "**/*.wasm" -q clear-zig build

.PHONY: run
run: build
	npx live-server . --ignore="src,.git,.zig-cache" --wait=0ms --watch=.

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
	@scp -q src/fragment.glslx ${HOSTNAME}:${SERVER_PATH}/src/fragment.glslx
	@scp -q src/vertex.glslx ${HOSTNAME}:${SERVER_PATH}/src/vertex.glslx
	@scp -q src/webgl.zig ${HOSTNAME}:${SERVER_PATH}/src/webgl.zig
	@scp -q ${BIN_PATH}/${NAME}.wasm ${HOSTNAME}:${SERVER_PATH}/${BIN_PATH}/${NAME}.wasm
	@echo "✔ Updated ${NAME} on ${ART_URL}"
