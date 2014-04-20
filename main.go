package main
import (
	//"github.com/lexlapax/graveldb/levelgraph"
	//"github.com/lexlapax/graveldb/core"
	"fmt"
	"os"
	"github.com/jmhodges/levigo"
	"bytes"
)

func writeToDb(db *levigo.DB) {
	fmt.Fprintln(os.Stdout, "putting value")
	wo := levigo.NewWriteOptions()
	alpha := []string{"a","b","c","d","e"}
	numb := []string{"1","2","3","4","5"}
	//var err error
	for _,a := range alpha {
		for _,b := range alpha {
			for _,n := range numb { 
				putKey := []byte(a + "::" + b + "::" + n)
				putValue := []byte("hello yo mama")
				_ = db.Put(wo, putKey, putValue)
			}
		}
	}
}

func readFromDb(db *levigo.DB){
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.NewIterator(ro)
	defer it.Close()
	prefix := append([]byte("a"), []byte("::")...)
	it.Seek(prefix)
	keys := [][]byte{}
	for it = it; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
		keys = append(keys, it.Key())
//		fmt.Printf("%s = %s\n", it.Key(), it.Value())
	}
	for _, k := range keys {
		fmt.Printf("%s = \n", k)
	}


}

func BytesTest() {
	a := []byte("a")
	b := []byte("b")
	c := []byte("c")
	d := []byte("")
	sep := []byte("\x1f")
	combined := [][]byte{}
	combined = append(combined,a,b,c,d)
	flat := bytes.Join(combined, sep)
	
	//combined = append(combined,b)
	fmt.Printf("combined=%v\n", combined)
	fmt.Printf("flat=%v\n", flat)
	split := bytes.Split(flat, sep)
	fmt.Printf("splitindiv=%v %v %v\n", split[0], split[1], split[2])
	for _, a := range bytes.Split(flat,sep) {
		fmt.Printf("%v\n", string(a[:]))
	}
	fmt.Printf("split=%v\n", split)

}

func DbTest() {
	fmt.Println("yo starting")
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	fmt.Fprintln(os.Stdout, "creating filter")
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)
	//ro := levigo.NewReadOptions()
	//wo := levigo.NewWriteOptions()
	db, _ := levigo.Open("./test.db", opts)
	//writeToDb(db)
	readFromDb(db)
	//var getValueBytes []byte
	//getValueBytes, err = db.Get(ro, putKey)
	//CheckGet(db, ro, putKey, putValue)
	//fmt.Print("value=")
	//fmt.Println(string(getValueBytes[:]))
	db.Close()
}


func main() {
	BytesTest()
}



