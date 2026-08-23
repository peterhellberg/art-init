NAME=art-canvas
ART_URL=https://${HOSTNAME}/${NAME}
BIN_PATH=zig-out/bin
SERVER_PATH=~/public_html/${NAME}
HOSTNAME=localhost

.PHONY: all
all: dev

.PHONY: dev
dev:
	@tmux new-session 'make watch' \; split-window 'make run'

.PHONY: build
build:
	@zig build

.PHONY: watch
watch:
	@zig build --watch --debounce 50 --error-style minimal

.PHONY: run
run: build
	@go run github.com/peterhellberg/live@latest

.PHONY: clean
clean:
	rm -rf .zig-cache
	rm -rf zig-out

.PHONY: deploy
deploy: build
	@ssh ${HOSTNAME} 'mkdir -p ${SERVER_PATH}/${BIN_PATH} ${SERVER_PATH}/src'
	@rsync -a build.zig build.zig.zon index.html script.js ${HOSTNAME}:${SERVER_PATH}/
	@rsync -a src/ ${HOSTNAME}:${SERVER_PATH}/src/
	@rsync -a ${BIN_PATH}/${NAME}.wasm ${HOSTNAME}:${SERVER_PATH}/${BIN_PATH}/
	@echo "✔ Updated ${NAME} on ${ART_URL}"
