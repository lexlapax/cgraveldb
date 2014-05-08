package levigo

import (
		"sync"
		"strings"
		"regexp"
		"github.com/jmhodges/levigo"
		"github.com/lexlapax/graveldb/util"
		// "fmt"
)

type InvertedIndex struct {
	wordtodoc *levigo.DB
	doctoword *levigo.DB
	dbdir string
	opts *levigo.Options
	ro *levigo.ReadOptions
	wo *levigo.WriteOptions
	recsep []byte
	isopen bool
	// rwlock *sync.RWMutex

	// wordtodoc map[string][]string
	// doctoword map[string][]string
	sync.RWMutex
}

func NewInvertedIndex() *InvertedIndex {
	idx := new(InvertedIndex)
	// idx.wordtodoc = make(map[string][]string)
	// idx.doctoword = make(map[string][]string)
	return idx
}

func (index *InvertedIndex) Clear() {
	index.Lock()
	defer index.Unlock()
	// index.wordtodoc = make(map[string][]string)
	// index.doctoword = make(map[string][]string)
}

func (index *InvertedIndex) DocCount() int {
	index.RLock()
	defer index.RUnlock()

	// return len(index.doctoword)
	return 0
}

func (index *InvertedIndex) Docs() []string {
	index.RLock()
	defer index.RUnlock()
	keys := []string{}
	// for k,_ := range index.doctoword {
	// 	keys = append(keys, k)
	// }
	return keys
}

func (index *InvertedIndex) Tokens() []string {
	index.RLock()
	defer index.RUnlock()
	keys := []string{}
	// for k,_ := range index.wordtodoc {
	// 	keys = append(keys, k)
	// }
	return keys
}

func (index *InvertedIndex) TokenCount() int {
	index.RLock()
	defer index.RUnlock()

	// return len(index.wordtodoc)
	return 0
}

func (index *InvertedIndex) AddDoc(id string, doc string) {
	if id == "" || doc == "" { return }
	index.Lock()
	defer index.Unlock()
	words := []string{}
	re := regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")
	for _, word := range strings.Fields(doc) {
		words = append(words,  re.ReplaceAllString(word, ""))
	}


	// var ids []string
	// for _, word := range words {
	// 	if curids, ok := index.wordtodoc[word]; ok {
	// 		idset := core.NewStringSet()
	// 		idset.AddArray(curids)
	// 		idset.Add(id)
	// 		ids = idset.Members()
	// 	} else {
	// 		ids = []string{id}
	// 	}
	// 	//fmt.Printf("ids=%v\n", ids)
	// 	index.wordtodoc[word] = ids
	// }
	// index.doctoword[id] = words
	// fmt.Printf("w2d=%v\n", index.wordtodoc)
	// fmt.Printf("d2w=%v\n", index.doctoword)
}

func (index *InvertedIndex) Search(keywords ...string) []string {
	idset := util.NewStringSet()
	if len(keywords) < 1 { return idset.Members() }
	index.RLock()
	defer index.RUnlock()
	// for _, keyword := range keywords {
	// 	if keyword == "" { continue }
	// 	if val, ok := index.wordtodoc[keyword]; ok {
	// 		idset.AddArray(val)
	// 	}
	// }
	return idset.Members()
}

func (index *InvertedIndex) DelDoc(id string) {
	if id == "" { return }
	index.Lock()
	defer index.Unlock()
	// if words, ok := index.doctoword[id]; ok {
	// 	for _, word := range words {
	// 		if ids, ok := index.wordtodoc[word]; ok {
	// 			idset := core.NewStringSet()
	// 			idset.AddArray(ids)
	// 			idset.Del(id)
	// 			if idset.Count() > 0 {
	// 				index.wordtodoc[word] = idset.Members()
	// 			} else {
	// 				delete(index.wordtodoc, word)
	// 			}
	// 		}
	// 	}
	// 	delete(index.doctoword, id)
	// }
}

