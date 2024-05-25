package utils

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/palp1tate/FlowFederate/api/global"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func UploadFile(data []byte, size int64, suffix string) (url string, err error) {
	key := fmt.Sprintf("%s%s", GenerateUUID(), suffix)
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", global.ServerConfig.QiNiuYun.Bucket, key),
	}
	mac := qbox.NewMac(global.ServerConfig.QiNiuYun.AccessKey, global.ServerConfig.QiNiuYun.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	uploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err = uploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), size, &putExtra)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s/%s", global.ServerConfig.QiNiuYun.Domain, ret.Key)
	return
}

func DeleteFile(url string) (err error) {
	mac := qbox.NewMac(global.ServerConfig.QiNiuYun.AccessKey, global.ServerConfig.QiNiuYun.SecretKey)
	bucketManager := storage.NewBucketManager(mac, nil)
	key := filepath.Base(url)
	err = bucketManager.Delete(global.ServerConfig.QiNiuYun.Bucket, key)
	return
}
