package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jontynewman/tpleng"
)

func main() {

	template, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	parsed, err := tpleng.Parse(string(template), map[string]string{"num": "42", "greetee": "Good Growth", "expression": "{{ .recurse }}", "recurse": "{{ .num }}"})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", parsed)
}
