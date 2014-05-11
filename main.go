/*
nothing to see here
these are just random go learnings
*/

package main

import (
	//"github.com/lexlapax/graveldb/levelgraph"
	"github.com/lexlapax/graveldb/util"
	"fmt"
	"os"
	"github.com/jmhodges/levigo"
	"bytes"
	"time"
	"reflect"
	"strconv"
	"encoding/binary"
	"regexp"
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

func testReflect() {
	a := []byte{}
	b := "stringvalue"
	c := 1
	d := uint(1)
	fmt.Printf("a: %v, %v\n", a, reflect.TypeOf(a))
	fmt.Printf("b: %v, %v\n", b, reflect.TypeOf(b))
	fmt.Printf("c: %v, %v\n", c, reflect.TypeOf(c))
	fmt.Printf("d: %v, %v\n", d, reflect.TypeOf(d))
	fmt.Printf("")
}

func testInterfaceArgs(args ...interface{}) {
	if len(args) > 0 {
		if aString, found := args[0].(string); found {
			if aString == "" {
				fmt.Printf("NoDirectory = %v\n", aString)
			} else {
				fmt.Printf("all good = %v\n", aString)
			}
		} else {
			fmt.Printf("InvalidParameterValue = %v\n", args[0])
		}
	}
}
func testVariadicArgs() {
	testInterfaceArgs(1)
	testInterfaceArgs(nil)
	testInterfaceArgs("")
	testInterfaceArgs("hello")
}

var testcounter = []byte{}

func incrByte(testcounter []byte) []byte {
	var intcounter uint64
	if testcounter == nil {
		intcounter = uint64(0)
	} else {
		intcounter, _ = binary.Uvarint(testcounter)
	}
	intcounter++ 
	bufsize := binary.Size(intcounter)
	testcounter = make([]byte, bufsize)
	binary.PutUvarint(testcounter, intcounter)
	return testcounter
}

func incrIntStr(counter string) string {
	var intcounter uint64
	var err error
	if counter == "" {
		intcounter = uint64(0)
	} else {
		intcounter, err = strconv.ParseUint(counter, 16, 64)
		if err != nil {intcounter = uint64(0)}
	}
	intcounter++
	retstr := strconv.FormatUint(intcounter, 16)
	return retstr
}

var lastidbyte = []byte(strconv.FormatUint(uint64(89891), 16))
func testStrconv() {
	if lastidbyte == nil { 
		lastidbyte = []byte(strconv.FormatUint(uint64(0), 16)) 
	}
	lastidstr := string(lastidbyte[:])
	intcounter, _ := strconv.ParseUint(lastidstr, 16, 64)
	intcounter++ 
	nextidstr := strconv.FormatUint(intcounter, 16)

	lastidbyte = []byte(nextidstr)

	//strconv.FormatUint(lastid, 16)
	lastid, _ := strconv.ParseUint(lastidstr, 16, 64)
	
	
	nextid, _ := strconv.ParseUint(nextidstr, 16, 64)
	//lastidstr := strconv.FormatUint(lastid, 16)
	//nextid := lastid + 1
	fmt.Printf("last=%v,lastint=%v,next=%v,nextint=%v\n",lastidstr, lastid, nextidstr, nextid)
}

func regexptest() {
	// re := regexp.MustCompile("(.*)[,;.-_]*$")
	re := regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")
	fmt.Printf("v=%v\n", re.ReplaceAllString(",ab-cde;,", ""))

}
var recsep = "\x1f"

func stringArrayToByteArray(record []string) []byte {
	if record == nil { return nil }
	if len(record) < 1 { return nil } 
	bytearr := [][]byte{}
	for _, s := range record {
		bytearr = append(bytearr, []byte(s))
	}
	return bytes.Join(bytearr, []byte(recsep))
}

func byteArrayToStringArray(record []byte) []string {
	strings := []string{}
	if record == nil {return strings}
	bytearray := bytes.Split(record, []byte(recsep))
	for _, arr := range bytearray {
		strings = append(strings, string(arr[:]))
	}
	return strings
}

func bytearraytest() {
	strarr := []string{"a", "b", "c"}
	byterecord := stringArrayToByteArray(strarr)
	newstrarr := byteArrayToStringArray(byterecord)
	fmt.Printf("s=%v,\nb=%v,\nn=%v\n", strarr, byterecord, newstrarr)
}

func uuidtest() {
	uuid, err := util.UUID()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("%s\n", uuid)
}

func testIterChannelWrite() <-chan string {
	ch := make(chan string)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- fmt.Sprintf("i am %v", i)
		}
		close(ch)
	}()
	return ch
}

func testIterChannel() {
	for s := range testIterChannelWrite() {
		fmt.Println(s)
	}
	time.Sleep(1 * 1.e9)
}

func testWriteChannel(cs chan int) {
	for i := 1; i <= 10; i+= 1 {
		cs <- i
		fmt.Printf("sent %v\n", i)

	}
}


func testReadChannel(cs <-chan int, arr *[]int) {
	for s:= range cs {
		*arr = append(*arr, s)
		fmt.Printf("received = %v\n",s)
	}
}


func testChannel() {
	cs := make(chan int)
	intarr := []int{}
	go testWriteChannel(cs)
	go testReadChannel(cs, &intarr)
	time.Sleep(1 * 1.e9)
	fmt.Printf("%v\n", intarr)
}

func main() {
	testIterChannel()
}



