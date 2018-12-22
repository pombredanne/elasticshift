/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
)

func TestMinioCreateBucket(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = mc.CreateBucket("elasticshift", "")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// var f *os.File
	// f, err = os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	// writers := []io.Writer{NewLogWriter(mc)}
	// log.SetOutput(io.MultiWriter(writers...))

	// go dumpdata()

	// _, err = mc.PutObject("test", "file1.txt", r, "text/plain")
	// if err != nil {
	// 	panic(err)
	// }
	// err = mc.CreateBucket(&itypes.CreateBucketOptions{
	// 	Name: "test/12345&LogWriter{},
	// })

	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }
	// ch := make(chan int)

	// <-ch
}

func TestMinioPutObject(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	r := strings.NewReader("this is sample text.")
	_, err = mc.PutObject("elasticshift", "logs/12345/test2.log", r, "text/plain")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestMinioPutFObject(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	_, err = mc.PutFObject("elasticshift", "cache/12345/.cache", "/Users/ghazni/sample.txt", "text/plain")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestMinioGetFObject(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	src := "cache/5a3a41f08011e098fb86b41f/github.com/nshahm/hybrid.test.runner/master/metadata"
	err = mc.GetFObject("elasticshift", src, "/tmp/cache/12345/.cache")
	// err = mc.GetFObject("elasticshift", "cache/12345/.cache", "/tmp/cache/12345/.cache")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestMinioGetObject(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	r, err := mc.GetObject("elasticshift", "cache/12345/.cache")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	newStr := buf.String()

	fmt.Println("Content=", newStr)
}

func TestMinioSetupStorage(t *testing.T) {

	mc, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	bucketName := "test7"
	workerURL := "http://127.0.0.1:9000/worker/worker-v0.0.1-alpha.tar.gz"

	wpath, err := mc.SetupStorage(bucketName, workerURL)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println("WorkerPath = ", wpath)
}

func connectToMinio() (StorageInterface, error) {

	loggr, err := logger.New("info", "text")
	if err != nil {
		return nil, err
	}
	l := loggr.GetLogger("minio_test")

	opts := types.Storage{}
	ms := &types.MinioStorage{
		Host:        "http://127.0.0.1:9000",
		Certificate: "",
		AccessKey:   "AKIAIOSFODNN7EXAMPLE",
		SecretKey:   "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}
	opts.StorageSource = types.StorageSource{Minio: ms}

	return ConnectMinio(l, opts)
}

func testCM(t *testing.T) {
	_, err := connectToMinio()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
