package mem

import (
	"testing"
	//"fmt"
	// "github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/stretchr/testify/assert"
)
func TestIndexInvertedIndex(t *testing.T) {
	// t.Skip()
	idx := NewInvertedIndex()
	assert.True(t, idx != nil)
	testdocs := make(map[string]string)
	stringset := core.NewStringSet()
	// testdocs["doc1"] = 
	idx.AddDoc("doc1", "my name is inverted index")
	assert.Equal(t, 1, idx.DocCount())
	assert.Equal(t, 5, idx.TokenCount())
	tokens := idx.Tokens()
	assert.Equal(t, 5, len(tokens))
	stringset.Clear()
	stringset.AddArray(tokens)
	assert.Equal(t, 5, stringset.Count())
	assert.True(t, stringset.Contains("my"))
	assert.True(t, stringset.Contains("name"))
	assert.True(t, stringset.Contains("is"))
	assert.True(t, stringset.Contains("inverted"))
	assert.True(t, stringset.Contains("index"))

	idx.AddDoc("doc2", "i index documents")
	assert.Equal(t, 2, idx.DocCount())
	assert.Equal(t, 7, idx.TokenCount())
	idx.AddDoc("doc3", "i am able to find documents with words")
	assert.Equal(t, 3, idx.DocCount())
	assert.Equal(t, 13, idx.TokenCount())
	tokens = idx.Tokens()
	assert.Equal(t, 13, len(tokens))
	stringset.Clear()
	stringset.AddArray(tokens)
	assert.Equal(t, 13, stringset.Count())

	testdocs["doc4"] = "i am case sensitive"
	testdocs["doc5"] = "names are not important, but a name is required"
	for k, v := range testdocs {
		idx.AddDoc(k, v)
	}
	assert.Equal(t, 5, idx.DocCount())
	stringset.Clear()
	stringset.AddArray(idx.Tokens())
	assert.Equal(t, 22, stringset.Count())
	assert.True(t, stringset.Contains("sensitive"))
	assert.True(t, stringset.Contains("important"))


	stringset.Clear()
	stringset.AddArray(idx.Search("somethingnotthere"))


	//single word searches
	stringset.Clear()
	stringset.AddArray(idx.Search("inverted"))
	assert.Equal(t, 1, stringset.Count())
	assert.True(t, stringset.Contains("doc1"))

	stringset.Clear()
	assert.Equal(t, 0, stringset.Count())
	stringset.AddArray(idx.Search("name"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc5"))

	stringset.Clear()
	stringset.AddArray(idx.Search("index"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc2"))

	stringset.Clear()
	stringset.AddArray(idx.Search("i"))
	assert.Equal(t, 3, stringset.Count())
	assert.True(t, stringset.Contains("doc2"))
	assert.True(t, stringset.Contains("doc3"))
	assert.True(t, stringset.Contains("doc4"))

	stringset.Clear()
	stringset.AddArray(idx.Search("am"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc3"))
	assert.True(t, stringset.Contains("doc4"))

	//multi word searches
	stringset.Clear()
	docids := idx.Search("am", "name")
	assert.Equal(t, 4, len(docids))
	stringset.AddArray(docids)
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc5"))
	assert.True(t, stringset.Contains("doc3"))
	assert.True(t, stringset.Contains("doc4"))

	stringset.Clear()
	docids = idx.Search("somethingnotthere", "am", "name")
	assert.Equal(t, 4, len(docids))
	stringset.AddArray(docids)
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc5"))
	assert.True(t, stringset.Contains("doc3"))
	assert.True(t, stringset.Contains("doc4"))

	stringset.Clear()
	docids = idx.Search("", "index", "name")
	assert.Equal(t, 3, len(docids))
	stringset.AddArray(docids)
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc2"))
	assert.True(t, stringset.Contains("doc5"))

	stringset.Clear()
	docids = idx.Search("name", "is", "i")
	assert.Equal(t, 5, len(docids))
	//fmt.Printf("v=%v\n", docids)
	stringset.AddArray(docids)
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc2"))
	assert.True(t, stringset.Contains("doc3"))
	assert.True(t, stringset.Contains("doc4"))
	assert.True(t, stringset.Contains("doc5"))

	//document deletions
	idx.DelDoc("doc6")
	assert.Equal(t, 5, idx.DocCount())
	stringset.Clear()
	stringset.AddArray(idx.Tokens())
	assert.Equal(t, 22, stringset.Count())

	idx.DelDoc("doc4")
	assert.Equal(t, 4, idx.DocCount())
	stringset.Clear()
	stringset.AddArray(idx.Search("i"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc2"))
	assert.True(t, stringset.Contains("doc3"))
	stringset.Clear()
	stringset.AddArray(idx.Search("index"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc2"))
	stringset.Clear()
	assert.Equal(t, 0, stringset.Count())
	stringset.AddArray(idx.Search("name"))
	assert.Equal(t, 2, stringset.Count())
	assert.True(t, stringset.Contains("doc1"))
	assert.True(t, stringset.Contains("doc5"))
	stringset.Clear()
	stringset.AddArray(idx.Tokens())
	assert.Equal(t, 20, stringset.Count())
	assert.False(t, stringset.Contains("sensitive"))
	assert.False(t, stringset.Contains("case"))


	idx.DelDoc("doc1")
	assert.Equal(t, 3, idx.DocCount())
	stringset.Clear()
	stringset.AddArray(idx.Search("name"))
	assert.Equal(t, 1, stringset.Count())
	assert.True(t, stringset.Contains("doc5"))

	stringset.Clear()
	stringset.AddArray(idx.Search("index"))
	assert.Equal(t, 1, stringset.Count())
	assert.True(t, stringset.Contains("doc2"))

	stringset.Clear()
	stringset.AddArray(idx.Tokens())
	assert.Equal(t, 18, stringset.Count())
	assert.False(t, stringset.Contains("inverted"))
	assert.False(t, stringset.Contains("my"))

	//index clear
	idx.Clear()
	assert.Equal(t, 0, idx.DocCount())
	assert.Equal(t, 0, idx.TokenCount())
}
