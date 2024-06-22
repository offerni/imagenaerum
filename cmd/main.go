package main

import (
	"github.com/offerni/imagenaerum/rest"
	"github.com/offerni/imagenaerum/utils"
)

func main() {
	ensureDirectories()
	rest.InitializeServer()
}

func ensureDirectories() {
	for _, dir := range utils.Directories {
		utils.EnsureDir(dir)
	}
}
