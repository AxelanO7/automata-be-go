package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go-backend/internal/llm"
	"go-backend/internal/telegram"

	"github.com/gocolly/colly"
)

func RunSniper() {
	log.Println("[SNIPER WORKER] Menginisiasi perburuan data (Hunt mode)...")
	c := colly.NewCollector()

	var pageTextBuffer []string

	c.OnHTML("body", func(e *colly.HTMLElement) {
		textChunk := strings.TrimSpace(e.Text)
		
		// Batasi buffer string untuk agar tidak membebani konteks Ollama
		if len(textChunk) > 1500 {
			textChunk = textChunk[:1500] 
		}
		if textChunk != "" {
			pageTextBuffer = append(pageTextBuffer, textChunk)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		content := strings.Join(pageTextBuffer, " \n")
		prompt := fmt.Sprintf(`Analisis teks berikut. Jika teks ini merupakan penawaran loker IT atau Programming dengan mode remote dan gaji tinggi (biasanya dalam USD atau disebutkan lumayan tinggi), kembalikan response dalam format JSON murni.
		{"is_target": true, "headline": "[buat headline FOMO yang sangat singkat dan meyakinkan]"}.
		Jika tidak relevan atau bukan loker remote berbayar tinggi, isi field "is_target" dengan false dan "headline" kosong.
		
		Teks: %s`, content)

		log.Println("[SNIPER WORKER] Menganalisis menggunakan struktur LLM...")
		jsonRawResponse, err := llm.AskOllama(prompt, true)
		if err != nil {
			log.Printf("[SNIPER ERROR] Gagal berdiskusi dengan LLM: %v\n", err)
			return
		}

		var analysis struct {
			IsTarget bool   `json:"is_target"`
			Headline string `json:"headline"`
		}

		if err := json.Unmarshal([]byte(jsonRawResponse), &analysis); err != nil {
			log.Printf("[SNIPER ERROR] Gagal parsing respons JSON LLM: %v | Raw: %s\n", err, jsonRawResponse)
			return
		}

		if analysis.IsTarget {
			log.Printf("[SNIPER BINGO] Target Dikonfirmasi! Menyebarkan alert telegram...\n")
			alertMsg := fmt.Sprintf("🚨 <b>HIGH PAYING REMOTE JOB DETECTED!</b> 🚨\n\n%s\n\n<b>Source URL:</b> %s", analysis.Headline, r.Request.URL.String())
			
			if err := telegram.SendAlert(alertMsg); err != nil {
				log.Printf("[SNIPER ERROR] Gagal mengirim alert TG: %v\n", err)
			}
		} else {
			log.Println("[SNIPER WORKER] Menarik mundur. Lokasi tidak memenuhi syarat (No target matched).")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[SNIPER HTTP ERROR] Request gagal - %v URL: %s\n", err, r.Request.URL)
	})

	c.Visit("https://weworkremotely.com/categories/remote-programming-jobs")
	log.Println("[SNIPER WORKER] Eksekusi tuntas.")
}
