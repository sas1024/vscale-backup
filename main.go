package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	vscale "github.com/vscale/go-vscale"
)

const (
	appName = "vscale-backup"
)

type Config struct {
	Token            string
	DeletionInterval string
}

var (
	flVerbose    = flag.Bool("verbose", false, "Enable verbose output")
	flConfigPath = flag.String("config", "config.cfg", "Path to config file")
	cfg          Config
)

var version = "dev"

func main() {
	flag.Parse()
	fixStdLog(*flVerbose)
	log.Printf("starting %v version=%v", appName, version)

	if _, err := toml.DecodeFile(*flConfigPath, &cfg); err != nil {
		die(err)
	}

	interval, err := time.ParseDuration(cfg.DeletionInterval)
	if err != nil {
		die(err)
	}
	backupValidTill := time.Now().Add(-interval)

	c := vscale.NewClient(cfg.Token)
	servers, _, err := c.Scalet.List()
	if err != nil {
		die(err)
	}

	backupDate := time.Now().Format("2006.01.02")

	for _, s := range *servers {
		backupName := fmt.Sprintf("%s_%s", s.Name, backupDate)
		_, _, err := c.Scalet.Backup(s.CTID, backupName)
		if err != nil {
			log.Printf("Error creating backup for %s: %v", s.Name, err)
			continue
		}
		log.Printf("Create new backup %s for %s", backupName, s.Name)
	}
	backups, _, err := c.Backup.List()
	if err != nil {
		die(err)
	}
	for _, b := range *backups {
		if b.Status != "finished" {
			log.Printf("Skip old backup %s deletion: unfinished status", b.Name)
			continue
		}
		createdAt, err := time.ParseInLocation("02.01.2006 15:04:05", b.Created, time.Local)
		if err != nil {
			log.Printf("Skip old backup %s deletion: error parsing creation time: %v", b.Name, err)
			continue
		}
		if createdAt.After(backupValidTill) {
			log.Printf("Skip old backup %s deletion: too young", b.Name)
			continue
		}
		_, _, err = c.Backup.Remove(b.ID)
		if err != nil {
			log.Printf("Error removing old backup %s: %v", b.Name, err)
			continue
		}
		log.Println("Delete old backup", b.Name)
	}
}

// die calls log.Fatal if err wasn't nil.
func die(err error) {
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}

// fixStdLog sets additional params to std logger.
func fixStdLog(verbose bool) {
	log.SetPrefix("D")
	log.SetFlags(log.LstdFlags)

	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}
