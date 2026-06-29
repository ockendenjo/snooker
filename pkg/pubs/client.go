package pubs

import (
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewClient(s3Client *s3.Client, bucket string) Client {
	return &client{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

type Client interface {
	GetPubGroups(ctx context.Context) (GroupMap, error)
	GetAllPubs(ctx context.Context) ([]*Pub, error)
}

type client struct {
	s3Client *s3.Client
	bucket   string
}

func (c *client) GetAllPubs(ctx context.Context) ([]*Pub, error) {
	res, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    new("pubs.json"),
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var pf pubFile
	if err = json.NewDecoder(res.Body).Decode(&pf); err != nil {
		return nil, err
	}
	return pf.Pubs, nil
}

func (c *client) GetPubGroups(ctx context.Context) (GroupMap, error) {
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    new("groups.json"),
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(result.Body)

	var gf groupsFile
	if err = json.NewDecoder(result.Body).Decode(&gf); err != nil {
		return nil, err
	}
	return gf.Groups, nil
}
