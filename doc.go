/*
Package archiver provides an easy way of accessing files inside archives
such as tar.gz.

Example

	package main

	import "fmt"
	import "os"
	import "github.com/erggo/archiver"

	func main() {
		f, e := os.Open("archive.tar.bz2")
		if e != nil {
			fmt.Printf("%s \n", e)
			return
		}
		tar := archiver.NewTarBz2(f)
		info, e := tar.GetFile("file.xml")
		if e != nil {
			fmt.Printf("%s \n", e)
			return
		}
		fmt.Printf("%s", info)
	}

*/
package archiver