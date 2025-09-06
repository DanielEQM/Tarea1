package main

import (
	"fmt"
)

func main() {
	fmt.Print("Hola mundo")
	var aux int
	for i := 0; i < 5; i++ {
		aux = aux + 1
	}
	fmt.Print(aux)
}
