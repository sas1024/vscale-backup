# VScale backup service

This backup service creates new and automatically removes old backups for [Vscale](https://vscale.io) instances.

It uses Vscale API for all operations and requires API token with write access.

## Installation
To install vscale backup service use go get:
```bash
go get -u github.com/sas1024/vscale-backup
```
## Running via CLI

```bash
vscale-backup -token "vscale-api-token" -expiration "48h" -cron "3 15 * * *"
```

## Running via Docker
```bash
docker run -it -e "API_TOKEN=<vscale-api-token>" -e "BACKUP_EXPIRATION=48h" -e 'CRON="3 15 * * *"' sas1024/vscale-backup
```
