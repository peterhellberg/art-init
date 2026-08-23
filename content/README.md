# art-canvas :art:

Using [Zig](https://ziglang.org/) to compile a `.wasm` cart
for use in a [Canvas](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API)

## Development

`make dev` (and the default `make all`) requires
[tmux](https://github.com/tmux/tmux) to run the watcher and server side by side.
Without tmux, run `make watch` and `make run` in separate terminals.

File watcher can be started by calling:
```sh
make watch
```

Running the canvas:
```sh
make run
```

Deploy:
```
make deploy
```

> [!Note]
> `script.js` uses `WebAssembly.instantiateStreaming`, which requires the
> `.wasm` file to be served with the `Content-Type: application/wasm` header.
