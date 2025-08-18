package backup

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MySqlBackup struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func (m MySqlBackup) Dump(ctx context.Context, opt ...DumpDbOpts) (string, error) {
	// final tar.gz file
	filename := fmt.Sprintf("%s_%s.tar.gz", m.Database, time.Now().Format("20060102_150405"))
	outputPath := filepath.Join(".", filename)

	// create output tar.gz file
	outfile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	// gzip writer
	gzw := gzip.NewWriter(outfile)
	defer gzw.Close()

	// tar writer
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// create a tar header for the .sql file inside archive
	sqlFileName := fmt.Sprintf("%s.sql", m.Database)
	header := &tar.Header{
		Name:    sqlFileName,
		Mode:    0600,
		Size:    0, // we'll use a pipe so size is unknown
		ModTime: time.Now(),
	}
	// we will fill Size later via pipe, so skip direct WriteHeader+Size.

	// build mysqldump args
	args := []string{
		"-h", m.Host,
		"-P", m.Port,
		"-u", m.User,
		fmt.Sprintf("--password=%s", m.Password),
		m.Database,
	}

	// prepare command
	cmd := exec.CommandContext(ctx, "mysqldump", args...)

	// pipe mysqldump stdout directly into tar entry
	pr, pw, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer pr.Close()

	cmd.Stdout = pw
	cmd.Stderr = os.Stderr

	// start mysqldump
	if err := cmd.Start(); err != nil {
		pw.Close()
		return "", fmt.Errorf("mysqldump start failed: %w", err)
	}

	// close write end after command finishes
	go func() {
		cmd.Wait()
		pw.Close()
	}()

	// since tar requires knowing size, we stream by reading
	// workaround: buffer mysqldump into tar without setting Size
	// (using io.Copy + WriteHeader with zero size + special flag)
	// Instead we can capture to temp buffer, but that uses memory.
	// Simpler: copy to tar directly, but we need Size…
	// Solution: use io.Pipe again to avoid full buffer.

	tmpFile, err := os.CreateTemp("", "mysqldump-*.sql")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// copy mysqldump output to tmp file
	if _, err := pr.WriteTo(tmpFile); err != nil {
		return "", fmt.Errorf("write to temp failed: %w", err)
	}

	// rewind temp file
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return "", err
	}

	// get size
	info, _ := tmpFile.Stat()
	header.Size = info.Size()

	// write tar header
	if err := tw.WriteHeader(header); err != nil {
		return "", err
	}

	// copy file into tar
	if _, err := tmpFile.WriteTo(tw); err != nil {
		return "", err
	}

	return outputPath, nil
}
