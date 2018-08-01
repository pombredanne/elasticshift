/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
)

func testCreateBucket(t *testing.T) {

	// logger := logrus.New()
	// logger.Out = os.Stdout

	// opts := types.Storage{
	// 	StorageSource: types.StorageSource{
	// 		&MinioStorage{
	// 			Host:      "127.0.0.1:9000",
	// 			AccessKey: "AKIAIOSFODNN7EXAMPLE",
	// 			SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	// 		},
	// 	},
	// }

	// mc, err := ConnectMinio(*logger, opts)

	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }

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
