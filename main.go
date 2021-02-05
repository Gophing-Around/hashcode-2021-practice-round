package main

import (
	"io/ioutil"
)

func main() {

	inputSet := readFile("./dataset/a_example")
}

func readFile(source string) string {
	in, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	return string(in)

}
