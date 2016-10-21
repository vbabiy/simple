package http

import (
	"github.com/vbabiy/simple/simple/store"
	"fmt"
	"strings"
	"net/http"
)

func StartServer(address string) error {
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

	return http.ListenAndServe(address, nil)

}
