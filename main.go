package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	switch {
	case len(os.Args) < 2:
		goto isError
	default:
		if _, err := os.Stat(os.Args[1]); err == nil {
			goto isGood
		}
	}
isError:
	fmt.Println("Provide a valid go file to obfuscate...")
	os.Exit(0)

isGood:
	fileName := os.Args[1]
	outName := strings.Split(fileName, ".")[0] + "_out"

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("There was an error reading your file")
		os.Exit(0)
	}
	obfuscator := NewObfuscator(string(fileContent))
	data := obfuscator.Obfuscate().Bytes()
	os.WriteFile(outName, data, fs.ModePerm)
}
