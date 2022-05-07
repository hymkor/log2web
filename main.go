package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"
)

var portNo = flag.Int("p", 8000, "portNo")

type Handler struct {
	notFound http.Handler
	path     string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, req.Method, req.URL.Path)
	if req.URL.Path == "/favicon.ico" {
		// log.Println(req.RemoteAddr, req.Method, req.URL.Path, "NOT FOUND")
		h.notFound.ServeHTTP(w, req)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	title := html.EscapeString(h.path)

	fmt.Fprintln(w, `<html>
<head><style><!--
	th{
		text-align: right
	}
	pre{
		background-color: #f8f8f8
	}
--></style>`)
	fmt.Fprintf(w, "<title>%s</title>\n", title)
	fmt.Fprintln(w, "</head><body>")
	defer fmt.Fprintln(w, "</body></html>")

	fmt.Fprintf(w, "<h1>%s</h1>\n", title)

	fd, err := os.Open(h.path)
	if err != nil {
		fmt.Fprintf(w, "<div><i>%s</i></div>\n", html.EscapeString(err.Error()))
		return
	}
	defer fd.Close()

	dropLineCount := 0
	stat, err := fd.Stat()
	if err == nil {
		size := stat.Size()
		fmt.Fprintln(w, "<table>")
		fmt.Fprintf(w, "<tr><th>ModTime:</th><td>%v</td></tr>\n",
			stat.ModTime())
		fmt.Fprintf(w, "<tr><th>Size:</th><td>%d bytes",
			size)
		if size > 1024 {
			fd.Seek(-1024, os.SEEK_END)
			dropLineCount = 1
			fmt.Fprintln(w, "(showing last 1K bytes)")
		}
		fmt.Fprintln(w, "</td></tr></table>")
	}

	fmt.Fprintln(w, `<pre>`)
	defer fmt.Fprintln(w, "</pre>")
	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		if dropLineCount > 0 {
			dropLineCount--
		} else {
			fmt.Fprintln(w, html.EscapeString(sc.Text()))
		}
	}
}

func mains(args []string) error {
	if len(args) < 1 {
		return errors.New("FILENAME is required")
	}
	handler := &Handler{
		path:     args[0],
		notFound: http.NotFoundHandler(),
	}
	service := &http.Server{
		Addr:           fmt.Sprintf(":%d", *portNo),
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
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
