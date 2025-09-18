package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// "github.com/lyr-2000/toylang/base/evaluator"
	"github.com/lyr-2000/toylang/base/compiler"
	"github.com/lyr-2000/toylang/base/evaluator/v2"
)

func main() {
	var (
		filePath string
	)

	flag.StringVar(&filePath, "f", "", "file path")
	flag.Parse()

	if filePath == "" {
		fmt.Println("file path is required")
		return
	}

	code, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file error:", err)
		return
	}
	// fmt.Println(string(code))
	vm := evaluator.New()
	byteCode := compiler.Compile(string(code))
	vm.SetReader(strings.NewReader(string(byteCode)))
	vm.Handle()
	if vm.ExitCode != 0 {
		log.Fatal("exit code:", vm.ExitCode)
		return
	}
	if vm.ErrCode != 0 {
		log.Fatal("error code:", vm.ErrCode, "error message:", vm.ErrMsg)
		return
	}

}
