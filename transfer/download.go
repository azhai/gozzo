package transfer

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Code-Hex/pget"
	"github.com/azhai/gozzo/filesystem"
)

type Downloader struct {
	procs   int
	dirname string
	*http.Client
}

func NewDownloader(dirname string, procs int) *Downloader {
	if procs <= 0 {
		procs = runtime.NumCPU()
	}
	return &Downloader{
		dirname: dirname, procs: procs,
		Client: http.DefaultClient,
	}
}

func (d *Downloader) GetConfig(fileUrl, savePath string) *pget.DownloadConfig {
	return &pget.DownloadConfig{
		Client:   d.Client,
		Procs:    d.procs,
		Dirname:  filepath.Dir(savePath),
		Filename: filepath.Base(savePath),
		URLs:     []string{fileUrl},
	}
}

func (d *Downloader) GetCheckConfig(fileUrl string) *pget.CheckConfig {
	return &pget.CheckConfig{
		Client:  d.Client,
		Timeout: 10 * time.Second,
		URLs:    []string{fileUrl},
	}
}

func (d *Downloader) GetLength(fileUrl string) (int64, error) {
	var size int64
	resp, err := d.Client.Head(fileUrl)
	if err == nil && resp != nil {
		size = resp.ContentLength
	}
	return size, err
}

func (d *Downloader) GetLengthContext(ctx context.Context, fileUrl string) (int64, error) {
	var size int64
	target, err := pget.Check(ctx, d.GetCheckConfig(fileUrl))
	if err == nil && target != nil {
		size = target.ContentLength
	}
	return size, err
}

func (d *Downloader) GetSavePath(fileUrl, fileName string) string {
	if fileName == "" {
		if u, err := url.Parse(fileUrl); err == nil {
			fileName = filepath.Base(u.Path)
		} else {
			fileName = filepath.Base(fileUrl)
		}
	}
	return filepath.Join(d.dirname, fileName)
}

// Download 下载文件到指定位置
func (d *Downloader) Download(fileUrl, fileName string, force bool) (size int64, err error) {
	savePath := d.GetSavePath(fileUrl, fileName)
	if !force && filesystem.NewFileHandler(savePath).IsExist() {
		return
	}
	if size, err = d.GetLength(fileUrl); err != nil {
		return
	}
	conf := d.GetConfig(fileUrl, savePath)
	conf.ContentLength = size
	ctx := context.Background()
	err = pget.Download(ctx, conf)
	return
}

// DownloadIgnore 下载文件到指定位置，同名文件存在时会被覆盖
func (d *Downloader) DownloadIgnore(fileUrl, fileName string) (int64, error) {
	return d.Download(fileUrl, fileName, true)
}

// DownloadIfNot 下载文件到指定位置，如果不存在
func (d *Downloader) DownloadIfNot(fileUrl, fileName string) (int64, error) {
	return d.Download(fileUrl, fileName, false)
}
