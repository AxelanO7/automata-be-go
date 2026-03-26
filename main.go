package main

import (
	"flag"
	"log"

	"go-backend/internal/db"
	"go-backend/workers"

	"github.com/joho/godotenv"
)

func main() {
	runFlag := flag.String("run", "", "Tentukan worker yang akan dieksekusi: sniper, pseo, b2b")
	flag.Parse()

	// Load file .env dari level root/parent apabila dijalankan di nested directory
	_ = godotenv.Load("../.env")

	if err := db.InitSupabase(); err != nil {
		log.Fatalf("Kesalahan Fatal: Gagal inisialisasi Supabase - %v\n", err)
	}

	switch *runFlag {
	case "sniper":
		workers.RunSniper()
	case "pseo":
		workers.RunPSEO()
	case "b2b":
		workers.RunB2B()
	case "":
		log.Println("Silakan berikan flag --run. Contoh eksekusi: go run main.go --run=sniper")
	default:
		log.Printf("Kesalahan: Worker '%s' tidak dikenali dalam sistem. Opsi: sniper, pseo, b2b\n", *runFlag)
	}
}
