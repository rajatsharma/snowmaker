package main

import (
	"fmt"
	"os"
	"path"
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

func main() {
	if exists(relativePath("stack.yaml")) {
		println("haskell project found")
		return
	}

	if exists(relativePath("go.mod")) {
		println("go project found")
		return
	}

	if exists(relativePath("Cargo.toml")) {
		println("rust project found")
		return
	}

	if exists(relativePath("package.json")) && exists(relativePath("pnpm-lock.yaml")) {
		println("node pnpm project found")
		return
	}

	println("unknown project")
}
