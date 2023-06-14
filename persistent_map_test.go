package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type PersistentMapSuite struct {
	datapath string
	dir      string

	m *PersistentMap[string]

	collidingKeys []string
}

var _ = Suite(&PersistentMapSuite{})

func (s *PersistentMapSuite) SetUpSuite(c *C) {
	s.datapath = _datapath

	//md5 collision found here: https://web.archive.org/web/20220812190633/https://www.links.org/?p=6
	hexKeys := []string{
		"d131dd02c5e6eec4693d9a0698aff95c2fcab58712467eab4004583eb8fb7f8955ad340609f4b30283e488832571415a085125e8f7cdc99fd91dbdf280373c5bd8823e3156348f5bae6dacd436c919c6dd53e2b487da03fd02396306d248cda0e99f33420f577ee8ce54b67080a80d1ec69821bcb6a8839396f9652b6ff72a70",
		"d131dd02c5e6eec4693d9a0698aff95c2fcab50712467eab4004583eb8fb7f8955ad340609f4b30283e4888325f1415a085125e8f7cdc99fd91dbd7280373c5bd8823e3156348f5bae6dacd436c919c6dd53e23487da03fd02396306d248cda0e99f33420f577ee8ce54b67080280d1ec69821bcb6a8839396f965ab6ff72a70",
	}

	for _, hexKey := range hexKeys {
		key, err := hex.DecodeString(hexKey)
		c.Assert(err, IsNil)
		s.collidingKeys = append(s.collidingKeys, string(key))
	}

	c.Check(s.collidingKeys[0], Not(Equals), s.collidingKeys[1])
	hasher := func(key string) string {
		return fmt.Sprintf("%s", md5.Sum([]byte(key)))
	}
	c.Check(hasher(s.collidingKeys[0]), Equals, hasher(s.collidingKeys[1]))
}

func (s *PersistentMapSuite) TearDownSuite(c *C) {
	_datapath = s.datapath
}

func (s *PersistentMapSuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	_datapath = s.dir
	s.m = NewPersistentMap[string]("unit-test")
}

func (s *PersistentMapSuite) TestIsEmptyUponCreation(c *C) {
	c.Check(s.m.Map, HasLen, 0)
}

func (s *PersistentMapSuite) TestValuePersistence(c *C) {
	s.m.Map["foo"] = "something"
	c.Assert(s.m.SaveKey("foo"), IsNil)
	opKey := s.m.opaqueKey("foo")
	_, err := os.Stat(filepath.Join(s.dir, "unit-test", opKey[:2], opKey+".json"))
	c.Check(err, IsNil)

	newMap := NewPersistentMap[string]("unit-test")
	value, ok := newMap.Map["foo"]
	c.Check(value, Equals, "something")
	c.Check(ok, Equals, true)

}

func (s *PersistentMapSuite) TestHashCollisionResolution(c *C) {
	s.m.Map[s.collidingKeys[0]] = "value0"
	s.m.SaveKey(s.collidingKeys[0])
	s.m.Map[s.collidingKeys[1]] = "value1"
	s.m.SaveKey(s.collidingKeys[1])

	newMap := NewPersistentMap[string]("unit-test")
	c.Assert(newMap.Map, HasLen, 2)
	c.Check(newMap.Map[s.collidingKeys[0]], Equals, "value0")
	c.Check(newMap.Map[s.collidingKeys[1]], Equals, "value1")
	c.Check(newMap.opaqueKeys, HasLen, len(s.m.opaqueKeys))
	for opKey, key := range s.m.opaqueKeys {
		c.Check(newMap.opaqueKeys[opKey], Equals, key)
	}

}
