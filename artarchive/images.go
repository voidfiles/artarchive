package artarchive

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	// _ "golang.org/x/image/webp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var acceptableContentType = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var fileEndingToContentType = map[string]string{
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"jpg":  "image/jpg",
	"jpeg": "image/jpg",
}

var fileContentTypeToEnding = map[string]string{
	"image/png":  "png",
	"image/gif":  "gif",
	"image/webp": "webp",
	"image/jpg":  "jpg",
}

type S3UploadInterface interface {
	Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type ImageUploader struct {
	s3upload S3UploadInterface
	sss      S3Interface
	bucket   string
	prefix   string
}

func MustNewImageUploader(s3upload S3UploadInterface, sss S3Interface, prefix, bucket string) *ImageUploader {
	return &ImageUploader{
		s3upload: s3upload,
		sss:      sss,
		bucket:   bucket,
		prefix:   prefix,
	}
}

var tumblrImage = regexp.MustCompile(`([^/]+)(/.*)_([\d]{0,3}\.)(jpg|jpeg|png|gif)`)

func FixTumblrURL(host string, parts []string) string {
	return fmt.Sprintf("https://%s/%s%s_1280.%s", host, parts[1], parts[2], parts[4])
}

func ParseImageURL(imageURL string) string {
	parsedImageURL, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}
	if tumblrImage.MatchString(parsedImageURL.Path) {
		return FixTumblrURL(parsedImageURL.Host, tumblrImage.FindStringSubmatch(parsedImageURL.Path))
	}

	if strings.HasPrefix(imageURL, "http://feeds.feedburner.com") {
		return ""
	}

	return imageURL
}

func (iu *ImageUploader) Archive(metadata map[string]*string, imageURL string) (*ImageInfo, error) {
	res, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Failed to fetch imageUrl %s", imageURL)
		return nil, err
	}
	defer res.Body.Close()

	tmpfile, err := ioutil.TempFile("", "images")
	if err != nil {
		log.Printf("Failed to create tmpfile %s", err)
		return nil, err
	}
	defer tmpfile.Close()
	defer func() {
		os.Remove(tmpfile.Name())
	}()
	h := sha256.New()
	mWriter := io.MultiWriter(tmpfile, h)

	io.Copy(mWriter, res.Body)
	tmpfile.Seek(0, 0)
	decodedImage, ending, err := image.DecodeConfig(tmpfile)
	if err != nil {
		log.Printf("Failed to decode image %s err %s", imageURL, err)
		return nil, nil
	}
	if ending == "jpeg" {
		ending = "jpg"
	}
	filename := fmt.Sprintf("%x.%s", h.Sum(nil), ending)
	imageInfo := &ImageInfo{
		URL:         imageURL,
		Width:       decodedImage.Width,
		Height:      decodedImage.Height,
		ContentType: fmt.Sprintf("image/%s", ending),
		Filename:    filename,
	}

	objHead, _ := iu.sss.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(iu.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", iu.prefix, imageInfo.Filename)),
	})

	if objHead != nil && objHead.ETag != nil {
		log.Printf(
			"Already downloaded image %s to: %s contentType: %s",
			imageURL,
			fmt.Sprintf("%s/%s", iu.prefix, imageInfo.Filename),
			imageInfo.ContentType)
		return imageInfo, nil
	}

	metadata["width"] = aws.String(fmt.Sprintf("%d", imageInfo.Width))
	metadata["height"] = aws.String(fmt.Sprintf("%d", imageInfo.Width))
	log.Printf("Saving image from: %s to: %s", imageURL, fmt.Sprintf("%s/%s", iu.prefix, imageInfo.Filename))
	_, err = iu.s3upload.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(iu.bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", iu.prefix, imageInfo.Filename)),
		Body:        tmpfile,
		ContentType: aws.String(imageInfo.ContentType),
		Metadata:    metadata,
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		log.Printf("Failed to upload key %v", imageInfo)
		return imageInfo, err
	}

	return imageInfo, nil
}

type SlideImageUploader struct {
	slideChanIn  chan Slide
	slideChanOut chan Slide
	iu           *ImageUploader
}

func NewSlideImageUploader(slideChanIn chan Slide, slideChanOut chan Slide, iu *ImageUploader) *SlideImageUploader {
	return &SlideImageUploader{
		slideChanIn:  slideChanIn,
		slideChanOut: slideChanOut,
		iu:           iu,
	}
}

func (sc *SlideImageUploader) Run() {
	for slide := range sc.slideChanIn {
		if slide.ArchivedImage == nil {
			imageInfo, err := sc.iu.Archive(make(map[string]*string), slide.SourceImageURL)
			if err != nil {
				log.Printf("Error while upload image: %v", err)
			} else {
				slide.ArchivedImage = imageInfo
			}
		}
		sc.slideChanOut <- slide
	}

	close(sc.slideChanOut)
}
