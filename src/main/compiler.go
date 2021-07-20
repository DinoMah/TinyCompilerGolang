package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		filePath := os.Args[1] //Obtiene la ruta del archivo
		fmt.Println("Argumentos: ", filePath)
		lexicAnalysis(filePath)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func lexicAnalysis(filePath string) {
	fmt.Println("Analyzing lexically...")
	f, err := os.Open(filePath)
	check(err)
	reader := bufio.NewReader(f)
	line, isPrefx, err := reader.ReadLine()
	aux := []byte{}
	_ = isPrefx
	var tokens []info
	for err == nil {
		line = append(line, aux...)
		aux, isPrefx, err = reader.ReadLine()
		aux = append(aux, []byte("\n")...)
	}
	tokens = analyze(string(line))
	for i := 0; i < len(tokens); i++ {
		fmt.Println("Token: ", tokens[i].token, ", Tipo: ", tokens[i].tokenType)
	}
	f.Close()
}
