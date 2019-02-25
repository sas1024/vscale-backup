# VScale backup tool

This backup tool is used to helps create and automaticaly remove old backups for [Vscale](https://vscale.io) instances.

It uses Vscale API for all operations and requires API token with write access.

## Installation
To install vscale backup tool use go get:
```go get -u github.com/sas1024/vscale-backup```

## Usage
1. Rename config.cfg.dist to config.cfg
2. Modify API token and old backups deletion interval in config.cfg
3. Run!

## Command-Line Options
```
Usage of ./vscale-backup:
  -config string
    	Path to config file (default "config.cfg")
  -verbose
    	Enable verbose output
```