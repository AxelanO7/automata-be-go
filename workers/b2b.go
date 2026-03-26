package workers

import (
	"fmt"
	"log"

	"go-backend/internal/db"
	"go-backend/internal/llm"

	"github.com/gocolly/colly"
)

type B2BLeadRecord struct {
	BusinessName  string `json:"business_name"`
	WebsiteStatus string `json:"website_status"`
	Niche         string `json:"niche"`
	AIEmailDraft  string `json:"ai_email_draft"`
}

func RunB2B() {
	log.Println("[B2B WORKER] Hunting Directory Bisnis untuk B2B Lead Gen...")
	c := colly.NewCollector()

	nicheProfile := "Klinik Kecantikan / Aesthetic Clinic"

	// Selector fleksibel, sesuaikan jika Anda scrape dari platform khusus (contoh: Yellowpages/Yelp)
	c.OnHTML(".business-card-element", func(e *colly.HTMLElement) {
		bizName := e.ChildText(".biz-title")
		websiteUrl := e.ChildAttr("a.website", "href")

		if bizName == "" {
			return
		}

		if websiteUrl == "" || websiteUrl == "#" {
			log.Printf("[B2B ENGINE] Menemukan target rentan secara digital (TIDAK ADA WEBSITE): %s", bizName)
			
			prompt := fmt.Sprintf(`Tuliskan draf cold email yang sopan, bersahabat, profesional, dan persuasif namun tanpa basa-basi berbunga-bunga, untuk menawarkan jasa pembuatan website profesional kepada bisnis "%s" dengan ceruk bisnis "%s". 
			Jangan sertakan placeholder semacam [Nama Anda], buat email utuh siap dikirim secara anonim atau dari 'Digital Architect Team'. Taruh di format teks biasa.`, bizName, nicheProfile)

			emailTextHTML, err := llm.AskOllama(prompt, false)
			if err != nil {
				log.Printf("[B2B ERROR] LLM menolak kalkulasi draft untuk %s: %v\n", bizName, err)
				return
			}

			payload := B2BLeadRecord{
				BusinessName:  bizName,
				WebsiteStatus: "NO_WEBSITE",
				Niche:         nicheProfile,
				AIEmailDraft:  emailTextHTML,
			}

			dbInsert := db.Client.From("b2b_leads").Insert(payload, false, "", "", "exact")
			if _, _, err := dbInsert.Execute(); err != nil {
				log.Printf("[B2B ERROR] Supabase menolak Lead Entry (%s): %v\n", bizName, err)
				return
			}
			
			log.Printf("[B2B HARVEST] Draft berhasil dikunci & dibakar ke database Lead: %s\n", bizName)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[B2B HTTP ERROR] Deteksi pemblokiran %v\n", err)
	})

	c.Visit("https://example-directory-site.com/indonesia/bali/klinik-kecantikan")
	log.Println("[B2B WORKER] Eksplorasi direktori selesai. Workers Terminated.")
}
