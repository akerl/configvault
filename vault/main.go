package vault

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Path describes a specific object in the vault
type Path struct {
	Key    string
	User   string
	Public bool
	Bucket string
}

// Query describes a key to find in a vault
type Query struct {
	Key    string
	Public bool
	Bucket string
}

var category = map[bool]string{
	true:  "public",
	false: "private",
}

// Read returns the value for a Path
func Read(p Path) (string, error) {
	c, err := aws.GetClient()
	if err != nil {
		return "", err
	}

	key := category[p.Public] + "/" + p.User + "/" + p.Key

	buffer := []byte{}
	writer := manager.NewWriteAtBuffer(buffer)

	downloader := manager.NewDownloader(c)
	_, err = downloader.Download(context.TODO(), writer, &s3.GetObjectInput{
		Bucket: &p.Bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

// Write sets the value for a Path
func Write(p Path, data string) error {
	c, err := aws.GetClient()
	if err != nil {
		return err
	}

	key := category[p.Public] + "/" + p.User + "/" + p.Key

	uploader := manager.NewUploader(c)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &p.Bucket,
		Key:    &key,
		Body:   strings.NewReader(data),
	})
	return err
}

// Search returns a list of users that match a Query
func Search(q Query) ([]string, error) {
	c, err := aws.GetClient()
	if err != nil {
		return []string{}, err
	}

	prefix := category[q.Public] + "/"

	paginator := s3.NewListObjectsV2Paginator(c, &s3.ListObjectsV2Input{
		Bucket: &q.Bucket,
		Prefix: &prefix,
	})

	users := []string{}

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return []string{}, err
		}
		for _, obj := range page.Contents {
			parts := strings.SplitN(*obj.Key, "/", 3)
			if parts[2] == q.Key {
				users = append(users, parts[1])
			}
		}
	}
	return users, nil
}

type awsHelper struct {
	client *s3.Client
}

func (a *awsHelper) GetClient() (*s3.Client, error) {
	if a.client == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		a.client = s3.NewFromConfig(cfg)
	}
	return a.client, nil
}

var aws = awsHelper{}
