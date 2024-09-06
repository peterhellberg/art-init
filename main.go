package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
)

//go:embed all:content
var content embed.FS

const (
	defaultHostname   = "art.c7.se"
	defaultServerPath = "/var/www/art.c7.se"
)

type config struct {
	dir        string
	title      string
	hostname   string
	serverPath string
}

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, _ io.Writer) error {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.Usage = func() {
		format := "Usage: %s [OPTION]... DIRECTORY\n\nOptions:\n"

		fmt.Fprintf(flags.Output(), format, os.Args[0])

		flags.PrintDefaults()
	}

	flags.StringVar(&cfg.title, "title", "", "The title of the ART project")
	flags.StringVar(&cfg.hostname, "hostname", defaultHostname, "The hostname to deploy the ART canvas to")
	flags.StringVar(&cfg.serverPath, "server-path", defaultServerPath, "The path on the server ART should be uploaded to")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return fmt.Errorf("no name given as the first argument")
	}

	cfg.dir = rest[0]

	if cfg.title == "" {
		cfg.title = cfg.dir
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

	entries, err := content.ReadDir("content")
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
				srcEntries, err := content.ReadDir("content/src")
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
	defer func() { f.Close() }()

	return err
}

func writeFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := content.ReadFile("content/" + name)
	if err != nil {
		return fmt.Errorf("writeFile: %w", err)
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
		return replaceOne(data, "art-canvas-title", cfg.title)
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	return bytes.Replace(data, []byte(old), []byte(new), 1)
}
