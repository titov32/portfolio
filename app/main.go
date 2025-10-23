package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	data "portfolio/modules"

	"github.com/jung-kurt/gofpdf"
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
		// Загружаем YAML → Resume
		resume, _ := data.LoadData("data.yml")

		pdfBytes, err := renderPDF(resume)
		if err != nil {
			http.Error(w, "PDF generation error: "+err.Error(), http.StatusInternalServerError)
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

func renderPDF(resume *data.Resume) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, resume.Main.Name)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 14)
	pdf.Cell(0, 10, resume.Main.JobTitle)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Email: "+resume.Main.Email)
	pdf.Ln(8)
	pdf.Cell(0, 10, "GitHub: "+resume.Main.GitHub)
	pdf.Ln(8)
	pdf.Cell(0, 10, "LinkedIn: "+resume.Main.LinkedIn)
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Skills")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8,
		"Programming Languages: "+resume.Skills.ProgrammingLanguages+"\n"+
			"Frameworks: "+resume.Skills.Frameworks+"\n"+
			"Tools: "+resume.Skills.Tools,
		"", "", false)

	// пишем в память
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
