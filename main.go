package main

import (
	_ "embed"

	"github.com/axelrindle/limacity-dns-update/cmd"
)

//go:embed banner.txt
var banner string

func main() {
	println(banner)
	println("  " + cmd.GetVersion())
	println()

	cmd.Run()
}
