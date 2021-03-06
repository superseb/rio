// Code generated by go-bindata.
// sources:
// stacks/coredns-stack.yml
// stacks/istio-stack.yml
// DO NOT EDIT!

package stacks

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _stacksCorednsStackYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xcd\x6e\xdb\x3c\x10\xbc\xfb\x29\xf6\x01\x2c\x4a\x82\x3e\x7f\x75\x79\x6c\xea\x43\x2e\x39\xb4\x45\xaf\x06\x43\xae\x2d\xc2\xfc\xc3\x2e\x95\xc6\x68\xfb\xee\x05\x4d\xcb\x50\xe3\x46\x17\x2d\x77\x86\xcb\xc5\xcc\xe8\x18\x0e\xf6\xc8\x72\x05\x50\xca\xf2\xbf\x54\x19\x43\x96\xf0\xab\xb9\x9c\x01\x84\xdc\x0c\xf0\xf3\x7a\x28\x1f\x12\x45\xe2\x45\x63\x44\xe5\xf2\xb8\x68\x9c\xa6\x67\xa4\x80\x19\x19\xb4\x9b\x38\x23\x09\x17\xb5\x72\x60\x43\xa3\x8c\x21\xa1\x28\x29\xb0\xe9\xff\x5a\x2c\x87\x97\x2f\x45\xc3\x60\x03\xa3\x9e\x08\xdf\x60\x53\xe2\x4c\xa8\xfc\x9b\xf6\x41\x39\x97\x47\x8a\xd3\x71\xfc\xf7\x23\x0b\xfe\xef\x45\x9d\x28\x7a\xcc\x23\x4e\x0c\xf2\x63\xbf\x19\xfe\x86\x5e\xcf\x20\xa0\xc5\xac\x5b\x42\x8e\xee\x45\x14\x9d\x16\x14\xad\xf4\x88\x30\x74\xab\x79\x30\x23\xbd\x58\x8d\x57\x4d\x09\x4d\xe0\x2a\x6b\x29\x60\x8c\x9c\xab\xc8\x2a\xed\x95\x31\x12\x9e\x76\xdf\xf6\x9f\x1e\x9f\x3e\xef\xbf\xee\xbe\x7c\x7f\x7c\xd8\xdd\x50\x43\x31\x49\x50\xce\x5d\x4d\xf1\x5e\x05\x53\x47\x35\xd0\xdc\xd6\x68\xea\x76\xd7\xa7\xda\x87\x48\x78\xb0\x0e\x67\x27\x67\x7b\x0b\xf1\x62\xf1\xfb\x6c\x7c\x4d\x91\x71\x26\xf7\x9d\xf8\x6f\x10\x9d\xe8\x3b\xb9\x19\xe4\x66\x68\x27\x93\xd6\x26\xf0\x3b\x70\x81\x9a\xac\xd3\x15\x2e\x42\xae\x3d\x66\xb2\xba\xde\x38\xba\xf8\xac\xdc\x3e\x21\x79\xcb\x6c\x63\xb8\x6d\xe5\x2c\xe7\xf5\x0f\x95\xf5\x08\x18\x4c\x8a\x36\x64\xbe\x87\x66\x59\xef\x91\x12\x95\xfb\x6e\x50\x1e\x39\xa9\xf9\x46\xcd\xa7\x1e\x51\x9f\x24\x8c\x39\x27\xd9\xb6\x97\x3c\x16\x43\xe4\xb6\xdb\x76\xed\x22\xc2\xd6\xab\x23\x4a\x38\x6d\x59\x1c\x35\x09\x1b\x67\xc1\x64\x2f\x7a\x51\x13\xe2\xd1\x47\x3a\x4b\xf8\xd0\xd5\x20\x12\x2a\xb3\x8f\xc1\x9d\x25\x64\x9a\xaa\xa2\xac\x95\x43\x09\xfd\xea\x4f\x00\x00\x00\xff\xff\xd5\x57\x8c\x0c\x68\x03\x00\x00")

func stacksCorednsStackYmlBytes() ([]byte, error) {
	return bindataRead(
		_stacksCorednsStackYml,
		"stacks/coredns-stack.yml",
	)
}

func stacksCorednsStackYml() (*asset, error) {
	bytes, err := stacksCorednsStackYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stacks/coredns-stack.yml", size: 872, mode: os.FileMode(436), modTime: time.Unix(1532625128, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _stacksIstioStackYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\xdf\x6f\xdb\x36\x10\x7e\xf7\x5f\x41\x68\x05\xba\x16\x91\xe5\x74\xcd\x5a\x08\xf0\x83\xe1\xb8\x68\x80\xd4\xf6\x62\x17\xc3\x9e\x0c\x9a\x3c\xdb\x44\x28\x52\x20\x4f\x4e\xbd\x1f\xff\xfb\x40\x89\x92\xa8\x44\x4e\x82\x61\x4f\x36\xf9\x7d\x77\xe4\xf1\xee\x3e\x1d\xd3\x6a\x27\xf6\x36\x1d\x10\x92\x81\x3d\xb8\x5f\x42\x98\x56\x08\x0a\x53\xf2\x77\x5c\xae\x09\xe1\xc2\xd2\xad\x84\xa5\x96\x82\x9d\xa6\x07\x60\xf7\x36\x25\x68\x0a\xf0\xb8\x50\x7b\x03\xd6\x4e\xb5\x42\xa3\xa5\x04\xf3\x4d\x73\x48\x49\xb4\xf8\xf2\x25\xf2\x14\x5a\xe0\xa1\xb2\x4f\xc9\x7c\x31\x9f\xf9\x6d\xc3\xed\x1d\xec\x0c\xd8\xc3\x35\x48\x7a\x4a\xc9\xe5\xc8\x7a\x48\x17\xb8\xd5\x85\xe2\x6b\x43\x77\x3b\xc1\xbc\xb1\x07\x09\xc9\xca\x23\x26\xb7\xb7\x8b\xdf\x37\x93\xf9\x1f\xf5\x4d\x61\x47\x0b\x89\xd3\x32\xae\x96\xcc\x85\x65\xfa\x08\xe6\x74\xee\xb0\x32\x6a\x05\x0c\xd7\x22\x03\x5d\x60\x4a\x7e\xe9\x62\x3b\xb1\x5f\x52\x3c\xa4\x24\x4a\x00\x59\x22\x2c\x0a\x9d\xe4\x46\xff\x38\x45\x0d\x6d\x2b\x14\x35\xa7\x9a\x56\x58\x93\x48\xcd\xa8\x4c\xb6\x42\x25\xa0\x8e\x3a\xa0\x5a\x30\x47\xc1\x60\x2a\x0b\x8b\x60\x52\x52\xfa\x8b\x4b\x7f\xed\xa5\x0d\x15\xea\xba\x30\x14\x85\x56\x29\xf9\x78\xd5\x5e\x28\xa7\x06\x14\xae\x0e\x05\x72\xfd\x10\x70\x2e\xb3\xe0\xd6\x42\x21\x18\x06\xb9\x83\xaa\x84\xdc\xcd\xae\x6f\xee\x66\xd3\x75\xeb\xc7\x1d\x38\xe1\x99\x50\x4b\x6d\x30\x25\x97\x57\xa3\xd1\x28\x0c\xdb\xe5\x73\x29\xa9\x82\x49\x7f\xfe\x82\xa7\x9d\x70\xee\x8a\xa0\x89\x45\x48\x8d\xc3\x37\x7f\xcd\x27\xdf\x66\xab\xe5\x64\x3a\xfb\x27\x75\xde\x3f\x0d\x06\x3e\xf6\xb2\xea\x02\x6e\x5d\x7c\x59\x46\x15\x4f\x5b\xbf\x83\x36\x03\xb6\xe2\xc4\x55\xb1\x06\x89\xa8\xe0\xc4\x6d\x97\x0c\x50\x47\x61\xb4\xca\x5c\x19\x7b\x93\xe5\xe2\x7a\xe3\xee\x32\x7e\xf3\xb3\x05\xb9\x4b\x14\xcd\xe0\xdd\x23\xac\xbc\x67\x48\xb0\x39\x65\x2d\xeb\xe6\x76\xb1\xde\xac\xbf\xde\x2d\xd6\xeb\xdb\xd9\xf8\xca\xbf\x54\x0d\x4c\x27\xd3\xaf\xb3\xcd\xea\xb7\xef\x93\xd5\xd7\xf1\x55\x09\xed\xa5\xde\x52\xb9\xc9\xc1\x64\xc2\x5a\xa1\x55\x13\xc0\xdb\xf7\x3e\xa4\x61\x19\xc0\x50\xe8\xe4\xfd\xdb\x16\x53\x80\x0f\xda\xdc\x0b\x75\x06\x77\xdd\x04\x0a\x05\x2b\xf3\x7e\x86\x93\x0b\xf8\x81\xa0\xca\x73\x87\xf7\x9f\xad\x23\xb0\xc2\xa2\xce\x0c\x58\x5d\x18\x06\x1c\x76\x42\x09\xe7\xc2\x06\x86\xad\x55\x82\x07\x61\x78\x4e\x0d\x9e\x6a\x93\xd7\x13\x87\x2d\xdc\x6f\xe3\x55\xe3\x9c\xcb\x06\x4e\x2c\x52\x2c\x1a\x16\x33\x40\x11\x2e\xf6\x80\x17\x52\x58\xbc\x78\xa0\xc8\x0e\x17\x45\xce\x29\x82\x7f\xd4\x8c\xe6\xd6\xb3\x41\xf1\x5c\x0b\x85\xf5\x3a\xd7\xbc\xfe\x5b\xd7\xa1\x5f\x36\xf9\x6e\x36\x34\x87\x96\xcb\x0c\x78\x27\x22\xa3\x7b\xf0\x65\x9e\x54\xa5\x3b\x1a\x7e\x1e\x56\xd5\xe0\x89\x29\x11\xdc\x25\x08\x4f\x55\x99\x32\x30\xde\xda\x0a\x0e\xf7\xc2\x29\x68\x2d\x9e\x6d\xeb\xb7\x82\x05\x3f\x72\x6d\xa1\x5d\xc7\x65\x6b\x7e\x4a\x0e\x88\x79\x77\xf3\x72\x94\xec\x4d\xce\xda\xae\xef\xdc\xce\x79\x3d\x7e\x08\xee\x17\xf6\x58\xe0\xa8\xab\x3c\x31\x89\xe3\xae\x42\x05\x50\xd0\xb2\x1d\x03\x84\x2c\x97\x14\xe1\x8b\x90\x10\x00\x8f\xe5\xb2\x52\xc2\x4d\x25\x0f\x27\x9a\xc9\x21\x66\xb9\xec\x78\xea\xd7\x9d\x80\xd2\x11\xa0\x27\xbd\xfe\x52\xbf\xbf\xbe\xe7\x2b\xe6\xcd\x7c\xb5\x9e\xcc\xa7\xb3\xcd\xcd\xb2\xe6\x89\xfc\x5d\x20\xe5\xcf\x64\xbc\x91\x38\x26\x90\x72\x90\x69\x58\x41\x91\x17\x2f\x0f\x95\x49\x8a\x06\x4f\x32\xe4\x9e\x84\xe6\x39\x28\x1e\x73\x65\xe3\xf2\x8e\xe3\xe6\xeb\xeb\x50\x57\x00\x71\xae\x0d\x8e\x3f\x8f\x7e\x1d\x75\xb7\x0f\xda\xa2\x33\x19\xfb\x63\x06\x6d\x82\xe5\x2e\xb6\x62\xaf\x80\xc7\x8c\x76\x1d\x7a\x6e\x6c\x51\x1b\xba\x87\xb8\x79\x97\x71\x15\x8d\x3d\x59\x84\xec\x05\x8d\x7b\x30\x02\xa1\xd3\x39\x4d\xcf\x51\xc6\x74\xa1\x1e\xef\x56\xcb\x1e\x4f\xae\xd7\x39\x48\x40\x20\x1c\x72\xa9\x4f\x2e\xd9\xf6\x29\xd8\xef\xfd\x29\xa1\x07\xe9\xf7\x5b\x0a\x4c\x25\x2d\x35\x31\x50\xa8\xff\x62\x63\x20\x97\x82\x51\x5b\x3f\xc9\xeb\xaa\x67\x4f\x11\x1e\xa8\xd7\x07\x49\xb7\x20\x1b\xf9\x88\x3c\x16\xa5\x24\x72\xe7\x18\x45\x65\xd4\x53\x65\x1d\x29\xa8\x08\x0a\x30\x25\xae\x3e\xca\x15\x57\x36\x25\x2c\x68\xf7\x47\x45\xd8\x4a\x44\x4c\x8c\x2e\x6a\x56\x4c\xe2\xa3\xff\x13\x7d\x88\x9a\x0a\xea\x9d\xb8\x6a\x9d\xbf\xb4\x6f\xc9\x4f\xcf\x31\xe2\xb8\x33\xfc\xd4\x76\x1f\xaf\x4a\xc3\x1e\x28\x8e\xfb\x47\xa2\xe6\xc4\x6c\xe4\x4c\x9f\x25\x95\xc2\x13\xcc\x80\x8d\x6d\x69\xda\x8b\xf5\xca\x64\xfc\x64\x9a\x73\xc4\x3f\x45\x7e\x2f\x94\x1f\x92\xea\x07\x6b\xdf\xcb\x7d\xe3\x2c\xff\xce\xf3\xb3\x8c\xee\xac\x56\xe3\xe5\xc4\x16\x0d\x5e\x21\x9d\x81\x6c\x76\x12\xd4\x3d\x30\x1c\xc8\xaa\x79\xcd\x6b\xec\xff\x37\x47\x3d\xa3\xa7\x31\xb9\x59\xad\x6f\x16\x9b\x6f\xb3\xf5\x64\x73\xfe\xa4\x17\x3f\xb2\xe7\x45\x29\x72\x2d\x5a\x8d\x0c\x41\xa3\xbe\x30\xca\x44\xcf\x1b\xbf\x4f\x8e\xc2\x60\x41\x65\x2d\x2f\x2f\xf2\x39\x58\x14\xaa\x2c\x3e\x53\xc8\x57\x18\xf8\x26\xb7\xd1\xe0\xdf\x00\x00\x00\xff\xff\xe9\x71\x98\x65\xb3\x0d\x00\x00")

func stacksIstioStackYmlBytes() ([]byte, error) {
	return bindataRead(
		_stacksIstioStackYml,
		"stacks/istio-stack.yml",
	)
}

func stacksIstioStackYml() (*asset, error) {
	bytes, err := stacksIstioStackYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "stacks/istio-stack.yml", size: 3507, mode: os.FileMode(436), modTime: time.Unix(1532623795, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"stacks/coredns-stack.yml": stacksCorednsStackYml,
	"stacks/istio-stack.yml": stacksIstioStackYml,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"stacks": &bintree{nil, map[string]*bintree{
		"coredns-stack.yml": &bintree{stacksCorednsStackYml, map[string]*bintree{}},
		"istio-stack.yml": &bintree{stacksIstioStackYml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

