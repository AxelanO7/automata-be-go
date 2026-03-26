package db

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func InitSupabase() error {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if url == "" || key == "" {
		return fmt.Errorf("variabel SUPABASE_URL atau SUPABASE_SERVICE_ROLE_KEY belum diset")
	}

	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return err
	}

	Client = client
	return nil
}
