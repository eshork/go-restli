package main

import (
	"log"

	"github.com/PapaCharlie/go-restli/internal/codegen/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := cmd.CodeGenerator().Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
