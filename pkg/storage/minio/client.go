package minio

import (
	"io"

	"github.com/minio/minio-go/v6"
)

type Client struct {
	c      *minio.Client
	region string
	bucket string
}

func New(endpoint, accessKeyID, secretAccessKey, region, bucket string, useSSL bool) (*Client, error) {
	c, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}
	return &Client{c: c, region: region, bucket: bucket}, nil
}

func (c *Client) Upload(name string, r io.Reader) (size int64, err error) {
	return c.c.PutObject(c.bucket, name, r, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
}

func (c *Client) Download(name string) (io.Reader, error) {
	return c.c.GetObject(c.bucket, name, minio.GetObjectOptions{})
}
