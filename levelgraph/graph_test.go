package levelgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
	"reflect"
	//"github.com/jmhodges/levigo"
	//"fmt"
	//"github.com/lexlapax/graveldb/core"	
)

func TestOpenGraph(t *testing.T){
	dbdir := "./testing.db"
	gdb,_ := opengraph(dbdir)
	if assert.NotNil(t, &gdb) {
		assert.Equal(t, "./testing.db", gdb.dbdir)
		assert.Equal(t, "*levigo.DB", reflect.TypeOf(gdb.meta).String())
		assert.NotNil(t, gdb.meta)
		assert.NotNil(t, gdb.elements)
		assert.NotNil(t, gdb.hs)
		assert.NotNil(t, gdb.props)
		fi, _ := os.Lstat(dbdir)
		assert.True(t, fi.IsDir(), "dbdir should be a directory")
		assert.Equal(t, "testing.db", fi.Name(), "dbdir name should match" )
		assert.Equal(t, "dbdir=./testing.db", gdb.String())
	}
}

func TestAddVertex(t *testing.T) {
	
}