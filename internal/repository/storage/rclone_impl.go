package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"ez-snapshot/internal/entity"
)

type rCloneImpl struct {
	host   string
	fs     string
	remote string
	client *http.Client
}

func newRCloneImpl(host, fs, remote string) Repository {
	return &rCloneImpl{
		host:   host, // e.g. "http://localhost:5572"
		fs:     fs,   // e.g. "s3remote:mybucket"
		remote: remote,
		client: &http.Client{},
	}
}

func (rc *rCloneImpl) Upload(ctx context.Context, key string, r io.Reader) (string, error) {
	endpoint := rc.host + "/operations/uploadfile"
	values := url.Values{}
	values.Set("fs", rc.fs)
	values.Set("remote", rc.remote) // directory (like your working curl: remote=db-backup)

	// Pipe so we can stream multipart directly to the request (no full buffering)
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Build request that reads from the pipe
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"?"+values.Encode(), pr)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Printf("Begin upload to %s/%s\n", rc.remote, key)

	// Write multipart in a goroutine so Do(req) can read concurrently
	go func() {
		var err error
		defer func() {
			// Close the multipart writer first to flush the boundary,
			// then close the pipe (propagate any error up the read side).
			if cerr := writer.Close(); err == nil && cerr != nil {
				err = cerr
			}
			if err != nil {
				_ = pw.CloseWithError(err)
			} else {
				_ = pw.Close()
			}
		}()

		// Create the single unnamed form file field (like curl: --form '=@file')
		part, e := writer.CreateFormFile("", key)
		if e != nil {
			err = e
			return
		}

		// Progress counter prints periodically while data flows
		pc := &progressCounter{}
		tee := io.TeeReader(r, pc)

		_, e = io.Copy(part, tee)
		if e != nil {
			err = e
			return
		}
		// final newline after the \r-updating line
		fmt.Printf("\rUploaded %d bytes\n", pc.n)
	}()

	// Send request (streams from pr)
	resp, err := rc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(b))
	}

	fmt.Println("✅ Backup has been uploaded")
	return key, nil
}

// progressCounter prints a running byte count without spamming the console.
type progressCounter struct {
	n    int64
	last time.Time
}

func (pc *progressCounter) Write(p []byte) (int, error) {
	pc.n += int64(len(p))
	now := time.Now()
	if now.Sub(pc.last) >= 500*time.Millisecond {
		fmt.Printf("\rUploaded %d bytes", pc.n)
		pc.last = now
	}
	return len(p), nil
}

func (rc *rCloneImpl) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	endpoint := rc.host + "/operations/publiclink"
	values := url.Values{}
	values.Set("fs", rc.fs)
	values.Set("remote", key)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("publiclink failed: %s", string(b))
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.URL == "" {
		return nil, fmt.Errorf("publiclink returned empty url")
	}

	// Now fetch the actual file using the signed URL
	fileReq, err := http.NewRequestWithContext(ctx, "GET", result.URL, nil)
	if err != nil {
		return nil, err
	}

	fileResp, err := rc.client.Do(fileReq)
	if err != nil {
		return nil, err
	}

	if fileResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(fileResp.Body)
		fileResp.Body.Close()
		return nil, fmt.Errorf("download failed: %s", string(b))
	}

	return fileResp.Body, nil // caller must Close()
}

func (rc *rCloneImpl) Delete(ctx context.Context, key string) error {
	endpoint := rc.host + "/operations/deletefile"
	values := url.Values{}
	values.Set("fs", rc.fs)
	values.Set("remote", key)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"?"+values.Encode(), nil)
	if err != nil {
		return err
	}

	resp, err := rc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(b))
	}

	return nil
}

func (rc *rCloneImpl) List(ctx context.Context) ([]*entity.Backup, error) {
	endpoint := rc.host + "/operations/list"
	values := url.Values{}
	values.Set("fs", rc.fs)
	values.Set("remote", rc.remote) // root path, or change to subdir if needed

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed: %s", string(b))
	}

	var result struct {
		List []*entity.Backup `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	backups := make([]*entity.Backup, 0, len(result.List))
	for _, f := range result.List {
		if f.IsDir {
			continue
		}
		backups = append(backups, f)
	}

	return backups, nil
}
