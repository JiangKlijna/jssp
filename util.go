package main

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
)

func getFile(fs http.FileSystem, name string) http.File {
	f, err := fs.Open(name)
	if err != nil {
		return nil
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil
	}
	if stat.IsDir() {
		f.Close()
		return nil
	} else {
		return f
	}
}

func readFile(f http.File) ([]byte, error) {
	defer f.Close()
	return ioutil.ReadAll(f)
}

func fileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return !stat.IsDir()
	}
	return false
}

func jssp_jsjs(data []byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(`echo("`)
	isJsjs := false
	isPrint := false
	for i, n := 0, len(data); i < n; i++ {
		c := data[i]
		if isJsjs {
			if c == '%' && data[i+1] == '>' {
				if isPrint {
					isPrint = false
					buf.WriteString(`);echo("`)
				} else {
					buf.WriteString(`;echo("`)
				}
				i++
				isJsjs = false
			} else {
				buf.WriteByte(c)
			}
		} else {
			if c == '<' && data[i+1] == '%' {
				buf.WriteString(`");`)
				if data[i+2] == '=' {
					i++
					isPrint = true
					buf.WriteString(`echo(`)
				}
				i++
				isJsjs = true
			} else {
				switch c {
				case '\n':
					buf.WriteString(`\n");`)
					buf.WriteByte(c)
					buf.WriteString(`echo("`)
				case '\r':
					continue
				case '"':
					buf.WriteString(`\"`)
				default:
					buf.WriteByte(c)
				}
			}
		}
	}
	buf.WriteString(`");`)
	return buf.Bytes()
}

func getUUID() string {
	return uuid.NewV4().String()
}
