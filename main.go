package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vbabiy/simple/simple/http"
	"github.com/vbabiy/simple/simple/sfile"
	"github.com/vbabiy/simple/simple/store"
)

func main() {
	if len(os.Args) >= 3 {
		component, task := os.Args[1], os.Args[2]
		component = strings.ToLower(strings.Trim(component, " "))
		task = strings.ToLower(strings.Trim(task, " "))
		switch component {
		case "server":
			err := handleServer(task)
			if err != nil {
				fmt.Println("An error occured while running the server", err)
				os.Exit(1)
			}
			return
		case "simple":
			err := handleSimpleFile(task, os.Args[3:])
			if err != nil {
				fmt.Println("An error occured in Simple file handling", err)
				os.Exit(1)
			}
			return
		}
	}
	usage()
}

func usage() {
	fmt.Println(`
	Command is <component> <task>

	To start the server:
		simple server start 

	To add a file:
		simple simple add <path-to-file>
`)
}

func handleServer(task string) error {
	if task != "start" {
		return fmt.Errorf("Missing Task...")
	}
	log.Println("Starging webserver...")
	return http.StartServer(":9999")
}

func handleSimpleFile(task string, args []string) error {
	if task != "add" {
		return fmt.Errorf("Simple task must be add")
	}

	if len(args) == 0 {
		return fmt.Errorf("Missing file path.")
	}
	filename := args[0]

	log.Println("Processing", filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add to current store
	meta, err := store.MetaStore.Add(file)
	if err != nil {
		return err
	}
	outputName := sfile.SwapExt(file.Name())
	sfile.WriteSimpleFile(outputName, meta)
	log.Println(outputName, "Has been added to simple, Thank you!")
	return nil
}
