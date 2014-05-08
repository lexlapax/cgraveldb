package levigo

import (
		"sync"
		"strings"
		"regexp"
		"errors"
		"path"
		"os"
		"bytes"
		"encoding/binary"
		"github.com/jmhodges/levigo"
		"github.com/lexlapax/graveldb/util"
		// "fmt"
)

var (
	ErrNilValue = errors.New("value passed cannot be nil")
	ErrNoDirectory     = errors.New("need to pass a valid path for a db directory")
	ErrInvalidParameter = errors.New("the value of parameter passed is invalid")
	ErrDBNotOpen = errors.New("db is not open")
	wordtodoc = "reverse.db"	
	doctoword = "forward.db"
	metadb = "meta.db"	
	recsep = "\x1f" //the ascii unit separator dec val 31
	proprevcount = "revcount"
	propforcount = "forcount"
)


type IndexContainer struct {
	ridb *levigo.DB //reverse indexc
	fidb *levigo.DB //forward indexc
	metadb *levigo.DB //	metadb = "meta.db"	

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

type InvertedIndex struct {
	indexc *IndexContainer
	indexkey string
}

func (index *InvertedIndex) DocCount() uint {
	return index.indexc.DocCount(index.indexkey)
}

func (index *InvertedIndex) Docs() []string {
	return index.indexc.Docs(index.indexkey)
}

func (index *InvertedIndex) Tokens() []string {
	return index.indexc.Tokens(index.indexkey)
}

func (index *InvertedIndex) TokenCount() uint {
	return index.indexc.TokenCount(index.indexkey)
}

func (index *InvertedIndex) AddDoc(id string, doc string) {
	index.indexc.AddDoc(index.indexkey, id, doc)	
}

func (index *InvertedIndex) Search(keywords ...string) []string {
	return index.indexc.Search(index.indexkey, keywords...)	
}

func (index *InvertedIndex) DelDoc(id string) {
	index.indexc.DelDoc(index.indexkey, id)
}

func (index *InvertedIndex) Delete() {
	index.indexc.DelIndex(index.indexkey)
}

func (index *InvertedIndex) Clear() {
	index.indexc.ClearIndex(index.indexkey)
}


func (indexc *IndexContainer) Open(args ...interface{}) error {
	if indexc.isopen == true {return nil }//errors.New("db already open") }
	if len(args) > 0 {
		if aString, found := args[0].(string); found {
			if aString == "" {
				return ErrNoDirectory
			} else {
				indexc.dbdir = aString
			}
		} else {
			return ErrInvalidParameter
		}
	}
	if indexc.dbdir == "" { return ErrNoDirectory }

	err := os.MkdirAll(indexc.dbdir, 0755)
	if err != nil { return err }

	//indexc.recsep = []byte(recsep)
	//indexc.rwlock = &sync.RWMutex{}
	indexc.Lock()
	defer indexc.Unlock()

	indexc.opts = levigo.NewOptions()
	cache := levigo.NewLRUCache(3<<30)
	indexc.opts.SetCache(cache)
	indexc.opts.SetCreateIfMissing(true)
	filter := levigo.NewBloomFilter(10)
	indexc.opts.SetFilterPolicy(filter)

	indexc.ro = levigo.NewReadOptions()
	indexc.wo = levigo.NewWriteOptions()
	indexc.recsep = []byte(recsep)

	indexc.metadb, err = levigo.Open(path.Join(indexc.dbdir, metadb), indexc.opts)
	if err != nil {return err}

	indexc.ridb, err = levigo.Open(path.Join(indexc.dbdir, wordtodoc), indexc.opts)
	if err != nil {return err}

	indexc.fidb, err = levigo.Open(path.Join(indexc.dbdir, doctoword), indexc.opts)
	if err != nil {return err}

	indexc.isopen = true
	return nil
}

func(indexc *IndexContainer) getDbProperty(prop string) ([]byte, error){
	if prop == "" {return nil, ErrNilValue}
	val, err := indexc.metadb.Get(indexc.ro, []byte(prop))
	if err != nil {return nil, err}
	return val, nil
}

func(indexc *IndexContainer) putDbProperty(prop string, val []byte) ([]byte, error){
	if prop == "" {return nil, ErrNilValue}
	key := []byte(prop)
	oldval, err := indexc.metadb.Get(indexc.ro, key)
	if err != nil {return nil, err}
	err2 := indexc.metadb.Put(indexc.wo, key, val)
	if err2 != nil {return nil, err}
	return oldval, nil
}

func (indexc *IndexContainer) keepcount(key string, upordown int) (uint) {
	var storedcount, returncount uint

	val, _ := indexc.getDbProperty(key)
	if val == nil {
		storedcount = 0
		upordown = 0
	} else {
		tempint,_ := binary.Uvarint(val)
		storedcount = uint(tempint)
	}
	switch upordown {
		case -1:
			returncount = storedcount - 1
		case 1:
			returncount = storedcount + 1
		default:
			returncount = storedcount
	}
	if returncount != storedcount || val == nil {
		bufsize := binary.Size(uint64(returncount))
		buf := make([]byte, bufsize)
		binary.PutUvarint(buf, uint64(returncount))
		indexc.putDbProperty(key, buf)
	}
	return returncount
}


func OpenIndexContainer(dbdir string) (*IndexContainer, error ) {
	if dbdir == "" { return nil, ErrNoDirectory }

	indexc := new(IndexContainer)
	indexc.dbdir = dbdir
	indexc.isopen = false

	err := indexc.Open()
	if err != nil {return nil, err }
	return indexc, nil
}

func (indexc *IndexContainer) AddIndex(indexkey string) (*InvertedIndex, error) {
	if !indexc.isopen {return nil, ErrDBNotOpen }
	if indexkey == "" {return nil, ErrNilValue}
	index := &InvertedIndex{indexc, indexkey}
	indexc.keepcount(indexkey + recsep + "propforcount", 0)
	indexc.keepcount(indexkey + recsep + "proprevcount", 0)
	return index, nil
}

func (indexc *IndexContainer) ClearIndex(indexkey string) error {
	if !indexc.isopen {return ErrDBNotOpen }
	if indexkey == "" {return ErrNilValue}
	val, err := indexc.getDbProperty(indexkey + recsep + "propforcount")
	if val == nil { return nil }
	if err != nil { return err }
	
	indexc.Lock()
	defer indexc.Unlock()
	
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := indexc.ridb.NewIterator(ro)

	prefix := []byte(indexkey + recsep + "r" + recsep )
	it.Seek(prefix)
	propkeys := [][]byte{}
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		propkeys = append(propkeys, it.Key())
	}
	wb := levigo.NewWriteBatch()
	for _, propkey := range propkeys {
		wb.Delete(propkey)
	}
	indexc.ridb.Write(indexc.wo, wb)

	wb.Close()
	it.Close()
	ro.Close()
	
	ro = levigo.NewReadOptions()
	ro.SetFillCache(false)
	it = indexc.fidb.NewIterator(ro)

	prefix = []byte(indexkey + recsep + "f" + recsep )
	it.Seek(prefix)
	propkeys = [][]byte{}
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		propkeys = append(propkeys, it.Key())
	}
	wb = levigo.NewWriteBatch()
	for _, propkey := range propkeys {
		wb.Delete(propkey)
	}
	indexc.fidb.Write(indexc.wo, wb)
	// if err != nil { return err }
	wb.Close()
	it.Close()
	ro.Close()

	bufsize := binary.Size(uint64(0))
	buf := make([]byte, bufsize)
	binary.PutUvarint(buf, uint64(0))
	indexc.putDbProperty(indexkey + recsep + "propforcount", buf)
	indexc.putDbProperty(indexkey + recsep + "proprevcount", buf)

	return nil
}

func (indexc *IndexContainer) DelIndex(indexkey string) error {
	err := indexc.ClearIndex(indexkey)
	if err != nil { return err }
	indexc.Lock()
	defer indexc.Unlock()
	indexc.metadb.Delete(indexc.wo, []byte(indexkey + recsep + "propforcount"))
	indexc.metadb.Delete(indexc.wo, []byte(indexkey + recsep + "proprevcount"))
	return nil
}

func (indexc *IndexContainer) Clear() error {
	dbdir := indexc.dbdir
	indexc.Close()
	os.RemoveAll(dbdir)
	return indexc.Open()
}

func (indexc *IndexContainer) Close() error {
	indexc.Lock()
	defer indexc.Unlock()
	indexc.isopen = false
	indexc.metadb.Close()
	indexc.ridb.Close()
	indexc.fidb.Close()
	indexc.opts.Close()
	indexc.ro.Close()
	indexc.wo.Close()
	indexc.ridb = nil
	indexc.fidb = nil
	indexc.opts = nil
	indexc.ro = nil
	indexc.wo = nil
	return nil
}


func (indexc *IndexContainer) DocCount(indexkey string) uint {
	return indexc.keepcount(indexkey + recsep + "propforcount", 0)
}

func (indexc *IndexContainer) TokenCount(indexkey string) uint {
	return indexc.keepcount(indexkey + recsep + "proprevcount", 0)
}

func (indexc *IndexContainer) Docs(indexkey string) []string {
	indexc.RLock()
	defer indexc.RUnlock()
	keys := []string{}
	// for k,_ := range indexc.doctoword {
	// 	keys = append(keys, k)
	// }
	return keys
}

func (indexc *IndexContainer) Tokens(indexkey string) []string {
	indexc.RLock()
	defer indexc.RUnlock()
	tokens := []string{}

	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := indexc.ridb.NewIterator(ro)
	defer it.Close()
	defer ro.Close()
	prefix := []byte(indexkey + recsep + "r" + recsep )
	it.Seek(prefix)
	propkeys := [][]byte{}
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		propkeys = append(propkeys, it.Key())
	}
	for _, keyrecord := range propkeys {
		keys := byteArrayToStringArray(recsep, keyrecord)
		tokens = append(tokens, keys[len(keys) - 1])
	}

	return tokens
}

func stringArrayToByteArray(sep string, record []string) []byte {
	if record == nil { return nil }
	if len(record) < 1 { return nil } 
	bytearr := [][]byte{}
	for _, s := range record {
		bytearr = append(bytearr, []byte(s))
	}
	return bytes.Join(bytearr, []byte(sep))
}

func byteArrayToStringArray(sep string, record []byte) []string {
	strings := []string{}
	if record == nil {return strings}
	bytearray := bytes.Split(record, []byte(sep))
	for _, arr := range bytearray {
		strings = append(strings, string(arr[:]))
	}
	return strings
}

func (indexc *IndexContainer) AddDoc(indexkey string, id string, doc string) {
	if indexkey == "" || id == "" || doc == "" { return }

	indexc.Lock()
	defer indexc.Unlock()
	revprefix := indexkey + recsep + "r"
	forprefix := indexkey + recsep + "f"
	forkey := []byte(forprefix + recsep + id)
	keywords, _ := indexc.fidb.Get(indexc.ro, forkey)
	if keywords != nil { return } // already exists

	words := []string{}
	re := regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")
	for _, word := range strings.Fields(doc) {
		words = append(words,  re.ReplaceAllString(word, ""))
	}

	for _, word := range words {
		revkey := []byte(revprefix + recsep + word)
		ids := []string{}
		curids, _ := indexc.ridb.Get(indexc.ro, revkey)
		if curids != nil {
			ids = byteArrayToStringArray(recsep, curids)
			idset := util.NewStringSet()
			idset.AddArray(ids)
			idset.Add(id)
			ids = idset.Members()
		} else {
			ids = []string{id}
			indexc.keepcount(indexkey + recsep + "proprevcount", 1)
		}
		//fmt.Printf("ids=%v\n", ids)
		indexc.ridb.Put(indexc.wo, revkey, stringArrayToByteArray(recsep, ids))
		//increase word count
	}
	indexc.fidb.Put(indexc.wo, forkey, stringArrayToByteArray(recsep, words))
	indexc.keepcount(indexkey + recsep + "propforcount", 1)
	// // fmt.Printf("w2d=%v\n", indexc.wordtodoc)
	// // fmt.Printf("d2w=%v\n", indexc.doctoword)
}

func (indexc *IndexContainer) Search(indexkey string, keywords ...string) []string {
	idset := util.NewStringSet()
	if len(keywords) < 1 { return idset.Members() }
	indexc.RLock()
	defer indexc.RUnlock()
	revprefix := indexkey + recsep + "r"
	for _, keyword := range keywords {
		if keyword == "" { continue }
		revkey := []byte(revprefix + recsep + keyword)
		curids, err := indexc.ridb.Get(indexc.ro, revkey)
		if curids != nil && err == nil {
			idset.AddArray(byteArrayToStringArray(recsep, curids))
		}
	}
	return idset.Members()
}

func (indexc *IndexContainer) DelDoc(indexkey string, id string) {
	if id == "" { return }
	indexc.Lock()
	defer indexc.Unlock()
	revprefix := indexkey + recsep + "r"
	forprefix := indexkey + recsep + "f"
	forkey := []byte(forprefix + recsep + id)

	wordrec, err := indexc.fidb.Get(indexc.ro, forkey)
	if wordrec != nil && err == nil {
		words := byteArrayToStringArray(recsep, wordrec)
		for _, word := range words {
			revkey := []byte(revprefix + recsep + word)
			idrec, err := indexc.ridb.Get(indexc.ro, revkey)
			if idrec != nil && err == nil {
				idset := util.NewStringSet()
				idset.AddArray(byteArrayToStringArray(recsep, idrec))
				idset.Del(id)
				if idset.Count() > 0 {
					indexc.ridb.Put(indexc.wo, revkey, stringArrayToByteArray(recsep, idset.Members()))
				} else {
					indexc.ridb.Delete(indexc.wo, revkey)
					indexc.keepcount(indexkey + recsep + "proprevcount", -1)
				}
			}
		}
		indexc.fidb.Delete(indexc.wo, forkey)
		indexc.keepcount(indexkey + recsep + "propforcount", -1)
	}
}

