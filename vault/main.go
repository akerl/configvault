package vault

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

func lookupUser(given string) (string, error) {
	if given != "" {
		return given, nil
	}
	c, err := ah.GetStsClient()
	if err != nil {
		return "", err
	}
	resp, err := c.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	parts := strings.Split(*resp.Arn, "/")
	return parts[len(parts)-1], nil
}

// Read returns the value for a Path
func Read(p Path) (string, error) {
	c, err := ah.GetS3Client()
	if err != nil {
		return "", err
	}

	user, err := lookupUser(p.User)
	if err != nil {
		return "", err
	}

	key := category[p.Public] + "/" + user + "/" + p.Key

	writer := manager.NewWriteAtBuffer([]byte{})

	downloader := manager.NewDownloader(c)
	_, err = downloader.Download(context.TODO(), writer, &s3.GetObjectInput{
		Bucket: &p.Bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	return string(writer.Bytes()), nil
}

// Write sets the value for a Path
func Write(p Path, data string) error {
	c, err := ah.GetS3Client()
	if err != nil {
		return err
	}

	user, err := lookupUser(p.User)
	if err != nil {
		return err
	}

	key := category[p.Public] + "/" + user + "/" + p.Key

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
	c, err := ah.GetS3Client()
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
	config    *aws.Config
	s3Client  *s3.Client
	stsClient *sts.Client
}

func (a *awsHelper) getConfig() (*aws.Config, error) {
	if a.config == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		a.config = &cfg
	}
	return a.config, nil
}

func (a *awsHelper) GetS3Client() (*s3.Client, error) {
	if a.s3Client == nil {
		cfg, err := a.getConfig()
		if err != nil {
			return nil, err
		}
		a.s3Client = s3.NewFromConfig(*cfg)
	}
	return a.s3Client, nil
}

func (a *awsHelper) GetStsClient() (*sts.Client, error) {
	if a.stsClient == nil {
		cfg, err := a.getConfig()
		if err != nil {
			return nil, err
		}
		a.stsClient = sts.NewFromConfig(*cfg)
	}
	return a.stsClient, nil
}

var ah = awsHelper{}
