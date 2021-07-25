package proto

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

var _proto_ds_store = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x99\x4f\x8b\xd3\x40\x18\xc6\x9f\x49\x22\x04\x45\xc9\x41\xa4\x78\x90\x80\x9f\xa0\xdf\x20\xd6\x08\xbd\xf7\xe2\x41\xa8\xb6\x53\x6a\x69\x42\x4a\xd2\x1e\x3c\x08\xfd\x16\xea\x5d\x7a\x16\x3c\x08\xf5\x22\xde\xc4\x83\x5e\x3c\x79\xd9\xef\xb0\xb0\xe7\x25\x33\xef\xb2\x09\xcd\x74\x49\xda\xfd\x43\xf7\xfd\x41\xf9\x95\x26\xef\xd3\xce\x4b\x3a\xe9\x4c\x01\x88\xce\x42\xb6\x81\x16\x00\x17\xda\xde\x3d\x54\xe2\xd2\x63\x03\x8b\xfc\x30\xcf\x03\x3c\x20\x41\x1f\x31\xde\x20\x42\x54\x9d\xc5\x30\x0c\xc3\x30\xcc\x35\x23\xb4\x5c\xc3\x7d\x9f\x61\x98\x5b\x4c\x3e\x3f\xf8\xe4\x80\xbc\xd4\x16\x74\xdc\x22\x3b\x85\x1a\x8f\xec\x93\x03\xf2\x52\x5b\xd0\x79\x16\xd9\x21\xbb\x64\x8f\xec\x93\x03\xf2\x52\x9b\x26\x2d\x41\x8b\x0f\x41\xef\x2c\x68\x85\x22\x3c\xb2\x4f\x0e\x2e\xa7\x37\x0c\x73\x08\xe4\x6b\xf7\xc7\x88\x31\xc1\x10\x69\x69\xfd\xde\x47\x84\x04\x63\x4c\x90\x61\xae\x8e\x67\xe8\x63\xa6\xce\x9a\x23\x89\xc6\xed\xde\x30\x89\x67\x2a\x45\x74\xd7\x4d\x73\xe2\x24\x0c\x07\x51\x32\xc8\x7f\x88\xb4\x7e\x7f\x3c\xbe\xbf\xfa\xf9\xac\x79\x96\xdc\x5b\xd6\xec\x6d\x71\x7c\x4a\x8f\x0c\x39\x09\x52\x48\x8c\x90\x1a\xfb\x63\x3f\xf8\x57\xb7\xbe\xd8\x17\x74\xff\x3e\x19\x7f\x53\x63\xa9\x99\x21\x77\xce\x28\xf5\xc1\x0e\xf4\x35\x53\x5d\x9f\x61\x8a\x85\xb1\x07\xd6\x87\x55\x9d\xda\xe2\xf8\xa7\x4f\xd7\x9f\x5e\xe9\xcf\x5e\xa3\x5e\xee\x54\x5f\x1a\xb7\xf5\x03\x6a\x93\xac\xba\x76\xae\x5e\xc9\x9f\x8f\xf0\xce\x38\x7e\x7c\xfd\xd2\x24\xa3\xd8\x87\xef\x2f\x9f\xff\xb9\xfb\x59\x8d\xa3\x41\x8e\xdc\x4b\x4e\xa9\x2f\x38\xda\xf6\xbd\x58\x20\x53\x57\x93\x79\xde\x70\xfe\xff\xaa\x5b\x5f\xec\xc7\xc9\x7b\x3b\x95\xdb\xaf\x69\x43\x86\xdc\x39\xa3\xd4\x87\x3b\xde\x95\x4d\xd8\x0c\xb3\x27\x6c\xad\x56\xbe\xfe\x7f\x61\xde\xff\x67\x18\xe6\x80\x11\x4e\xd8\x0b\x3b\xe7\x1b\x82\x1b\x58\xb4\x11\xf0\xfa\xac\xe0\x82\x8d\x00\x51\xf8\xc3\xf0\xc6\x6d\x04\xf0\xfd\x9f\x61\x70\x1a\x00\x00\xff\xff\x8f\x43\x39\x65\x04\x20\x00\x00")

func proto_ds_store() ([]byte, error) {
	return bindata_read(
		_proto_ds_store,
		"proto/.DS_Store",
	)
}

var _proto_micro_mall_users_proto_license = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x92\xc1\x6e\xdc\x3c\x0c\x84\xef\x06\xfc\x0e\x73\xcb\xff\x03\xdb\x45\x93\xde\xda\x53\xd1\x53\x80\x1e\x7a\xc8\x0b\xd0\x12\xbd\x66\x23\x8b\x2e\x25\xad\xbb\x7d\xfa\x82\xda\x34\x4d\x90\xde\x6c\x89\x1c\x0e\xbf\xd1\xaf\x24\x13\xbe\x4a\xe0\x5c\x78\x1c\xc6\xe1\xbf\x2f\xff\xe3\xee\xfd\xdd\x2d\x82\x49\xa9\x42\xb9\x9f\x3e\x2c\x52\x50\x74\xae\x3b\x19\x43\x0a\x36\xd3\xb3\x44\x8e\xb8\xa1\xf2\x4e\xca\xcd\x01\xbb\xd4\x45\x5b\x05\xe5\x0b\xf8\xe7\x66\x5c\x0a\xd4\x20\xeb\x96\x84\xe3\x38\xec\x64\x46\xb9\x5e\x8e\xc0\x7d\x46\x56\xf0\x99\x73\xc5\x2e\x29\xa1\x2e\x0c\x6a\x75\x51\x2b\x98\x18\x0b\xa7\x88\x24\x34\x25\xc6\xac\xd6\x25\x23\xad\x74\xe2\x32\x0e\x64\x52\x24\x9f\x30\x9b\xae\xbd\xb1\x15\x86\xce\xa8\x2f\x2d\x1e\xdd\xf4\x37\xb6\x55\x4a\x11\xcd\xee\xf8\xe4\xd3\x39\xa2\xaa\xeb\x69\x66\xff\xf2\xde\x57\x8d\xcf\xf3\xb6\x66\x9b\x16\x3e\x8c\x83\xe4\x90\x5a\xf4\x91\x41\xd7\x95\x2d\x08\x25\xd0\xb6\x25\x09\x54\x45\x73\x39\x80\xf2\x55\x37\x55\x36\x48\xed\xff\xc6\x51\x4a\x35\x99\x5a\x65\x48\x1d\x87\xd9\x98\xd3\xe5\x80\xd2\xa6\xef\x1c\xaa\x37\xb8\xfd\x59\x53\xd2\xdd\xe5\x8d\xbd\x3e\x74\xcd\x8f\xbe\xc0\xed\x11\x0f\x0b\x43\x4d\x4e\x92\xdf\xec\x88\xb5\x95\x8a\xac\xd5\x91\xad\x52\x8c\x1d\x3a\xfb\x92\x9f\x70\xd1\xf6\x7c\x3f\x0e\x00\x42\x22\x71\x5e\x54\xfb\xdd\x6e\x5a\xb9\xcf\xbf\xaa\x53\xfa\xcb\x0e\xf7\x73\xaf\x79\x03\xa7\x0b\x49\x06\x79\xfc\xb1\x85\xea\x9b\x83\xc2\x63\xd6\x3d\x71\x3c\xad\x1e\xa8\xe4\x2e\xfb\x54\x81\xa8\xa1\xf9\x79\x47\x85\x5d\x5b\x8a\x98\xae\x4a\xb4\x6d\xc6\x41\xc8\x53\x99\x5a\xf5\x90\x7c\x1b\xe3\x1f\x4d\x8c\xe3\x71\x1c\xee\x8e\xf8\xec\x4c\x39\xa2\x68\xb3\xc0\x38\xb3\x79\xa0\xe5\xba\xdc\xc4\xd8\x12\x49\x4e\x17\xac\x64\x8f\x1c\x41\x05\xa5\x85\xe5\x1a\xc9\x0b\x40\x7d\xe0\x6b\x48\x5e\x3b\xb1\x83\xff\x37\x86\x71\xf8\xe0\xfc\xaf\xae\x24\x30\x56\xba\xfc\xc1\x6d\xbc\xea\x99\xa3\x3f\x70\x7a\x32\xd8\xdf\xa3\xbf\x9c\x27\xa7\xcf\xe9\x8b\xe6\xe3\x38\xfc\x0e\x00\x00\xff\xff\x4b\x39\xd9\x60\x68\x03\x00\x00")

func proto_micro_mall_users_proto_license() ([]byte, error) {
	return bindata_read(
		_proto_micro_mall_users_proto_license,
		"proto/micro_mall_users_proto/LICENSE",
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
	"proto/.DS_Store": proto_ds_store,
	"proto/micro_mall_users_proto/LICENSE": proto_micro_mall_users_proto_license,
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
	"proto": &_bintree_t{nil, map[string]*_bintree_t{
		".DS_Store": &_bintree_t{proto_ds_store, map[string]*_bintree_t{
		}},
		"micro_mall_users_proto": &_bintree_t{nil, map[string]*_bintree_t{
			"LICENSE": &_bintree_t{proto_micro_mall_users_proto_license, map[string]*_bintree_t{
			}},
		}},
	}},
}}
