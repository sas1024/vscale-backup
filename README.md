# VScale backup tool

This backup tool creates new and automatically removes old backups for [Vscale](https://vscale.io) instances.

It uses Vscale API for all operations and requires API token with write access.

## Installation
To install vscale backup tool use go get:
```bash
go get -u github.com/sas1024/vscale-backup
```
## Running via CLI

```bash
vscale-backup -token "vscale-api-token" -expiration "48h"
```

## Running via Docker
```bash
docker run -it -e "API_TOKEN=<vscale-api-token>" -e "BACKUP_EXPIRATION=48h" sas1024/vscale-backup
```
