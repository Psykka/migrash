package main

import (
	"fmt"

	"migrash/cmd"
)

var version = "1.0.0"

func main() {
	fmt.Printf("Migrash v%s\n", version)
	cmd.Execute()
}
