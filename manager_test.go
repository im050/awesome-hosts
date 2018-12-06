package main

import (
	"awesome-hosts/manager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPattern(t *testing.T) {
	m := manager.New("")
	assert.Error(t, m.CheckIP("123"))
	assert.NoError(t, m.CheckIP("127.0.0.1"))
	assert.NoError(t, m.CheckIP("fe80::181c:65cd:d4ff:9885"))
	assert.NoError(t, m.CheckIP("::1"))
	assert.NoError(t, m.CheckDomain("localhost"))
	assert.NoError(t, m.CheckDomain("www.baidu.com"))
	assert.Error(t, m.CheckDomain("good,|>asdf.asdf"))
	assert.Error(t, m.CheckDomain(".host"))
	assert.Error(t, m.CheckDomain("host."))
	assert.Error(t, m.CheckGroupName(",hel?lo"))
	assert.Error(t, m.CheckGroupName("h<e>ll|\\o"))
	assert.NoError(t, m.CheckGroupName("he.llo"))
	assert.NoError(t, m.CheckGroupName("he-_llo"))
	assert.NoError(t, m.CheckGroupName("Good boy"))
}

func TestRemoveRepeat(t *testing.T) {
	iSlc := []int{1,2,3,4,4,5,6,7,8,8,9}
	rSlc := manager.RemoveRepeatNumber(iSlc)
	expected := []int{1,2,3,4,5,6,7,8,9}
	assert.Equal(t, expected, rSlc)
}

func BenchmarkCheckIP(b *testing.B) {
	m := manager.New("")
	for i := 0; i < b.N; i++ {
		m.CheckIP("::1")
	}
}

func BenchmarkCheckDomain(b *testing.B) {
	m := manager.New("")
	for i := 0; i < b.N; i++ {
		m.CheckDomain("::1")
	}
}

func BenchmarkCheckGroupName(b *testing.B) {
	m := manager.New("")
	for i := 0; i < b.N; i++ {
		m.CheckGroupName("aaaaaaaaaaaaaa")
	}
}

//func BenchmarkSync(b *testing.B) {
//	m := manager.New(manager.GetUserHome() + "/.awesohosts")
//	m.SudoPassword = "your root password here"
//	m.SyncSystemHostsUnix()
//}
