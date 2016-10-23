package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vbabiy/simple/simple/store"
)

func StartServer(s *store.Store, address string) error {
	http.HandleFunc("/what", func(w http.ResponseWriter, r *http.Request) {
		out := `
				<h1>What?</h1>
				<ul>
				%s
				</ul>
			`
		lis := []string{}
		for _, value := range s.All() {
			lis = append(lis, fmt.Sprintf("<li>%s - %s</li>", value.UUID, value.Filename))
		}
		fmt.Fprintf(w, out, strings.Join(lis, ""))
	})

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		s.Reload()
		w.Write([]byte("Done..."))
	})

	return http.ListenAndServe(address, nil)

}
