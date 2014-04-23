package main
import (
	//"github.com/lexlapax/graveldb/levelgraph"
	//"github.com/lexlapax/graveldb/core"
	"fmt"
	"os"
	"github.com/jmhodges/levigo"
	"bytes"
	"time"
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


func joinBytes(sep []byte, elements ...[]byte) ([]byte) {
	if len(elements) < 1 { return []byte{} } 
	return bytes.Join(elements, sep)
}

func splitBytes(sep []byte, elements []byte) (int, [][]byte) {
	if elements == nil {return 0, nil}
	elementarr := bytes.Split(elements, sep)
	n := len(elementarr)
	return n, elementarr
}



func BytesTest() {
	a := []byte("a")
	b := []byte("b")
	c := []byte("c")
	d := []byte("d")
	sep := []byte("\x1f")
	
	joint := joinBytes(sep)
	fmt.Printf("0=%v, %v\n", joint, string(joint[:]))

	joint = joinBytes(sep, a)
	fmt.Printf("1=%v, %v\n", joint, string(joint[:]))
	n, split := splitBytes(sep, joint)
	fmt.Printf("1=%v, %v\n", n, split)


	joint = joinBytes(sep, a, b)
	fmt.Printf("2=%v, %v\n", joint, string(joint[:]))
	n, split = splitBytes(sep, joint)
	fmt.Printf("2=%v, %v\n", n, split)

	joint = joinBytes(sep, a, b, c)
	fmt.Printf("3=%v, %v\n", joint, string(joint[:]))
	n, split = splitBytes(sep, joint)
	fmt.Printf("3=%v, %v\n", n, split)

	joint = joinBytes(sep, a, b, c, d)
	fmt.Printf("4=%v, %v\n", joint, string(joint[:]))
	n, split = splitBytes(sep, joint)
	fmt.Printf("4=%v, %v\n", n, split)


	// combined := [][]byte{}
	// combined = append(combined,a,b,c,d)
	// flat := bytes.Join(combined, sep)
	
	// //combined = append(combined,b)
	// fmt.Printf("combined=%v\n", combined)
	// fmt.Printf("flat=%v\n", flat)
	// split := bytes.Split(flat, sep)
	// fmt.Printf("splitindiv=%v %v %v\n", split[0], split[1], split[2])
	// for _, a := range bytes.Split(flat,sep) {
	// 	fmt.Printf("%v\n", string(a[:]))
	// }
	// fmt.Printf("split=%v\n", split)

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
func testOutChannel(cs chan int) {
	for i := 1; i <= 10; i+= 1 {
		cs <- i
		fmt.Printf("sent %v\n", i)

	}
}


func iterEdgeSet(cs <-chan int, edges []int) {
	for s:= range cs {
		edges = append(edges, s)
	}
}

func testReadChannel(cs chan int) {
	for i := range cs {
		fmt.Printf("received = %v\n",i)
	}
}

func main() {
	cs := make(chan int)
	edges := []int{}
	go testOutChannel(cs)
	go iterEdgeSet(cs, edges)
	time.Sleep(1 * 1.e9)
	fmt.Printf("%v\n", edges)
}



