// Code generated by "esc -o bindata/esc.go -pkg=bindata templates"; DO NOT EDIT.

package bindata

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/templates/call.tmpl": {
		name:    "call.tmpl",
		local:   "templates/call.tmpl",
		size:    241,
		modtime: 1617338274,
		compressed: `
H4sIAAAAAAAC/0SOQWrDQAxFryKMFy0YHaDQA3hTSlvatRjLrsCeFo2SEITuHsY4mdWHP2/el/vEs2SG
LtG6dhHuF7FfwA9OLGfW2sgM+c8Ax/JpekoWYYbunKf6eicBI1qLb7RxxJO7Ul4Yehmg5xVeXgHfSWlj
Yy2HvZeIAR5/296PitUbzJB0KU2/K+riTuPX9Z9xLN+kQpOkCMTG7vF85C0AAP//ZQi8iPEAAAA=
`,
	},

	"/templates/function.tmpl": {
		name:    "function.tmpl",
		local:   "templates/function.tmpl",
		size:    3816,
		modtime: 1617453812,
		compressed: `
H4sIAAAAAAAC/7RX3W/bNhB/lv+Kq5EWUqCy2B4d+GFt0yHAmhRpsA4ohkKVT64QmfLIU4KA4P8+HElJ
lD+SYB8vlkTe9+93R9qYFVa1RJhXnSypbuXc2pkxr+GkgsUShLWzGW+BMeIGNV0WG7Q2JTgl1FTLtbjJ
wMwSVrmv6QeIayyxvkNl7Sxxy3UF4kJ/JtWV5BaH1Q81Nivt1xJ62CJUbgW0E2a7QVoVco07CgAAo/0L
SaiqosSwm7hdp8BZwOmmLW+/GSMGQX0hq1Z8ul37nMTHtrzd3/eb3qQx2OjIQWyf6/OwjUXlKnxwkOPX
UIF+KXrfeeXMuOifClVskFC51F2hCrWelCkq0r6Gc+iW9kKNPE79M7yaKWCMi4PrsLJ2U2y/alK1XP/p
vfc1+Tp+s4EQFKumsqWgn1kriw2CtzDW6Ah7+iIUcjVSaIcFgTH+4Va2CreFQmDapmH71D9zX7dT/s0G
HEZMjen972N6BLJHgEoS541/DuhEaF2j7hrSvZ8vhaTHgBpcXiN1SupzpdpQsPtC0rlS8L1tmx10GZA3
b+Dm6v3VAn5ZrYARhrLQqAWLsjqDs4A5L/40z48DM229uLXH9QlKAZPFk6CEnl+33K7iQl6pFar0X8+B
kw+L5dCr0Xqwt9PzH5F+tCt+/a3WFClA4Jlgi9aK8z8+nb+7STNhTBgU6WDTcaGPMUkYTGvz8cvhkgUM
RzX//aTaTb1Bnf6c5VEy44yJxw9M91wWmZ9EeY/YsBu3ZP/GYnaWVK0CY7ifWwXic/fdTwjf3Rf6srjF
VZY51kbzghnVj4hvwWYORH6y+FHhM/fmGH6ixZJoliQlqYYFAxsu8f5dK0m1TYMqJU5ihRUqYDnxoZa1
/pFms4EQ/wVx/znpiMRAlvGYWMKT59Al3j92FKWcbfYUzEldcQT9LHyxBFk3obkg2klfDXHmwO/DaEzs
s9jRF6EnxGSRe6BpsLHWI050tqfphnsvOOENa/DzzBieYuC0SFx33C5HOEYkwqdnmps3e7eVgxEO72k2
DTLK//AxdeSas3eeOM67s+Jhi064UNa+Ct7CtBe/F02H1prexGEC8k5gxOIw1/LRQDQZDlxJ9j7iCTa9
RPRpflE1DdlPLheLJbz6/kCoxduuqlCZ5zgMNPCz5Eo2D/HJlu2vX0l0VcpgiIxws20KQpgrf5rO4aRy
02jcKYum8cvHojhwpHInedR2A7MWUCmP6iEnZ8OES1nO92DGTyLRn9QBZhLOZJXOY1sb1LpYY0gFWQKW
8PIuh1795d08n7iv5bYbkkel8shZNjKiv7FMrh5uL25msFa5hPv2KltJtexw57B59H776G3nKKdc1X9t
aWycgWPis7s+ptlZJOKrGl+eRt41GoOPt4Wuy/ie3oN7Uh3iFzflJIa4zk0tcRfoZ8fzP/l/obBqsCTx
HnF7/ldXNOlgIZ8GlMURDeg9h4d9wCHYj11D9baZBBviGbn6BFGPBnn8L9MOT4ET8iLw2v+dsbNZz9O/
AwAA//+ruFUU6A4AAA==
`,
	},

	"/templates/header.tmpl": {
		name:    "header.tmpl",
		local:   "templates/header.tmpl",
		size:    262,
		modtime: 1617444815,
		compressed: `
H4sIAAAAAAAC/2SOMW7EIBBFe04xokqK4EOkchOlSB/NwhgjG7AwW43m7iuwvVssFZr/3tdndjSFRKBn
QkdFiyjmgskTmO8cI6W6izCbHlBy8CWiNrQLegJm83t8RZQKcculwserYeyXXUQBQKN/MJII6C7WWUSf
SZjAjPuYKpUJLZ1GezHb5Z/Z/GFY3/ShpVcHJXetbHO0D3W+34zNcfB5xeQ7Pfh8SJ9P9BEAAP//6YLz
QwYBAAA=
`,
	},

	"/templates/inline.tmpl": {
		name:    "inline.tmpl",
		local:   "templates/inline.tmpl",
		size:    49,
		modtime: 1617338274,
		compressed: `
H4sIAAAAAAAC/6quTklNy8xLVVDKzMvJzEtVqq1VqK4uSc0tyEksSVVQSk7MyVFS0AOLpual1NYCAgAA
//+q60H/MQAAAA==
`,
	},

	"/templates/inputs.tmpl": {
		name:    "inputs.tmpl",
		local:   "templates/inputs.tmpl",
		size:    177,
		modtime: 1617338274,
		compressed: `
H4sIAAAAAAAC/0yNMaoDMQxE+38KsWz58QECOUCaEMgJFCwvLqwESVsJ3T1YpHAlzWN4416pdSbYOn9O
0y3CfW9wuUKZb2/Ab4PyPF9GarqyOw6qEWbFnbhGMA76h1/I3t7KQzrbLeUTCvJByVFwkJFoKlAOLe5J
5/TiWc/fNwAA//94+RPrsQAAAA==
`,
	},

	"/templates/message.tmpl": {
		name:    "message.tmpl",
		local:   "templates/message.tmpl",
		size:    201,
		modtime: 1617338274,
		compressed: `
H4sIAAAAAAAC/zyN4WqDQBCE//sUiyi0oPsAhT5A/xRpS/9f4mgW9GLuTkNY9t2DB/HXDDPDN6o9BvGg
ckaMbkRJrVmhKgP5ayL+XU8JMUWz+sakCt+bqd4lXYh/cIZsCHvCf48F/O+mFWZ8DPnbzTB7y0Tugvj0
5Zd1B6oG50dQJQ1VmOjjk7hzwc1ICLmXgSoxa16/9XZws7wXqi1l+wwAAP//kC65UskAAAA=
`,
	},

	"/templates/results.tmpl": {
		name:    "results.tmpl",
		local:   "templates/results.tmpl",
		size:    168,
		modtime: 1617338274,
		compressed: `
H4sIAAAAAAAC/1yNTQrCQAyFr/Iosyw9gOBS3HsDoRkJlAy8ma5C7i6pRcFVfr4vee6rVDXBROn7NvoU
AXc+7SUoOqPIhssVy+ODI9y1omjEDHexNTf3NrBkc85a82DstH4jG1MW8uQ4hMbv0385A3/uUd8BAAD/
/7BPz2GoAAAA
`,
	},

	"/templates": {
		name:  "templates",
		local: `templates`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"templates": {
		_escData["/templates/call.tmpl"],
		_escData["/templates/function.tmpl"],
		_escData["/templates/header.tmpl"],
		_escData["/templates/inline.tmpl"],
		_escData["/templates/inputs.tmpl"],
		_escData["/templates/message.tmpl"],
		_escData["/templates/results.tmpl"],
	},
}
