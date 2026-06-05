package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

//go:embed zoneinfo.zip
var tzData []byte

var tzZip = func() *zip.Reader {
	r, err := zip.NewReader(bytes.NewReader(tzData), int64(len(tzData)))
	if err != nil {
		panic(err)
	}
	return r
}()

func zoneNames() []string {
	var names []string
	for _, f := range tzZip.File {
		if f.FileInfo().IsDir() {
			continue
		}
		n := f.Name
		if n == "Factory" || n == "posixrules" || strings.HasPrefix(n, "SystemV/") {
			continue
		}
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func loadZone(name string) (*time.Location, error) {
	for _, f := range tzZip.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		b, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		return time.LoadLocationFromTZData(name, b)
	}
	return nil, fmt.Errorf("zone not found: %s", name)
}
