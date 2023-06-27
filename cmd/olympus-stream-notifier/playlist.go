package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ReadPlaylist(filename string) (map[string]bool, error) {
	res := make(map[string]bool)
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return res, nil
		}
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		l, err := reader.ReadString('\n')
		if err == io.EOF {
			return res, nil
		}
		if err != nil {
			return nil, err
		}
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == '#' {
			continue
		}
		hostname := strings.TrimSuffix(path.Base(l), ".m3u8")
		if len(hostname) == 0 {
			continue
		}
		res[hostname] = true
	}

}

func PlaylistFilename(basepath string) string {
	return filepath.Join(basepath, "master.m3u8")
}

func WritePlaylist(filename string, hosts map[string]bool) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, "#EXTM3U")
	if err != nil {
		return err
	}
	for host, _ := range hosts {
		_, err = fmt.Fprintf(f, "/hls/%s.m3u8\n", host)
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterPlaylist(basepath string, hostname string) error {
	hosts, err := ReadPlaylist(PlaylistFilename(basepath))
	if err != nil {
		return err
	}
	hosts[hostname] = true
	return WritePlaylist(PlaylistFilename(basepath), hosts)
}

func UnregisterPlaylist(basepath string, hostname string) error {
	hosts, err := ReadPlaylist(PlaylistFilename(basepath))
	if err != nil {
		return err
	}
	delete(hosts, hostname)
	return WritePlaylist(PlaylistFilename(basepath), hosts)
}
