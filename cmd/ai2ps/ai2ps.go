package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/fpagyu/illustrator"
)

var (
	input  = flag.String("i", "", "-i <input file path>")
	output = flag.String("o", "", "-o <output file path>")
)

func main() {
	flag.Parse()

	if len(*output) == 0 {
		*output = strings.TrimSuffix(*input, ".ai") + "-ai.ps"
	}

	r, err := illustrator.NewFileReader(*input)
	if err != nil {
		log.Fatal(err)
	}

	priData, err := r.GetAIPrivateData()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(priData)
}
