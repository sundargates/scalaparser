package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sundargates/scalaparser/parser"
)

func lexFile(filename string) error {
	// fmt.Println("Processing ", filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	input := string(bytes)
	// fmt.Println(input)

	lexer := parser.Lexer(input)
	tokens := lexer.LexTillDone()
	for _, token := range tokens {
		// fmt.Println(token.String())
		if token.Typ == parser.ERROR {
			return errors.New(token.String())
		}
	}
	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	for _, arg := range args {
		name, err := filepath.Abs(arg)
		if err != nil {
			fmt.Printf("%s: %s", arg, err)
			continue
		}
		filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
			if _, elem := filepath.Split(path); elem != "" {
				// Skip various temporary or "hidden" files or directories.
				if elem[0] == '.' || elem[0] == '#' || elem[0] == '~' || elem[len(elem)-1] == '~' {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			if err != nil {
				fmt.Printf("%s: %s", path, err)
				return nil
			}
			// fmt.Println(info.Mode(), os.ModeT)
			if info != nil && info.Mode()&os.ModeType == 0 && filepath.Ext(path) == ".scala" {
				err := lexFile(path)
				if err != nil {
					fmt.Println("Failed processing ", path, "with error", err.Error())
				}
			}
			return nil
		})
	}
}
