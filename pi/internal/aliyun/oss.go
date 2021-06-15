package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 获取存储空间
func (os *ClientOSS) getBucket(bucketName string) *oss.Bucket {
	bucket, err := os.Bucket(bucketName)
	if err != nil {
		logIO.Error(err)
		return nil
	}
	return bucket
}

// 上传本地文件。
func (os *ClientOSS) Upload(objectName, fileName, bucketName string) {
	err := os.getBucket(bucketName).PutObjectFromFile(objectName, fileName)
	if err != nil {
		logIO.Error(err)
		return
	}

}

// 签名直传，下载到流。
func (os *ClientOSS) GetURL(objectName, bucketName string) string {
	signedURL, err := os.getBucket(bucketName).SignURL(objectName, oss.HTTPGet, 60)
	if err != nil {
		logIO.Error(err)
	}
	return signedURL
}

//空间是否存在
func (os *ClientOSS) IsExist(bucketName string) bool {
	isExist, err := os.IsBucketExist(bucketName)
	if err != nil {
		logIO.Error(err)
		return false
	}
	return isExist
}
func (os *ClientOSS) GetListBuckets() ([]oss.BucketProperties, error) {
	var list []oss.BucketProperties
	marker := ""
	for {
		lsRes, err := os.ListBuckets(oss.Marker(marker))
		if err != nil {
			logIO.Error(err)
			return nil, err
		}
		// 默认情况下一次返回100条记录。
		for _, bucket := range lsRes.Buckets {
			list = append(list, bucket)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return list, nil
}
