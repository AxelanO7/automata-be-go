package workers

import (
	"encoding/json"
	"fmt"
	"log"

	"go-backend/internal/db"
	"go-backend/internal/llm"
)

type PseoArticleRecord struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	ContentHTML string `json:"content_html"`
	MetaDesc    string `json:"meta_desc"`
}

func RunPSEO() {
	log.Println("[PSEO WORKER] Memulai operasi massal Programmatic SEO Generation...")
	topics := []string{"Jasa Web Dev Bali", "SEO Agency Jakarta", "Software House Surabaya"}

	for _, topic := range topics {
		log.Printf("[PSEO WORKER] Sedang memproses Keyword: %s", topic)
		
		prompt := fmt.Sprintf(`Buatkan satu artikel Programmatic SEO (PSEO) berkualitas dengan struktur silokan standar untuk Topik: "%s".
		Wajib kembalikan format JSON murni TANPA TEKS LAIN sama sekali untuk parameter ini:
		{
			"slug": "url-slug-kebab-case",
			"title": "Judul Menarik (H1 Setara)",
			"content_html": "<article><h2>...</h2><p>...</p><ul><li>...</li></ul></article>",
			"meta_desc": "Deskripsi meta super singkat untuk standar SEO Google (Max 160 char)"
		}`, topic)

		jsonRaw, err := llm.AskOllama(prompt, true)
		if err != nil {
			log.Printf("[PSEO ERROR] LLM Gagal merespon pada topik %s: %v\n", topic, err)
			continue
		}

		var article PseoArticleRecord
		if err := json.Unmarshal([]byte(jsonRaw), &article); err != nil {
			log.Printf("[PSEO ERROR] Unmarshall error untuk JSON pada keyword %s: %v\nRaw Content: %s\n", topic, err, jsonRaw)
			continue
		}

		// Inject to Database Supabase
		insertResp := db.Client.From("pseo_articles").Insert(article, false, "", "", "exact")
		_, _, dbErr := insertResp.Execute()
		
		if dbErr != nil {
			log.Printf("[PSEO ERROR] Gagal push record ke layer Database (%s): %v\n", article.Slug, dbErr)
			continue
		}

		log.Printf("[PSEO SUCCESS] Artikel Ditanam secara sukses ke database: /%s\n", article.Slug)
	}
	log.Println("[PSEO WORKER] Batch run selesai.")
}
