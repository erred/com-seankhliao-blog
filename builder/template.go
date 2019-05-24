package main

import (
	"log"
	"os"
	"path"
	"sort"
	"sync"
	"text/template"
)

type Post struct {
	Title       string
	URL         string
	Description string
	Date        string
	Content     string
}

const PostTemplate = `
<!doctype html>
<html lang="en">
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
  <meta name="theme-color" content="#000000">
  <script async src="https://www.googletagmanager.com/gtag/js?id=UA-114337586-4"></script>
  <script>window.dataLayer = window.dataLayer || [];function gtag() {dataLayer.push(arguments);};gtag("js", new Date());gtag("config", "UA-114337586-4");</script>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Inconsolata:400,700&display=swap">
  <link rel="stylesheet" href="./base.css">

  <title>{{ .Title }} | blog | seankhliao</title>

  <link rel="canonical" href="https://blog.seankhliao.com/{{ .URL }}">
  <meta name="description" content="{{ .Description }}">

  <hgroup>
    <h1>{{ .Title }}</h1>
    <p>
      <a href="https://seankhliao.com">seankhliao</a> / <a href="https://blog.seankhliao.com">blog</a> /
      <span>{{ .Date }}</span>
    </p>
  </hgroup>

  {{ .Content }}

</html>
`

type Index struct {
	Posts     []Post
	wg        sync.WaitGroup
	templates *template.Template
}

func NewIndex(dir string) *Index {

	w := NewWalker()
	i := &Index{}

	i.templates = template.Must(template.Must(template.New("post").Parse(PostTemplate)).New("index").Parse(IndexTemplate))

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

const IndexTemplate = `
<!doctype html>
<html lang="en">
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
  <meta name="theme-color" content="#000000">
  <script async src="https://www.googletagmanager.com/gtag/js?id=UA-114337586-4"></script>
  <script>window.dataLayer = window.dataLayer || [];function gtag() {dataLayer.push(arguments);};gtag("js", new Date());gtag("config", "UA-114337586-4");</script>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Inconsolata:400,700&display=swap">
  <link rel="stylesheet" href="./base.css">
  
  <title>blog | seankhliao</title>
  
  <link rel="canonical" href="https://blog.seankhliao.com" />
  <meta name="description" content="blog of seankhliao, it probably makes sense to somebody" />
  
  <h1><a href="https://seankhliao.com">seankhliao</a> / blog</h1>
  <p>Artisanal, hand-crafted blog posts imbued with delayed regrets</p>
  
  <ul>
    {{ range .Posts }}
        <li><span>{{ .Date }}</span> | <a href="./{{ .URL }}">{{ .Title }}</a></li>
    {{ end }}
  </ul>
</html>
`
