package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
)

//go:embed all:content
var content embed.FS

//go:embed all:shaders
var shaders embed.FS

const (
	defaultHostname   = "art.c7.se"
	defaultServerPath = "/var/www/art.c7.se"
	defaultShaders    = false
)

type config struct {
	dir        string
	title      string
	hostname   string
	serverPath string
	shaders    bool
	zon        ZON
}

type ZON struct {
	name        string
	fingerprint string
}

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func parse(args []string, stderr io.Writer) (config, error) {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.SetOutput(stderr)

	flags.Usage = func() {
		format := "Usage: %s [OPTION]... DIRECTORY\n\nOptions:\n"

		fmt.Fprintf(flags.Output(), format, os.Args[0])

		flags.PrintDefaults()
	}

	flags.StringVar(&cfg.title, "title", "", "The title of the ART project")
	flags.StringVar(&cfg.hostname, "hostname", defaultHostname, "The hostname to deploy the ART canvas to")
	flags.StringVar(&cfg.serverPath, "server-path", defaultServerPath, "The path on the server ART should be uploaded to")
	flags.BoolVar(&cfg.shaders, "shaders", defaultShaders, "Should the shaders template be used or not")

	if err := flags.Parse(args[1:]); err != nil {
		return cfg, err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return cfg, fmt.Errorf("no name given as the first argument")
	}

	cfg.dir = rest[0]

	if cfg.title == "" {
		cfg.title = cfg.dir
	}

	if cfg.shaders {
		cfg.serverPath = filepath.Join(cfg.serverPath, "shaders")
	}

	cfg.zon = generateZON(cfg.dir)

	return cfg, nil
}

func run(args []string, stderr io.Writer) error {
	cfg, err := parse(args, stderr)
	if err != nil {
		return err
	}

	// Make sure that dir does not already exist
	if _, err := os.Stat(cfg.dir); !os.IsNotExist(err) {
		return fmt.Errorf("%q already exists", cfg.dir)
	}

	// Create the dir and dir/src
	if err := os.MkdirAll(cfg.dir+"/src", os.ModePerm); err != nil {
		return err
	}

	// Enter the new directory
	if err := os.Chdir(cfg.dir); err != nil {
		return err
	}

	var (
		writeFile = contentWriteFile
		srcFS     = content
		srcBase   = "content"
	)

	if cfg.shaders {
		writeFile = shadersWriteFile
		srcFS = shaders
		srcBase = "shaders"
	}

	entries, err := srcFS.ReadDir(srcBase)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			if err := writeFile(cfg, e.Name(), replacer); err != nil {
				return err
			}
		} else {
			if e.Name() == "src" {
				srcEntries, err := srcFS.ReadDir(srcBase + "/src")
				if err != nil {
					return err
				}

				for _, e := range srcEntries {
					if !e.IsDir() {
						if err := writeFile(cfg, "src/"+e.Name(), replacer); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	if err := createFile("favicon.ico"); err != nil {
		return err
	}

	return nil
}

func createFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
}

type writeFileFunc func(cfg config, name string, dataFuncs ...dataFunc) error

func contentWriteFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := content.ReadFile("content/" + name)
	if err != nil {
		return fmt.Errorf("contentWriteFile: %w", err)
	}

	for i := range dataFuncs {
		data = dataFuncs[i](cfg, name, data)
	}

	return os.WriteFile(name, data, 0o644)
}

func shadersWriteFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := shaders.ReadFile("shaders/" + name)
	if err != nil {
		return fmt.Errorf("shadersWriteFile: %w", err)
	}

	for i := range dataFuncs {
		data = dataFuncs[i](cfg, name, data)
	}

	return os.WriteFile(name, data, 0o644)
}

type dataFunc func(config, string, []byte) []byte

func replacer(cfg config, name string, data []byte) []byte {
	switch name {
	case "Makefile":
		data = replaceOne(data, `NAME=art-canvas`, `NAME=`+cfg.dir)
		data = replaceOne(data, `HOSTNAME=localhost`, `HOSTNAME=`+cfg.hostname)
		data = replaceOne(data, `SERVER_PATH=~/public_html`, `SERVER_PATH=`+cfg.serverPath)
		return data
	case "README.md", "build.zig", "script.js":
		return replaceOne(data, "art-canvas", cfg.dir)
	case "index.html":
		data = replaceOne(data, "art-canvas-title", cfg.title)

		if cfg.shaders {
			data = replaceOne(data, "zig-out/bin/webgl.wasm", fmt.Sprintf("zig-out/bin/%s.wasm", cfg.dir))
		}

		return data
	case "build.zig.zon":
		data = replaceOne(data, ".art_canvas_name", cfg.zon.name)
		data = replaceOne(data, "0x7f6ba5038cf6243c", cfg.zon.fingerprint)

		return data
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	if !bytes.Contains(data, []byte(old)) {
		panic(fmt.Sprintf("art-init: placeholder %q missing from template", old))
	}

	return bytes.Replace(data, []byte(old), []byte(new), 1)
}

func generateZON(dir string) ZON {
	name := sanitizeName(dir)

	var id uint32

	for id == 0x00000000 || id == 0xffffffff {
		id = rand.Uint32()
	}

	return ZON{
		name:        "." + name,
		fingerprint: fmt.Sprintf("0x%08x%08x", crc32.ChecksumIEEE([]byte(name)), id),
	}
}

func sanitizeName(dir string) string {
	b := []byte(dir)

	for i, c := range b {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '_':
		default:
			b[i] = '_'
		}
	}

	if len(b) > 0 && b[0] >= '0' && b[0] <= '9' {
		b = append([]byte{'_'}, b...)
	}

	if len(b) > 32 {
		b = b[:32]
	}

	return string(b)
}
