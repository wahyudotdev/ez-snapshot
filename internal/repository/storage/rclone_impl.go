package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	endpoint := rc.host + "/operations/uploadfile"
	values := url.Values{}
	values.Set("fs", rc.fs)
	values.Set("remote", key)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint+"?"+values.Encode(), bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := rc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(b))
	}

	return key, nil
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
