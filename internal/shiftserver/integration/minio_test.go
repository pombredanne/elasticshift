/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"testing"

	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
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

func connectToMinio() (StorageInterface, error) {

	loggr, err := logger.New("info", "text")
	if err != nil {
		return nil, err
	}
	l := loggr.GetLogger("minio_test")

	opts := types.Storage{}
	ms := &types.MinioStorage{
		Host:        "127.0.0.1:9000",
		Certificate: "",
		AccessKey:   "AKIAIOSFODNN7EXAMPLE",
		SecretKey:   "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}
	opts.StorageSource = types.StorageSource{Minio: ms}

	return ConnectMinio(l, opts)
}

type LogWriter struct {
	storage      StorageInterface
	buf          bytes.Buffer
	chunksize    int
	bytesWritten int
	mu           sync.RWMutex
	datach       chan []byte
	// r *io.Reader
	// w *bufio.Writer
	pr *io.PipeReader
	pw *io.PipeWriter
}

func NewLogWriter(storage StorageInterface) *LogWriter {

	lw := &LogWriter{storage: storage, chunksize: 1024}
	// lw.buf = &bytes.Buffer{}
	// lw.w = bufio.NewWriter(&lw.buf)
	lw.pr, lw.pw = io.Pipe()
	lw.datach = make(chan []byte)

	go func(pw *io.PipeWriter, datach chan []byte) {

		select {
		case data := <-datach:
			pw.Write(data)
			fmt.Println(string(data))
		}
	}(lw.pw, lw.datach)

	go lw.putobject()
	return lw
}

func (lw LogWriter) putobject() {
	fmt.Println("put object started")
	_, err := lw.storage.PutObject("test", "file1.txt", lw.pr, "text/plain")
	if err != nil {
		panic(err)
	}
	fmt.Println("put object ended")
}

func (lw LogWriter) Write(b []byte) (int, error) {

	// lw.datach <- b

	// lw.mu.Lock()
	// defer lw.mu.Unlock()

	// fmt.Println("Writing...")
	//s := string(b)
	// fmt.Fprintf(&lw.buf, s)

	// var err error
	// ln, err := lw.buf.Write(b)
	// ln, err := lw.pw.Write(b)
	// if ln != len(b) || err != nil {
	// 	panic(err)
	// }

	// lw.bytesWritten = lw.bytesWritten + ln

	// fmt.Printf("\n---------------------written (%d)  : current len(%d)", lw.bytesWritten, ln)
	// fmt.Println(string(lw.buf.Bytes()))

	// if lw.chunksize < lw.bytesWritten {
	// fmt.Println("--------------------------------------writing to minio ----------------------------")
	//_, err := lw.storage.PutObject("test", "file1.txt", bytes.NewReader(b), -1, "text/plain")
	// lw.bytesWritten = 0
	// lw.buf.Reset()
	// }

	// fmt.Println(lw.buf.String())
	return len(b), nil
	// return lw.w.Write(b)
	// l := len(b)
	// _, err := lw.storage.PutObjectStreaming("test", "file1.txt", bytes.NewReader(b))
	// return l, err
}

func dumpdata() {

	i := 1
	for {
		log.Println(fmt.Sprintf("this is line no %d", i))
		i++
	}
}
