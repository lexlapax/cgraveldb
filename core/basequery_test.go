package core

import (
	"bytes"
	"testing"
	// "fmt"
	"github.com/lexlapax/graveldb/util"
	"github.com/stretchr/testify/assert"
)
func TestBaseQuery(t *testing.T) {
	// t.Skip()

	query := newBaseQuery()
	assert.Equal(t, 0, len(query.Conditions))
	assert.Equal(t, util.MaxInt, query.Max)
	query.Limit(10)
	assert.Equal(t, 10, query.Max)

	query.HasKey("testkey")
	assert.Equal(t, 1, len(query.Conditions))

	query = newBaseQuery()
	query.NoKey("nokey")
	assert.Equal(t, 1, len(query.Conditions))

	query = newBaseQuery()
	query.HasKeyValue("testkeyval", []byte("somevalue"))
	assert.Equal(t, 1, len(query.Conditions))

	query = newBaseQuery()
	query.NoKeyValue("nokeyval", []byte("somevalue"))
	assert.Equal(t, 1, len(query.Conditions))

	query = newBaseQuery()
	query.Limit(10).HasKey("testkey").NoKey("nokey").HasKeyValue("testkeyval", []byte("somevalue")).NoKeyValue("nokeyval", []byte("somevalue"))
	assert.Equal(t, 10, query.Max)
	assert.Equal(t, 4, len(query.Conditions))

	for _, cond := range query.Conditions {
		switch cond.Key {
		case "testkey":
			assert.True(t, cond.Value == nil)
			assert.True(t, cond.Has == true)
		case "nokey":
			assert.True(t, cond.Value == nil)
			assert.True(t, cond.Has == false)
		case "testkeyval":
			assert.True(t, bytes.Compare(cond.Value, []byte("somevalue")) == 0)
			assert.True(t, cond.Has == true)
		case "nokeyval":
			assert.True(t, bytes.Compare(cond.Value, []byte("somevalue")) == 0)
			assert.True(t, cond.Has == false)
		}
	}
}