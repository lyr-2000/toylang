package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "github.com/lyr-2000/toylang/base/evaluator"
	interpreter "github.com/lyr-2000/toylang/base/evaluator/v2"
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
	vm := interpreter.New()
	vm.ParseAndRun(string(code))
	if vm.ExitCode != 0 {
		log.Fatal("exit code:", vm.ExitCode)
		return
	}
	if vm.ErrCode != 0 {
		log.Fatal("error code:", vm.ErrCode, "error message:", vm.ErrMsg)
		return
	}

}
