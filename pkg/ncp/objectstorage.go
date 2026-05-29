package ncp

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// NCP Object Storage is S3-compatible. Endpoint/region default to the KR region
// and can be overridden via environment variables for other regions.
const (
	defaultObjectStorageEndpoint = "https://kr.object.ncloudstorage.com"
	defaultObjectStorageRegion   = "kr-standard"
)

// Bucket represents an Object Storage bucket.
type Bucket struct {
	Name string
}

func objectStorageEndpoint() string {
	if v := os.Getenv("NCP_OBJECT_STORAGE_ENDPOINT"); v != "" {
		return v
	}
	return defaultObjectStorageEndpoint
}

func objectStorageRegion() string {
	if v := os.Getenv("NCP_OBJECT_STORAGE_REGION"); v != "" {
		return v
	}
	return defaultObjectStorageRegion
}

// newS3Client builds an S3-compatible client pointed at NCP Object Storage,
// reusing the account's API access/secret key.
func (c *Client) newS3Client() *s3.Client {
	cfg := aws.Config{
		Region:      objectStorageRegion(),
		Credentials: credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, ""),
	}
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(objectStorageEndpoint())
		o.UsePathStyle = true
	})
}

// ListBuckets returns all Object Storage buckets for the account.
func (c *Client) ListBuckets() ([]Bucket, error) {
	cli := c.newS3Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	out, err := cli.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, b := range out.Buckets {
		buckets = append(buckets, Bucket{Name: aws.ToString(b.Name)})
	}
	return buckets, nil
}

// DeleteBucket empties a bucket (all object versions, delete markers and
// in-progress multipart uploads) and then deletes the bucket itself.
func (c *Client) DeleteBucket(bucket string, logFn func(string)) error {
	cli := c.newS3Client()

	if err := c.abortMultipartUploads(cli, bucket, logFn); err != nil {
		logFn(fmt.Sprintf("    [경고] 멀티파트 업로드 정리 오류: %v", err))
	}
	if err := c.deleteAllObjectVersions(cli, bucket, logFn); err != nil {
		return fmt.Errorf("객체 삭제: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := cli.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return fmt.Errorf("버킷 삭제: %w", err)
	}
	return nil
}

// deleteAllObjectVersions removes every object version and delete marker in a
// bucket. ListObjectVersions covers both versioned and non-versioned buckets
// (non-versioned objects are returned with a "null" version id).
func (c *Client) deleteAllObjectVersions(cli *s3.Client, bucket string, logFn func(string)) error {
	paginator := s3.NewListObjectVersionsPaginator(cli, &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
	})

	deleted := 0
	for paginator.HasMorePages() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		page, err := paginator.NextPage(ctx)
		cancel()
		if err != nil {
			return err
		}

		var objs []types.ObjectIdentifier
		for _, v := range page.Versions {
			objs = append(objs, types.ObjectIdentifier{Key: v.Key, VersionId: v.VersionId})
		}
		for _, m := range page.DeleteMarkers {
			objs = append(objs, types.ObjectIdentifier{Key: m.Key, VersionId: m.VersionId})
		}
		if len(objs) == 0 {
			continue
		}

		// DeleteObjects accepts up to 1000 keys per call; a single page is <= 1000.
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		out, err := cli.DeleteObjects(dctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &types.Delete{Objects: objs, Quiet: aws.Bool(true)},
		})
		dcancel()
		if err != nil {
			return err
		}
		for _, e := range out.Errors {
			logFn(fmt.Sprintf("    [경고] 객체 삭제 실패 %s: %s", aws.ToString(e.Key), aws.ToString(e.Message)))
		}
		deleted += len(objs)
		logFn(fmt.Sprintf("    객체 %d개 삭제됨...", deleted))
	}
	return nil
}

// abortMultipartUploads aborts any in-progress multipart uploads so the bucket
// can be deleted.
func (c *Client) abortMultipartUploads(cli *s3.Client, bucket string, logFn func(string)) error {
	paginator := s3.NewListMultipartUploadsPaginator(cli, &s3.ListMultipartUploadsInput{
		Bucket: aws.String(bucket),
	})
	for paginator.HasMorePages() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		page, err := paginator.NextPage(ctx)
		cancel()
		if err != nil {
			return err
		}
		for _, u := range page.Uploads {
			actx, acancel := context.WithTimeout(context.Background(), 30*time.Second)
			_, err := cli.AbortMultipartUpload(actx, &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      u.Key,
				UploadId: u.UploadId,
			})
			acancel()
			if err != nil {
				logFn(fmt.Sprintf("    [경고] 멀티파트 업로드 중단 실패 %s: %v", aws.ToString(u.Key), err))
			}
		}
	}
	return nil
}
