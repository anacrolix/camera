package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func groupWithoutExt(ss []string) (ret map[string][]string) {
	ret = make(map[string][]string, len(ss))
	for _, s := range ss {
		ext := filepath.Ext(s)
		gKey := s[:len(s)-len(ext)]
		gExts := ret[gKey]
		gExts = append(gExts, ext)
		ret[gKey] = gExts
	}
	return
}

func filterStringSlice(ss []string, f func(s string) bool) (ret []string) {
	ret = make([]string, 0, len(ss))
	for _, s := range ss {
		if f(s) {
			ret = append(ret, s)
		}
	}
	return
}

func hashGroup(key string, exts []string) ([]byte, error) {
	h := md5.New()
	for _, ext := range exts {
		f, err := os.Open(key + ext)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(h, f)
		f.Close()
		if err != nil {
			return nil, err
		}
	}
	return h.Sum(nil), nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	names = filterStringSlice(names, func(s string) bool { return !strings.HasPrefix(s, ".") })
	groups := groupWithoutExt(names)
	log.Println(groups)
	for key, exts := range groups {
		sum, err := hashGroup(key, exts)
		if err != nil {
			log.Printf("error hashing group %q: %s", key, err)
			continue
		}
		for _, ext := range exts {
			oldpath := key + ext
			newpath := filepath.Join(filepath.Dir(key), hex.EncodeToString(sum)+ext)
			log.Println("rename", oldpath, "->", newpath)
			err := os.Rename(oldpath, newpath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
