package main

import (
	"bufio"
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	path string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "<html><body>")
	defer fmt.Fprintln(w, "</body></html>")

	fmt.Fprintf(w, "<h1>%s</h1>\n", html.EscapeString(h.path))

	fd, err := os.Open(h.path)
	if err != nil {
		fmt.Fprintf(w, "<div><i>%s</i></div>\n", html.EscapeString(err.Error()))
		return
	}
	defer fd.Close()

	fmt.Fprintln(w, "<pre>")
	defer fmt.Fprintln(w, "</pre>")
	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		fmt.Fprintln(w, html.EscapeString(sc.Text()))
	}
}

func mains(args []string) error {
	if len(args) < 1 {
		return errors.New("FILENAME is required")
	}
	handler := &Handler{
		path: args[0],
	}
	service := &http.Server{
		Addr:           ":8000",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := service.ListenAndServe()
	closeErr := service.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}