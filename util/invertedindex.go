package util

import (
		"sync"
		// "strings"
		// "regexp"
		"errors"
		// "fmt"
)

var (
	ErrNilValue = errors.New("value passed cannot be nil")
	ErrInvalidParameter = errors.New("the value of parameter passed is invalid")
	ErrDBNotOpen = errors.New("db is not open")
)

type InvertedIndex struct {
	indexc *IndexContainer
	wordtodoc map[string][]string
	doctoword map[string][]string
	sync.RWMutex
}

type IndexContainer struct {
	indices map[string]*InvertedIndex
	sync.RWMutex
	isopen bool
}

func (indexc *IndexContainer) Open(args ...interface{}) error {
	if indexc.isopen == true {return nil }//errors.New("db already open") }
	indexc.indices = make(map[string]*InvertedIndex)
	return nil
}

func (indexc *IndexContainer) AddIndex(indexkey string) (*InvertedIndex, error) {
	if !indexc.isopen {return nil, ErrDBNotOpen }
	if indexkey == "" {return nil, ErrNilValue}
	if _,ok := indexc.indices[indexkey]; ok {
		return nil, nil
	}
	index := NewInvertedIndex()
	index.indexc = indexc
	indexc.indices[indexkey] = index
	return index, nil
}

func (indexc *IndexContainer) ClearIndex(indexkey string) error {
	if !indexc.isopen {return ErrDBNotOpen }
	delete(indexc.indices, indexkey)
	return nil
}

func NewInvertedIndex() *InvertedIndex {
	idx := new(InvertedIndex)
	idx.wordtodoc = make(map[string][]string)
	idx.doctoword = make(map[string][]string)
	return idx
}

func (index *InvertedIndex) Clear() {
	index.Lock()
	defer index.Unlock()
	index.wordtodoc = make(map[string][]string)
	index.doctoword = make(map[string][]string)
}

func (index *InvertedIndex) DocCount() int {
	index.RLock()
	defer index.RUnlock()

	return len(index.doctoword)
}

func (index *InvertedIndex) Docs() []string {
	index.RLock()
	defer index.RUnlock()
	keys := []string{}
	for k,_ := range index.doctoword {
		keys = append(keys, k)
	}
	return keys
}

func (index *InvertedIndex) Tokens() []string {
	index.RLock()
	defer index.RUnlock()
	keys := []string{}
	for k,_ := range index.wordtodoc {
		keys = append(keys, k)
	}
	return keys
}

func (index *InvertedIndex) TokenCount() int {
	index.RLock()
	defer index.RUnlock()

	return len(index.wordtodoc)
}

func (index *InvertedIndex) AddDoc(id string, doc string) {
	if id == "" || doc == "" { return }
	index.Lock()
	defer index.Unlock()
	// words := []string{}
	// re := regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")
	// for _, word := range strings.Fields(doc) {
	// 	words = append(words,  re.ReplaceAllString(word, ""))
	// }
	words := WhiteSpaceTokenize(doc)

	var ids []string
	for _, word := range words {
		if curids, ok := index.wordtodoc[word]; ok {
			idset := NewStringSet()
			idset.AddArray(curids)
			idset.Add(id)
			ids = idset.Members()
		} else {
			ids = []string{id}
		}
		//fmt.Printf("ids=%v\n", ids)
		index.wordtodoc[word] = ids
	}
	index.doctoword[id] = words
	// fmt.Printf("w2d=%v\n", index.wordtodoc)
	// fmt.Printf("d2w=%v\n", index.doctoword)
}

func (index *InvertedIndex) Search(keywords ...string) []string {
	idset := NewStringSet()
	if len(keywords) < 1 { return idset.Members() }
	index.RLock()
	defer index.RUnlock()
	searchwords := NewStringSet()

	for _, keyword := range keywords {
		if keyword == "" { continue }
		searchwords.AddArray(WhiteSpaceTokenize(keyword))
	}
	for _, word := range searchwords.Members() {
		if val, ok := index.wordtodoc[word]; ok {
			idset.AddArray(val)
		}
	}
	return idset.Members()
}

func (index *InvertedIndex) DelDoc(id string) {
	if id == "" { return }
	index.Lock()
	defer index.Unlock()
	if words, ok := index.doctoword[id]; ok {
		for _, word := range words {
			if ids, ok := index.wordtodoc[word]; ok {
				idset := NewStringSet()
				idset.AddArray(ids)
				idset.Del(id)
				if idset.Count() > 0 {
					index.wordtodoc[word] = idset.Members()
				} else {
					delete(index.wordtodoc, word)
				}
			}
		}
		delete(index.doctoword, id)
	}
}

