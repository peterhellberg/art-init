# art-init :sparkles:

`art-init` is used to create a directory containing code that
allows you to promptly get started with using [Zig](https://ziglang.org/)
to draw into a [Canvas](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API).

The Zig build `.target` is declared as `.{ .cpu_arch = .wasm32, .os_tag = .freestanding }`
and `.optimize` is set to `.ReleaseSmall`

> [!Important]
> No need to specify `-Doptimize=ReleaseSmall`

## Installation

(Requires you to have [Go](https://go.dev/) installed)

```sh
go install github.com/peterhellberg/art-init@latest
```

> [!Note]
> Live reloading when using `make run` is handled by
> <https://github.com/peterhellberg/live> via `go run`.

## Usage

(Requires you to have an up to date (_nightly_) version of
[Zig](https://ziglang.org/download/#release-master) installed.

```sh
art-init mycanvas
cd mycanvas
make run
```

> [!Tip]
> If you want to write WebGL shaders instead, then use the `-shaders` flag
>
> ```sh
> art-init -shaders myshaders
> cd myshaders
> make run
> ```

:seedling:
