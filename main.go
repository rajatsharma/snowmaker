package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

func relativePath(file string) string {
	currentDir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path.Join(currentDir, file)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	panic(fmt.Sprintf("error while reading %s occured: %v", path, err))
}

//go:embed flake.txt
var flake string

type Flake struct {
	language string
}

func main() {
	var language string
	if exists(relativePath("stack.yaml")) {
		println("haskell project found")
		language = "haskell"
	}

	if exists(relativePath("go.mod")) {
		println("go project found")
		language = "go"
	}

	if exists(relativePath("Cargo.toml")) {
		println("rust project found")
		language = "rust"
	}

	if exists(relativePath("package.json")) && exists(relativePath("pnpm-lock.yaml")) {
		println("node pnpm project found")
		language = "node"
	}

	if language == "" {
		println("unsupported language found")
		os.Exit(1)
	}

	tmpl, err := template.New("nix").Parse(flake)

	if err != nil {
		panic("unable to parse template")
	}

	buf := new(bytes.Buffer)
	tmpl.Execute(buf, Flake{language: language})

	err = os.WriteFile(relativePath("flake.nix"), []byte(buf.String()), 0644)

	if err != nil {
		panic("unable to write flake.nix")
	}

	err = os.WriteFile(relativePath(".envrc"), []byte("use flake"), 0644)

	if err != nil {
		panic("unable to write .envrc")
	}
}
