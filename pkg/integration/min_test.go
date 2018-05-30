/*
Copyright 2017 The Elasticshift Authors.
*/
package integration

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	minio "github.com/minio/minio-go"
)

func testMin(t *testing.T) {

	Host := "127.0.0.1:9000"
	AccessKey := "AKIAIOSFODNN7EXAMPLE"
	SecretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	cli, err := minio.New(Host, AccessKey, SecretKey, false)
	if err != nil {
		panic(err)
	}
	// url := `http://127.0.0.1:9000/test/build.log?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20180527%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20180527T083834Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=34936c53a15f26a62c36e2f90d042a6c529edcfc71dd15bca2c0362bc6a92574`

	// Generates a url which expires in a day.
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	presignedURL, err := cli.PresignedPutObject("test", "build.log", expiry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	pr, pw := io.Pipe()

	go func(pw *io.PipeWriter) {

		data := []byte("this is the test data\n")
		for {
			pw.Write(data)
		}
	}(pw)

	// io.Copy(os.Stdout, pr)

	req, err := http.NewRequest("PUT", presignedURL.String(), pr)
	req.ContentLength = 987636783
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(message))

	// cli.PutObject("test", "build.log", pr, "application/octet-stream")

}
