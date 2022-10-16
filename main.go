package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/gobeam/stringy"
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
	Language string
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
	if err = tmpl.Execute(buf, Flake{Language: stringy.New(language).CamelCase()}); err != nil {
		panic("unable to create flake.nix")
	}

	if err = os.WriteFile(relativePath("flake.nix"), []byte(buf.String()), 0644); err != nil {
		panic("unable to write flake.nix")
	}

	if err = os.WriteFile(relativePath(".envrc"), []byte("use flake"), 0644); err != nil {
		panic("unable to write .envrc")
	}

	gitExclude, err := os.OpenFile(relativePath(".git/info/exclude"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}

	defer gitExclude.Close()

	if _, err = gitExclude.WriteString("\nflake.nix\nflake.lock\n.envrc"); err != nil {
		panic(err)
	}
}
