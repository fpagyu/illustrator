package main

import (
	"bytes"
	"flag"
	"log"
	"strings"

	"github.com/fpagyu/illustrator"
	"github.com/fpagyu/illustrator/svg"
)

var (
	input  = flag.String("i", "", "-i <input file path>")
	output = flag.String("o", "", "-o <output file path>")
)

func main() {
	flag.Parse()

	if len(*input) == 0 {
		*input = "/Users/lemi/Downloads/029.ai"
		// *input = "/Users/lemi/Downloads/Zebra_giraffe_brown.ai"
		// *input = "/Users/lemi/Downloads/spring-flower-collection/4947645.ai"
	}

	if len(*output) == 0 {
		*output = strings.TrimSuffix(*input, ".ai") + "-ai.svg"
	}

	r, err := NewReader(*input)
	if err != nil {
		log.Fatal(err)
	}

	var svg svg.SVG
	err = r.Draw(&svg)
	if err != nil {
		log.Fatal(err)
	}

	err = svg.Save(*output)
	if err != nil {
		log.Fatal(err)
	}
}

func NewReader(path string) (*illustrator.AIReader, error) {
	r, err := illustrator.NewFileReader(path)
	if err != nil {
		return nil, err
	}

	data, err := r.GetAIPrivateData()
	if err != nil {
		return nil, err
	}

	return illustrator.NewAIReader(bytes.NewBuffer(data))
}
