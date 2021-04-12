package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	vscale "github.com/vscale/go-vscale"
)

const (
	appName = "vscale-backup"
)

var (
	flToken      = flag.String("token", os.Getenv("API_TOKEN"), "Vscale API token")
	flExpiration = flag.String("expiration", os.Getenv("BACKUP_EXPIRATION"), "Backups expiration time, older backups will be removed")
	flCron       = flag.String("cron", os.Getenv("CRON"), "Cron spec for periodically backup run")
)

var (
	version = "dev"

	c        *vscale.WebClient
	interval time.Duration
)

func main() {
	flag.Parse()

	if *flToken == "" || *flExpiration == "" || *flCron == "" {
		fmt.Println(`Usage example: vscale-backup -token "vscale-api-token" -expiration "48h" -cron "15 3 * * *"`)
		os.Exit(1)
	}

	log.Printf("starting %v version=%q", appName, version)

	var err error
	interval, err = time.ParseDuration(*flExpiration)
	if err != nil {
		log.Fatalf("Invalid backup expiration '%s': %v", *flExpiration, err)
	}

	c = vscale.NewClient(*flToken)
	cr := cron.New()

	_, err = cr.AddFunc(*flCron, processBackups)
	if err != nil {
		log.Fatalf("Error creating cron job with spec '%s': %v", *flCron, err)
	}
	cr.Start()

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}

// processBackups do create/delete jobs for vscale backups
func processBackups() {
	backupValidTill := time.Now().Add(-interval)
	servers, _, err := c.Scalet.List()
	if err != nil {
		log.Fatalf("Failed to retrieve servlets: %v", err)
	}

	backupDate := time.Now().Format("2006-01-02")

	for _, s := range *servers {
		backupName := fmt.Sprintf("%s_%s", s.Name, backupDate)
		_, _, err := c.Scalet.Backup(s.CTID, backupName)
		if err != nil {
			log.Printf("Error creating backup for servlet %s: %v", s.Name, err)
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
			log.Printf("Skipping backup %s deletion: invalid creation time: %v", b.Name, err)
			continue
		}
		if createdAt.After(backupValidTill) {
			log.Printf("Skipping backup %s deletion: not expired", b.Name)
			continue
		}
		_, _, err = c.Backup.Remove(b.ID)
		if err != nil {
			log.Printf("Error removing old backup %s: %v", b.Name, err)
			continue
		}
		log.Printf("Backup '%s' successfully deleted", b.Name)
	}
}
