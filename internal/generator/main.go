package main

import (
	"flag"
	"io"
	"log"

	"github.com/Newbluecake/bootstrap/clang"
)

func main() {
	verbose := flag.Bool("v", false, "verbose trace logging")
	flag.Parse()

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	log.Println("Bindings generator")
	log.Printf("libclang: %s", clang.GetClangVersion())
	log.Printf("platform args: %v", getPlatformArgs())
	log.Printf("system includes: %v", getSystemIncludes())

	m := Parse()

	m.structs["AVRational"].ByValue = true
	m.enums["AVOptionType"].Comment = ""

	Gen(m)
}
