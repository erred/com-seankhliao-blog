package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func main() {
	idx := NewIndex("src")
	idx.Write("dst")

	d, err := os.Create("dst/base.css")
	if err != nil {
		log.Fatalf("error creating dst/base.css: %v\n", err)
	}
	defer d.Close()
	s, err := os.Open("base.css")
	if err != nil {
		log.Fatalf("error opening base.css: %v\n", err)
	}
	defer s.Close()
	io.Copy(d, s)
}

type Walker struct {
	p      chan Post
	prefix string
}

func NewWalker() Walker {
	return Walker{
		p: make(chan Post, 5),
	}
}

func (w Walker) walk(dir string) {
	w.prefix = dir
	filepath.Walk(dir, w.walker)
	close(w.p)
}

func (w Walker) walker(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("error reading file %v: %v\n", path, err)
		return err
	}
	w.p <- parsePost(strings.TrimPrefix(path, w.prefix+"/"), string(b))
	return nil
}

func parsePost(filename, text string) Post {
	t := strings.SplitN(text, "\n---\n", 2)
	if len(t) < 2 {
		log.Printf("parseText --- not found in %v\n", filename)
	}

	p := Post{
		URL:     strings.ReplaceAll(strings.ToLower(strings.TrimSuffix(filename, ".md")), " ", "-"),
		Content: string(blackfriday.Run([]byte(t[1]))),
	}

	header := strings.Split(strings.TrimSpace(t[0]), "\n")
	if len(header) < 3 {
		log.Printf("parseText in %v epected 3 headers, found %v\n", filename, len(header))
	}
	for _, head := range header {
		h := strings.SplitN(head, "=", 2)
		if len(h) < 2 {
			log.Printf("parseText in %v in headers expected split by = in >%v<\n", filename, head)
		}
		v := strings.TrimSpace(h[1])
		switch strings.TrimSpace(h[0]) {
		case "title":
			p.Title = v
		case "date":
			p.Date = v
		case "description", "desc":
			p.Description = v
		}
	}
	return p
}
