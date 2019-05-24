package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/russross/blackfriday/v2"
)

func main() {
	idx := NewIndex("src")
	idx.Write("dst")
}

type Post struct {
	Title       string
	URL         string
	Description string
	Date        string
	Content     string
}

type Index struct {
	Posts     []Post
	wg        sync.WaitGroup
	templates *template.Template

	// dummy for template
	Title       string
	Description string
	URL         string
}

func NewIndex(dir string) *Index {

	w := NewWalker()
	i := &Index{
		Description: "blog of seankhliao",
	}

	t := template.Must(template.New("head").Parse(HeadTemplate))
	t = template.Must(t.New("foot").Parse(FootTemplate))
	t = template.Must(t.New("post").Parse(PostTemplate))
	t = template.Must(t.New("index").Parse(IndexTemplate))

	i.templates = t

	go w.walk(dir)
	for p := range w.p {
		i.Posts = append(i.Posts, p)
	}
	sort.Sort(i)
	return i
}

func (i *Index) Write(dir string) error {
	if err := os.MkdirAll("dst", 0744); err != nil {
		log.Printf("error creating dir %v: %v\n", "dst", err)
		return err
	}

	i.wg.Add(1)
	go i.writeFile("index", path.Join(dir, "index.html"), i)

	for _, p := range i.Posts {
		i.wg.Add(1)
		go i.writeFile("post", path.Join(dir, p.URL+".html"), p)
	}

	i.wg.Wait()
	return nil
}

func (i *Index) writeFile(tmp, fn string, data interface{}) error {
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("error creating file %v: %v\n", fn, err)
		return err
	}
	defer f.Close()
	if err := i.templates.ExecuteTemplate(f, tmp, data); err != nil {
		log.Printf("error executing template %v for %v: %v\n", tmp, fn, err)
		return err
	}
	i.wg.Done()
	return nil
}

// newer first
func (idx *Index) Less(i, j int) bool { return idx.Posts[i].Date > idx.Posts[j].Date }
func (idx *Index) Len() int           { return len(idx.Posts) }
func (idx *Index) Swap(i, j int)      { idx.Posts[i], idx.Posts[j] = idx.Posts[j], idx.Posts[i] }

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
