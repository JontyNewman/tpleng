package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jontynewman/tpleng"
)

func main() {
	dog := struct {
		name   string
		isGood bool
	}{
		"Rex",
		true,
	}

	template, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	parsed, err := tpleng.Parse(string(template), map[string]any{"int": 42, "string": "Hello, world!", "struct": dog})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", parsed)
}
