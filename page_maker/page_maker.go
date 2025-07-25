package page_maker

import (
	"html/template"
	"io"
	"log"
)

type PageMaker struct {
	writer io.Writer
}

type Page struct {
	link string
	mime string
}

func (r PageMaker) Do(link string, mime string) string {
	if tmpl, err := template.New("video").ParseFiles("template/image.html.template"); err == nil {
		vars := map[string]string{"url": link}
		if err := tmpl.Execute(r.writer, vars); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Нет шаблона для страницы")
	}
	return ""
}
