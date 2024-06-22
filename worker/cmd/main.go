package main

import (
	"github.com/offerni/imagenaerum/rest"
	"github.com/offerni/imagenaerum/utils"
)

func main() {
	utils.EnsureDirectories()
	rest.InitializeServer()
}
