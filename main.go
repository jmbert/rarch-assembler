package main

import (
	"fmt"
	"log"
	"os"

	rarch_parser "github.com/jmbert/rarch-assembler/parser"
)

func Assemble(file string) ([]byte, error) {

	fcontents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var parser rarch_parser.Parser

	parser.Buffer = string(fcontents) + "\n"
	parser.Pretty = true

	err = parser.Init()
	if err != nil {
		return nil, err
	}
	err = parser.Parse()
	if err != nil {
		return nil, err
	}

	bytes := parser.Codegen()

	return bytes, nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Must provide input and output file")
		os.Exit(1)
	}

	bytes, err := Assemble(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	os.Remove(os.Args[2])
	file, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
