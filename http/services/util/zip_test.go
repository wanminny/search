package util

import "testing"

func TestZipDir(t *testing.T) {
	ZipDir("/www/bootstrap/","/gopath/ziptest/bbbb.zip")
}
