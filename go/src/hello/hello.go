// +build cgo

package main

import "C"

import "fmt"

//export hello
func hello() {
	fmt.Println("A hello from go!")
}
