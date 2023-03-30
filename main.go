package main

import (
	"cligpt/cligpt"
	"os"
)

func main() {
	cligpt.CLI(os.Args[1:])
}
