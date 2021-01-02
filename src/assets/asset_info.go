package assets

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func (assets *Assets) Get(assetName string) (*Info, error) {
	info := assets.cache[assetName]
	if info != nil {
		return info, nil
	}

	fileName := filepath.Join(assets.assetsPath, assetName[assets.assetsUriLen:])

	contentType := mime.TypeByExtension(path.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	info, err := getInfo(fileName + ".gz", contentType, GZip)
	if err != nil {
		return nil, err
	}
	if info != nil {
		assets.cache[assetName] = info
		return info, nil
	}

	info, err = getInfo(fileName, contentType, Identity)
	if err != nil {
		return nil, err
	}
	if info != nil {
		assets.cache[assetName] = info
		return info, nil
	}

	return nil, nil
}

func getInfo(fileName, contentType, contentEncoding string) (*Info, error) {
	stat, err := getStat(fileName)
	if err != nil {
		return nil, err
	}
	if stat == nil {
		return nil, nil
	}
	etag, err := getETag(fileName)
	if err != nil {
		return nil, err
	}
	info := &Info{
		FilePath: fileName,
		ContentType: contentType,
		ContentLength: strconv.FormatInt(stat.Size(), 10),
		ContentEncoding: contentEncoding,
		ETag: etag,
	}
	return info, nil
}

func getStat(fileName string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(fileName)
	switch {
	case fileInfo != nil && fileInfo.IsDir():
		return nil, nil
	case os.IsNotExist(err):
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return fileInfo, err
	}
}

func getETag(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)[:12]), nil
}
