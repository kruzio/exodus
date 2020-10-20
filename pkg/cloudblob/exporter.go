package cloudblob

import (
	"context"
	"gocloud.dev/blob"
	log "k8s.io/klog"

	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"
)

type CloudBlob struct {
	BucketUrl string
	DestName  string
}

func (c *CloudBlob) Export(data []byte) error {
	bucket, err := blob.OpenBucket(context.Background(), c.BucketUrl)
	if err != nil {
		log.Errorf("Failed to open bucket %v. %v", c.BucketUrl, err)
		return err
	}

	defer func() {
		err := bucket.Close()
		if err != nil {
			log.Errorf("Failed to close bucket %v. %v", c.BucketUrl, err)
		}
	}()

	bucketWriter, err := bucket.NewWriter(context.Background(), c.DestName, nil)
	if err != nil {
		log.Errorf("Failed to initialize bucket writer for bucket %v. %v", c.BucketUrl, err)
		return err
	}

	_, err = bucketWriter.Write(data)
	if err != nil {
		log.Errorf("Failed to write Advisor Report to bucket %v. %v", c.BucketUrl, err)
		return err
	}

	err = bucketWriter.Close()
	if err != nil {
		log.Errorf("Failed to close writer for bucket %v. %v", c.BucketUrl, err)
		return err
	}

	log.V(5).Infof("Successfully written %v bytes to bucket URL %v at bucket key %v", len(data), c.BucketUrl, c.DestName)

	return nil
}

func (c *CloudBlob) SetUploadUrl(url string) error {
	c.BucketUrl = url

	return nil
}

func (c *CloudBlob) SetDestName(name string) error {
	c.DestName = name

	return nil
}
