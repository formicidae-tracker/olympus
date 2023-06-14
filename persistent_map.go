package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PersistentMap[T any] struct {
	Map map[string]T

	opaqueKeys map[string]string
	path       string
}

func NewPersistentMap[T any](name string) *PersistentMap[T] {
	res := &PersistentMap[T]{
		Map:        make(map[string]T),
		opaqueKeys: make(map[string]string),
		path:       filepath.Join(_datapath, name),
	}
	res.restore()
	return res
}

type persistentEntry[T any] struct {
	Key   string
	Value T
}

var subdirRx = regexp.MustCompile(`\A[[:xdigit:]]{2}\z`)
var fileRx = regexp.MustCompile(fmt.Sprintf(`\A[[:xdigit:]]{%d}\.json\z`, 2*md5.Size))

func (m *PersistentMap[T]) restore() {
	subdirs, err := os.ReadDir(m.path)
	if err != nil {
		return
	}
	for _, subdir := range subdirs {
		if subdir.IsDir() && subdirRx.MatchString(subdir.Name()) {
			m.restoreSubdir(subdir.Name())
		}
	}
}

func (m *PersistentMap[T]) restoreSubdir(name string) {
	dirpath := filepath.Join(m.path, name)
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return
	}

	for _, entry := range files {
		if entry.IsDir() == false && fileRx.MatchString(entry.Name()) {
			m.restoreFile(name, entry.Name())
		}
	}
}

func (m *PersistentMap[T]) restoreFile(subdir, filename string) {
	filename = filepath.Join(m.path, subdir, filename)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	data := persistentEntry[T]{}
	if err := dec.Decode(&data); err != nil {
		return
	}
	opKey := strings.TrimSuffix(filepath.Base(filename), ".json")
	keyBytes, err := hex.DecodeString(data.Key)
	if err != nil {
		return
	}
	key := string(keyBytes)
	// restores the state of opaqueKeys to resolves hash collisions
	m.opaqueKeys[opKey] = key
	m.Map[key] = data.Value
}

func (m *PersistentMap[T]) Save() error {
	for key := range m.Map {
		if err := m.SaveKey(key); err != nil {
			return fmt.Errorf("saving key %s: %w", key, err)
		}
	}
	return nil
}

func (m *PersistentMap[T]) SaveKey(key string) (err error) {
	value, ok := m.Map[key]
	if ok == false {
		return nil
	}
	opKey := m.opaqueKey(key)
	defer func() {
		// if we save the key, we need to save the chosen opaqueKey
		if err != nil {
			return
		}
		m.opaqueKeys[opKey] = key
	}()
	subdir := opKey[:2]
	filename := filepath.Join(m.path, subdir, opKey+".json")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	data := persistentEntry[T]{
		Key:   hex.EncodeToString([]byte(key)),
		Value: value,
	}
	return enc.Encode(data)
}

func (m *PersistentMap[T]) opaqueKey(key string) string {
	// md5 collisions are very rare, but they do exists. Here it is
	// how we would resolve them. We store each saved opaqueKey in
	// memory, and if we found a collision, we try a modified hash by
	// appending a char to the string (potentially multiple
	// times). When restoring state, we must also restore the listing
	// of opaqueKeys. If we do not go for around 2^128 elements (it
	// would be stupid to do so), the loop will end. See unit test for
	// the difficulty to actually makes md5 collide (it is not null,
	// but still hard when using common string).
	actualKey := key
	for {
		opKey := fmt.Sprintf("%x", md5.Sum([]byte(actualKey)))
		if oldKey, ok := m.opaqueKeys[opKey]; ok == true && oldKey != key {
			// we got an hash collision, find a new hash.
			actualKey += "Z"
		} else {
			// no hash collision, either never seen or we are re-saving the key
			return opKey
		}
	}
}

var _datapath string

func init() {
	_datapath = os.Getenv("OLYMPUS_DATA_HOME")
	if len(_datapath) == 0 {
		_datapath = filepath.Join(os.TempDir(), "fort", "olympus")
	}
}
