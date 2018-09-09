package artarchive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestBuildKey(t *testing.T) {
	var testTable = []struct {
		prefix   string
		guidHash string
		output   string
	}{
		{"", "a", "slides/a/data.json"},
		{"b", "a", "b/slides/a/data.json"},
	}

	for _, tt := range testTable {
		slide := Slide{
			GUIDHash: tt.guidHash,
		}
		o := buildKey(tt.prefix, slide)
		assert.Equal(t, tt.output, o)
	}
}

type TestS3Interface struct {
	body io.ReadCloser
	err  error
}

func (si *TestS3Interface) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, nil
}

func (si *TestS3Interface) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: si.body,
	}, si.err
}

func (si *TestS3Interface) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}

func TestResolve(t *testing.T) {
	var testTable = []struct {
		body     string
		GUIDHash string
		err      error
	}{
		{"{\"guid_hash\": \"bbb\"}", "bbb", nil},
		{"{\"guid_hash\": \"bbb\"}", "bbb", fmt.Errorf("Yo")},
	}

	for _, tt := range testTable {
		sss := &TestS3Interface{
			body: ioutil.NopCloser(bytes.NewBufferString(tt.body)),
			err:  tt.err,
		}
		ss := NewSlideStorage(sss, "my.bucket.com", tt.GUIDHash)
		slide := Slide{
			GUIDHash: tt.GUIDHash,
		}
		resp := ss.Resolve(slide)
		assert.Equal(t, slide, resp)
	}

}
