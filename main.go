package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"github.com/vbabiy/simple/simple/sfile"
	_ "github.com/vbabiy/simple/simple/store"
	"github.com/vbabiy/simple/simple/store"
)

func main() {
	if len(os.Args) >= 3 {
		component, task := os.Args[1], os.Args[2]
		component = strings.ToLower(strings.Trim(component, " "))
		task = strings.ToLower(strings.Trim(task, " "))

		if component == "server" {
			handleServer(task)
		} else if component == "simple" {
			handleSimpleFile(task, os.Args[3:])
		}
	} else {
		fmt.Println("Command is <component> <task>")

	}
}
func handleServer(task string) {

	if task == "start" {
		http.HandleFunc("/what", func(w http.ResponseWriter, r *http.Request) {
			out := `
				<h1>What?</h1>
				<ul>
				%s
				</ul>
			`
			lis := []string{}
			for _, value := range store.MetaStore.All() {
				lis = append(lis, fmt.Sprintf("<li>%s - %s</li>", value.UUID, value.Filename))
			}
			fmt.Fprintf(w, out, strings.Join(lis, ""))
		})

		http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
			store.MetaStore.Reload()
			w.Write([]byte("Done..."))
		})

		log.Println("Starging webserver...")
		log.Fatal(http.ListenAndServe(":9999", nil))
	} else {
		log.Fatal("Missing Task...")
	}
}

func handleSimpleFile(task string, args []string) {
	if task == "add" {
		if len(args) == 0 {
			log.Fatal("Missing file path.")
		}
		filename := args[0]

		log.Println("Processing", filename)

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Add to current store
		if meta, err := store.MetaStore.Add(file); err != nil {
			log.Fatal(err)
		} else {
			outputName := sfile.SwapExt(file.Name())
			sfile.WriteSimpleFile(outputName, meta)
			log.Println(outputName, "Has been added to simple, Thank you!")
		}
	}
}


