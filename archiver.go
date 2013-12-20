// Package archiver provides an easy way of accessing files inside archives
// such as tar.gz.
package archiver

import "archive/tar"
import "bytes"
import "compress/bzip2"
import "compress/gzip"
import "fmt"
import "io"
import "os"
import "path/filepath"

// WalkFunc is a function type for interfacing with client code.
type WalkFunc func(path string, info os.FileInfo, content bytes.Buffer, err error) error

// Archive is a generic interface representing any archive.
// Supported suffixes right now are:
// .tgz	                .tar.gz
// .tbz, .tbz2 & .tb2   .tar.bz2
// TODO suffixes:
// .taz	                .tar.Z
// .tlz	                .tar.lz & .tar.lzma
// .txz	                .tar.xz
// .zip
// .rar ???
type Archive interface {
	WalkAllWithContent(walkFn WalkFunc) error
	WalkWithContent(root string, walkFn WalkFunc) error
	GetFile(path string) (*bytes.Buffer, error)
}

type tarBased struct {
	reader *tar.Reader
}

// NewTarBz2 is a contructor that takes a file of type 
// .tbz, .tbz2, .tb2 or .tar.bz2 and returns an initialised archive.
func NewTarBz2(path io.Reader) Archive {
	reader := tar.NewReader(bzip2.NewReader(path))
	return &tarBased{reader}
}

// WalkAllWithContent is a method that walks over all files inside an archive
// and executes given walkFn on each of them.
func (a *tarBased) WalkAllWithContent(walkFn WalkFunc) error {
	var err error
	for head, err := a.reader.Next(); err == nil; head, err = a.reader.Next() {
		var buffer bytes.Buffer
		r := []byte{0}
		for _, eof := a.reader.Read(r); eof != io.EOF; _, eof = a.reader.Read(r) {
			buffer.Write(r)
		}
		err = walkFn(head.Name, head.FileInfo(), buffer, err)
		if err != nil {
			return fmt.Errorf("could not walk: %s, because of %s", head.Name, err)
		}
	}
	if err != nil {
		return fmt.Errorf("walking error: %s ", err)
	}
	return nil
}

// WalkWithContent is a method that walks over all files in particular 
// subfolder inside an archive and executes give walkFn on each file.
// TODO(erggo): Rewrite.
func (a *tarBased) WalkWithContent(root string, walkFn WalkFunc) error {
	return a.WalkAllWithContent(func(path string, info os.FileInfo, content bytes.Buffer, err error) error {
		// TODO(erggo): Replace with function isParent.
		if filepath.HasPrefix(filepath.Clean(path), filepath.Clean(root)) {
			err = walkFn(path, info, content, err)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetFile is a method that returns a content of given file inside an archive.
// TODO(erggo): Improve efficiency.
func (a *tarBased) GetFile(path string) (*bytes.Buffer, error) {
	var ret *bytes.Buffer
	e := a.WalkAllWithContent(func(actualPath string, info os.FileInfo, content bytes.Buffer, err error) error {
		if filepath.Clean(actualPath) == filepath.Clean(path) {
			ret = &content
		}
		return nil
	})

	if e != nil {
		return nil, e
	}
	if ret == nil {
		return nil, fmt.Errorf("file not found")
	}
	return ret, nil
}

// NewTarGz is a contructor that takes a file of type 
// .tgz or .tar.gz and returns an initialised archive.
func NewTarGz(path io.Reader) (Archive, error) {
	g, e := gzip.NewReader(path)
	if e != nil {
		return nil, e
	}
	reader := tar.NewReader(g)
	return &tarBased{reader}, nil
}

// // NewArchive is a constructor, which tries to guess compression algorithm.
// func NewArchive(path io.Reader) (Archive, error) {
// 	return nil, nil
// }

// // NewFolder is a constructor that returns internal representation of
// // an archive from folder.
// func NewFolder(path io.Reader) (Archive, error) {
// 	return nil, nil
// }
