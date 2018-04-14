package main

import (
	"archive/tar"
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"os"
)

type compressUtil struct {
}
func (*compressUtil)compressFiles(files []string,output string)bool{
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWrite := tar.NewWriter(gzipWriter)
	for i := range files{
		fileInfo,err := os.Stat(files[i])
		if err != nil {
			return false
		}
		header,err := tar.FileInfoHeader(fileInfo,"") //put files in root directory
		tarWrite.WriteHeader(header)
		content,err := ioutil.ReadFile(files[i])
		if err != nil {
			return true

		}
		tarWrite.Write(content)
	}
	tarWrite.Close()
	gzipWriter.Close()
	os.Remove(output + ".tar")
	err := ioutil.WriteFile(output + ".tar",buffer.Bytes(),0666)
	return  err == nil
}