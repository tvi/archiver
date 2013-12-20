package archiver

import "testing"
import "os"
import "fmt"
import "bytes"
import "io"

func NewTarBz2Wrapper(path io.Reader) (Archive, error) {
	return NewTarBz2(path), nil
}

var tests = []struct {
	filepath    string
	constructor func(path io.Reader) (Archive, error)
}{
	{"./testdata/data.txt.gz",
		NewTarGz, // TODO(erggo): Finish.
	},
	{"./testdata/data.txt.tar.gz",
		NewTarGz,
	},
	{"./testdata/data.txt.tar.bz2",
		NewTarBz2Wrapper,
	},
	{"./testdata/data.txt.bz2",
		NewTarBz2Wrapper, // TODO(erggo): Finish.
	},
	{"./testdata/data.txt.tar",
		NewTarBz2Wrapper, // TODO(erggo): Finish.
	},
	// {"./testdata/folder/",
	// 	NewFolder,
	// },
}

func TestOpen(t *testing.T) {
	files := make(map[string][]byte)
	files["data.txt"] = bytes.NewBufferString(`A wonderful serenity has taken possession of my entire soul, like these sweet mornings of spring which I enjoy with my whole heart. I am alone, and feel the charm of existence in this spot, which was created for the bliss of souls like mine. I am so happy, my dear friend, so absorbed in the exquisite sense of mere tranquil existence, that I neglect my talents. I should be incapable of drawing a single stroke at the present moment; and yet I feel that I never was a greater artist than now.
`).Bytes()

	for _, test := range tests {
		f, e := os.Open(test.filepath)
		if e != nil {
			t.Fatalf("%s", e)
		}
		fmt.Printf("%s \n", test.filepath)
		arch, e := test.constructor(f)
		if e != nil {
			t.Fatalf("%s", e)
		}

		rets := make(map[string][]byte)
		arch.WalkAllWithContent(func(path string, info os.FileInfo, content bytes.Buffer, err error) error {
			rets[path] = content.Bytes()
			return nil
		})

		for file, content := range files {
			cont, ok := rets[file]
			if !ok {
				t.Errorf("File not found: %s in %s", file, test.filepath)
				break
			}
			if len(cont) != len(content) {
				t.Logf("File content mismatch: %s in %s", file, test.filepath)
				t.Logf("Expecting %s, got %s", len(content), len(cont))
			}
		}
	}
}

// These tests have to be only passing cases from previous test.
var tests2 = []struct {
	filepath    string
	constructor func(path io.Reader) (Archive, error)
}{
	{"./testdata/data.txt.tar.gz",
		NewTarGz,
	},
	{"./testdata/data.txt.tar.bz2",
		NewTarBz2Wrapper,
	},
}

func TestFindFile(t *testing.T) {
	files := make(map[string][]byte)
	files["data.txt"] = bytes.NewBufferString(`A wonderful serenity has taken possession of my entire soul, like these sweet mornings of spring which I enjoy with my whole heart. I am alone, and feel the charm of existence in this spot, which was created for the bliss of souls like mine. I am so happy, my dear friend, so absorbed in the exquisite sense of mere tranquil existence, that I neglect my talents. I should be incapable of drawing a single stroke at the present moment; and yet I feel that I never was a greater artist than now.
`).Bytes()

	for _, test := range tests2 {
		f, _ := os.Open(test.filepath)
		arch, _ := test.constructor(f)
		data, e := arch.GetFile("data.txt")
		if e != nil {
			t.Errorf("%s", e)
		}
		if len(data.Bytes()) != len(files["data.txt"]) {
			t.Errorf("Files are not equal.")
		}
	}
}
