# EZ-Snapshot

A simple Go-based CLI tool that wraps [`mysqldump`](https://dev.mysql.com/doc/refman/8.0/en/mysqldump.html) and [
`rclone`](https://rclone.org/) to provide easy database backup and restore operations.  
It is designed to help you back up MySQL databases to remote storage (e.g., S3, GCS, local FS) using rclone.

> **Warning**  
> This tools is not production ready since it is using mysqldump under the hood 
> that are known to be slow and intended only for development purpose.


## Features

- **List Backups** – view available backups stored in the configured remote.
- **Manual Backup** – trigger an immediate MySQL dump and upload it to remote storage.
- **Manual Restore** – download a backup from remote storage and restore it to MySQL.

## Requirements

- [Go 1.24+](https://go.dev/doc/install) (for building from source)
- [mysql](https://dev.mysql.com) and [mysqldump](https://dev.mysql.com/doc/refman/8.0/en/mysqldump.html) available in
  `$PATH`. For mac user you can install ```mysql-client``` by
  using [brew](https://formulae.brew.sh/formula/mysql-client)
- [rclone](https://rclone.org/) with [rc (remote control) API](https://rclone.org/rc/) enabled, for example:

  ```bash
  rclone rcd --rc-no-auth --rc-addr=:5572

## Installation

Clone and build:

```shell
git clone https://github.com/wahyudotdev/ez-snapshot.git
cd ez-snapshot
go build cmd/main.go
chmod +x main
./main
```

## How to use

1. Copy the config.example.yaml to config.yaml and adjust the mysql & rcloone configuration
2. You need to configure the RClone first, go to the RClone documentation for complete guideline for each online storage.
3. Run the binary, this will check for dependency and API connectivity to RClone. So make sure that you have run RClone
   RC API before.

## Project Roadmap

- ✅ Interactive CLI
- ✅ Backup compression using .tar.gz
- ✅ Support db restore for both .tar.gz and .sql file
- ✅ Support multiple storage using rclone
- ✅ Support MySQL backup and restore
- ✅️ Support non-interactive CLI
- ⌛️ Support scheduled backup (daemon mode)
- ⌛️ Support PostgresQL backup and restore
- ⌛️ Single binary release (homebrew / snap)
- ⌛️ Support RClone basic auth
