package feeds

// func FeedRunner() {
// 	sess, err := session.NewSession()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	sss := s3.New(sess)
// 	logger := logging.NewLogger(false, nil)
//
// 	feedToResolve := make(chan artarchive.Slide, 0)
// 	resolveToImageArchive := make(chan artarchive.Slide, 0)
// 	archiveToUpload := make(chan artarchive.Slide, 0)
// 	uploadToConsumer := make(chan artarchive.Slide, 0)
//
// 	// Fetch from feeds
// 	rssFetcher := NewFeedToSlideProducer([]string{
// 		"https://feedbin.com/starred/Cepxc9l63Bbn0RKef9J3MQ.xml",
// 		"https://feedbin.com/starred/3d5um7AVLzNCL-mMtMxKeg.xml",
// 	}, feedToResolve)
//
// 	// Resolve found slides with known versions
// 	slideStorage := NewSlideStorage(sss, "art.rumproarious.com", "v2")
// 	resolveTransform := NewSlideResolverTransform(feedToResolve, resolveToImageArchive, slideStorage)
//
// 	// Archive new images
// 	s3Upload := s3manager.NewUploader(sess)
// 	imageUploader := MustNewImageUploader(s3Upload, sss, "images", "art.rumproarious.com")
// 	imageArchiver := NewSlideImageUploader(resolveToImageArchive, archiveToUpload, imageUploader)
//
// 	// Upload slides
// 	slideUploader := NewSlideUploader(archiveToUpload, uploadToConsumer, slideStorage)
//
// 	// Dump things
// 	debugConsumer := NewDebugSlideConsumer(logger, uploadToConsumer)
// 	pipeline := NewPipeline(rssFetcher, debugConsumer, resolveTransform, imageArchiver, slideUploader)
// 	pipeline.Run()
// }
