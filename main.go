package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const namePostfix = "fields"

var (
	tName = flag.String("tag", "field",
		"tag name whose value will be used as the name of the field")
	sNames = flag.String("struct", "",
		"comma-separated list of struct names for which to generate fields; must be set")
	cNames = flag.String("custom_name", "",
		"comma-separated list of names, default is the struct names for which to generate fields")
	output = flag.String("output", "",
		"output file name; default file creates for each struct, with name: <struct_name>_fields.go")
)

// usage is a replacement usage function for the flags package.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tgen-struct-fields -struct=StructName [flags]\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/abramlab/gen-struct-fields\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

type parsedFlags struct {
	tagName       string
	neededStructs map[string]*options
	output        string
}

type options struct {
	customName string
}

func parseFlags() *parsedFlags {
	flag.Usage = usage
	flag.Parse()

	if *sNames == "" {
		flag.Usage()
		os.Exit(2)
	}
	structNames := strings.Split(*sNames, ",")
	customNames := make([]string, len(structNames))
	if *cNames != "" {
		cns := strings.Split(*cNames, ",")
		for i, name := range cns {
			customNames[i] = name
		}
	}

	neededStructs := make(map[string]*options)
	for i := 0; i < len(structNames); i++ {
		customName := structNames[i]
		if customNames[i] != "" {
			customName = customNames[i]
		}
		neededStructs[structNames[i]] = &options{
			customName: customName,
		}
	}
	return &parsedFlags{
		tagName:       *tName,
		neededStructs: neededStructs,
		output:        *output,
	}
}

func main() {
	flags := parseFlags()

	gen := &generator{
		tagName: flags.tagName,
		genTpls: basicTemplates,
	}

	if err := gen.parse(flags.neededStructs); err != nil {
		log.Fatalf("parsing files failed: %v", err)
	}

	if err := gen.generate(flags.output); err != nil {
		log.Fatalf("generating files failed: %v", err)
	}
}
