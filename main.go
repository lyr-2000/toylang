package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lyr-2000/toylang/base/evaluator"
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

	vm := evaluator.NewCodeRunner()
	exit,err := vm.ParseAndRunRecover(string(code))
	if err != nil {
		log.Fatalf("panic error: %v", err)
	}
	if exit != 0 {
		log.Fatalf("exit code: %d", exit)
	}
}
