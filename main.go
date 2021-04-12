package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	vscale "github.com/vscale/go-vscale"
)

const (
	appName = "vscale-backup"
)

var (
	flToken      = flag.String("token", os.Getenv("API_TOKEN"), "Vscale API token")
	flExpiration = flag.String("expiration", os.Getenv("BACKUP_EXPIRATION"), "Backups expiration time, older backups will be removed")
)

var version = "dev"

func main() {
	flag.Parse()

	if *flToken == "" || *flExpiration == "" {
		fmt.Println(`Usage: vscale-backup -token "vscale-api-token" -expiration "30d"`)
		os.Exit(1)
	}

	log.Printf("starting %v version=%q", appName, version)

	interval, err := time.ParseDuration(*flExpiration)
	if err != nil {
		log.Fatalf("Invalid backup expiration '%s': %q", *flExpiration, err)
	}
	backupValidTill := time.Now().Add(-interval)

	c := vscale.NewClient(*flToken)
	servers, _, err := c.Scalet.List()
	if err != nil {
		log.Fatalf("Failed to retrieve servlets: %q", err)
	}

	backupDate := time.Now().Format("2006-01-02")

	for _, s := range *servers {
		backupName := fmt.Sprintf("%s_%s", s.Name, backupDate)
		_, _, err := c.Scalet.Backup(s.CTID, backupName)
		if err != nil {
			log.Printf("Error creating backup for servlet %s: %q", s.Name, err)
			continue
		}
		log.Printf("Backup %s for servlet %s successfully created", backupName, s.Name)
	}
	backups, _, err := c.Backup.List()
	if err != nil {
		log.Fatalf("Failed to fetch backup list: %v", err)
	}
	for _, b := range *backups {
		if b.Status != "finished" {
			log.Printf("Skipping backup %s deletion: unfinished status", b.Name)
			continue
		}
		createdAt, err := time.ParseInLocation("02.01.2006 15:04:05", b.Created, time.Local)
		if err != nil {
			log.Printf("Skipping backup %s deletion: invalid creation time: %q", b.Name, err)
			continue
		}
		if createdAt.After(backupValidTill) {
			log.Printf("Skipping backup %s deletion: not expired", b.Name)
			continue
		}
		_, _, err = c.Backup.Remove(b.ID)
		if err != nil {
			log.Printf("Error removing old backup %s: %q", b.Name, err)
			continue
		}
		log.Printf("Backup '%s' successfully deleted", b.Name)
	}
}
