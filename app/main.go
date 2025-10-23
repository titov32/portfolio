package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	data "portfolio/modules"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {

	cfg, err := data.Load("data.yml")
	if err != nil {
		log.Fatal("Ошибка загрузки конфига:", err)
	}

	fmt.Println("Host:", cfg.Server.Host)
	fmt.Println("Port:", cfg.Server.Port)
	fmt.Println("DB User:", cfg.Database.User)

	// // Отдаём статический HTML (твоя страница)
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))

	// Раздача статики из /static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Тут ты можешь грузить данные из YAML через свою LoadData
		resume, err := data.LoadData("data.yml")
		if err != nil {
			http.Error(w, "failed to load resume", http.StatusInternalServerError)
			return
		}

		// Отдаем HTML с подставленными данными
		if err := tmpl.Execute(w, resume); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	})

	// PDF endpoint
	http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		pdfBytes, err := renderPDF("http://localhost:8080")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
		w.Write(pdfBytes)
	})

	// API для отдачи данных (можно заполнить resume.json)
	http.HandleFunc("/api/resume", func(w http.ResponseWriter, r *http.Request) {
		file, err := data.LoadData("data.yml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(file)
	})

	// Health-check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderPDF(url string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var pdfBuf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// берём три возвращаемых значения
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuf, nil
}
