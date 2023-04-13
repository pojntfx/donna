package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _migrations_1681420952_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8d\xc1\x0d\x02\x31\x0c\x04\xff\xae\x62\x9f\x20\xb8\x0a\xee\x4b\x0b\x14\x60\x12\x0b\x45\x24\x76\xe4\x18\xc1\x75\x8f\x80\x0f\xe1\xb7\xda\x19\x69\x96\x05\x87\x56\xae\xce\x21\x38\x77\x4a\x2e\xef\x15\x7c\xa9\x82\x64\x1a\x9c\x62\x60\x47\x00\x50\x32\x42\x9e\x81\xee\xa5\xb1\x6f\xb8\xc9\x06\xb5\x80\xde\x6b\x3d\x7e\x0c\xe5\x26\x5f\x67\xfe\x39\x67\x97\x31\x66\x44\xfb\x95\x7e\xeb\x27\x7b\x28\x65\xb7\xfe\x57\x5f\x5f\x01\x00\x00\xff\xff\x31\x48\x4a\x73\xa3\x00\x00\x00")

func migrations_1681420952_sql() ([]byte, error) {
	return bindata_read(
		_migrations_1681420952_sql,
		"../migrations/1681420952.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"../migrations/1681420952.sql": migrations_1681420952_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"..": &_bintree_t{nil, map[string]*_bintree_t{
		"migrations": &_bintree_t{nil, map[string]*_bintree_t{
			"1681420952.sql": &_bintree_t{migrations_1681420952_sql, map[string]*_bintree_t{
			}},
		}},
	}},
}}
