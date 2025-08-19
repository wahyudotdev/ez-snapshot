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

### Build from source

```shell
git clone https://github.com/wahyudotdev/ez-snapshot.git
cd ez-snapshot

# build binary
go build -o ez-snapshot cmd/main.go

# move to bin folder so it can be executed anywhere
sudo mv ez-snapshot /usr/local/bin/ez-snapshot
sudo chmod +x /usr/local/bin/ez-snapshot

# copy the config files
mkdir ~/.config/ez-snapshot
cp config.example.yaml ~/.config/ez-snapshot/config.yaml

# run the binary
ez-snapshot --help
```

## How to use

Make sure that you have followed the installation instruction above

1. Copy the config.example.yaml to config.yaml and adjust the mysql & rclone configuration. You can also place the
   configuration under ```$HOME/.config/ez-snapshot/config.yaml```.
2. You need to configure the RClone first, go to the RClone documentation for complete guideline for each online
   storage.
3. Run RClone RC API server
4. Run the binary by using ```ez-snapshot``` command, this will check for dependency and API connectivity to RClone.

## Configuration File

You should place it under ~/.config/ez-snapshot/config.yaml so it can be automatically read by ez-snapshot process

| Key              | Explanation                                                    |
|------------------|----------------------------------------------------------------|
| `mysql.host`     | `127.0.0.1` (MySQL Host)                                       |
| `mysql.port`     | `3306` (MySQL Port)                                            |
| `mysql.username` | `root` (MySQL username)                                        |
| `mysql.password` | `password` (MySQL password)                                    |
| `mysql.database` | `db` (MySQL DB schema         )                                |
| `rclone.host`    | `http://localhost:5572` (rclone API host, no auth)             |
| `rclone.fs`      | `s3:mybucket` → `s3` = rclone remote, `mybucket` = bucket name |
| `rclone.remote`  | `db-backup` remote path. Backup files would be stored here     |

## Non-Interactive CLI

You could also use non-interactive CLI to execute ```backup``` and ```restore``` command directly, that will be
useful for cron-job setup. Example crontab configuration

```text
# run backup every 00.00 AM
0 0 * * * /usr/local/bin/ez-snapshot --backup
```

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
- ⌛ Provide file encryption support