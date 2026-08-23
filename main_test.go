package main

import (
	"bytes"
	"embed"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cfg, err := parse([]string{"art-init", "-title", "My Art", "my-canvas"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if cfg.dir != "my-canvas" {
		t.Errorf("dir = %q, want %q", cfg.dir, "my-canvas")
	}

	if cfg.title != "My Art" {
		t.Errorf("title = %q, want %q", cfg.title, "My Art")
	}

	if cfg.hostname != defaultHostname {
		t.Errorf("hostname = %q, want %q", cfg.hostname, defaultHostname)
	}

	if cfg.serverPath != defaultServerPath {
		t.Errorf("serverPath = %q, want %q", cfg.serverPath, defaultServerPath)
	}

	if cfg.shaders {
		t.Error("shaders = true, want false")
	}

	if cfg.zon.name != ".my_canvas" {
		t.Errorf("zon.name = %q, want %q", cfg.zon.name, ".my_canvas")
	}
}

func TestParseTitleDefaultsToDir(t *testing.T) {
	cfg, err := parse([]string{"art-init", "canvas"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if cfg.title != "canvas" {
		t.Errorf("title = %q, want %q", cfg.title, "canvas")
	}
}

func TestParseShadersAppendsToServerPath(t *testing.T) {
	cfg, err := parse([]string{"art-init", "-shaders", "canvas"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if !cfg.shaders {
		t.Error("shaders = false, want true")
	}

	if want := defaultServerPath + "/shaders"; cfg.serverPath != want {
		t.Errorf("serverPath = %q, want %q", cfg.serverPath, want)
	}
}

func TestParseRequiresDirectory(t *testing.T) {
	if _, err := parse([]string{"art-init"}, io.Discard); err == nil {
		t.Error("parse without directory should fail")
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"mycanvas", "mycanvas"},
		{"my-art", "my_art"},
		{"my.canvas", "my_canvas"},
		{"MyArt", "MyArt"},
		{"9lives", "_9lives"},
		{"hello world", "hello_world"},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", strings.Repeat("a", 32)},
	}

	for _, tt := range tests {
		if got := sanitizeName(tt.in); got != tt.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGenerateZon(t *testing.T) {
	z := generateZon("my-art")

	if z.name != ".my_art" {
		t.Errorf("name = %q, want %q", z.name, ".my_art")
	}

	wantPrefix := crc32.ChecksumIEEE([]byte("my_art"))

	var got uint64

	if _, err := fmt.Sscanf(z.fingerprint, "0x%x", &got); err != nil {
		t.Fatalf("fingerprint %q: %v", z.fingerprint, err)
	}

	if id := got & 0xffffffff; id == 0x00000000 || id == 0xffffffff {
		t.Errorf("fingerprint %q uses reserved id %#x", z.fingerprint, id)
	}

	if uint32(got>>32) != wantPrefix {
		t.Errorf("fingerprint %q checksum = %#x, want %#x", z.fingerprint, got>>32, wantPrefix)
	}
}

func TestReplacer(t *testing.T) {
	cfg, err := parse([]string{"art-init", "-title", "My Art", "my-canvas"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	makefile := replacer(cfg, "Makefile", []byte("NAME=art-canvas\nHOSTNAME=localhost\nSERVER_PATH=~/public_html/${NAME}\n"))
	for _, want := range []string{"NAME=my-canvas", "HOSTNAME=art.c7.se", "SERVER_PATH=/var/www/art.c7.se"} {
		if !bytes.Contains(makefile, []byte(want)) {
			t.Errorf("Makefile missing %q:\n%s", want, makefile)
		}
	}

	index := replacer(cfg, "index.html", []byte("<title>art-canvas-title</title>"))
	if !bytes.Contains(index, []byte("<title>My Art</title>")) {
		t.Errorf("index.html title not replaced:\n%s", index)
	}

	shaderCfg := cfg
	shaderCfg.shaders = true
	shaderCfg.serverPath += "/shaders"

	shaderIndex := replacer(shaderCfg, "index.html", []byte(`<title>art-canvas-title</title><script>const bin = "zig-out/bin/webgl.wasm";</script>`))
	for _, want := range []string{"<title>My Art</title>", "zig-out/bin/my-canvas.wasm"} {
		if !bytes.Contains(shaderIndex, []byte(want)) {
			t.Errorf("shader index.html missing %q:\n%s", want, shaderIndex)
		}
	}

	zonData := replacer(cfg, "build.zig.zon", []byte(".name = .art_canvas_name,\n.fingerprint = 0x7f6ba5038cf6243c,"))
	if !bytes.Contains(zonData, []byte(cfg.zon.name)) || !bytes.Contains(zonData, []byte(cfg.zon.fingerprint)) {
		t.Errorf("build.zig.zon not replaced:\n%s", zonData)
	}
}

func TestReplacerPanicsOnMissingPlaceholder(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("replacer should panic when placeholder is missing")
		}
	}()

	cfg, _ := parse([]string{"art-init", "canvas"}, io.Discard)
	replacer(cfg, "index.html", []byte("<title>untouched</title>"))
}

func TestTemplatesScaffoldWithoutDrift(t *testing.T) {
	cfg, err := parse([]string{"art-init", "-title", "Drift Check", "drift-check"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	shaderCfg := cfg
	shaderCfg.shaders = true
	shaderCfg.serverPath += "/shaders"

	for _, tt := range []struct {
		fsys embed.FS
		base string
		cfg  config
	}{
		{content, "content", cfg},
		{shaders, "shaders", shaderCfg},
	} {
		err := fs.WalkDir(tt.fsys, tt.base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			data, err := tt.fsys.ReadFile(path)
			if err != nil {
				return err
			}

			name := strings.TrimPrefix(path, tt.base+"/")

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s: %v", name, r)
				}
			}()

			replacer(tt.cfg, name, data)

			return nil
		})

		if err != nil {
			t.Error(err)
		}
	}
}
