package main

import "fmt"

func main() {
	input := `
package main

import ("fmt")

func main() {
	fmt.Println("Hello world")
}
var x = [...]string{}
	`
	obf := NewObfuscator(input)
	fmt.Println(obf.Obfuscate().String())
}
