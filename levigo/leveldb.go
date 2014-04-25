package levigo


import ("bytes"
		"fmt"
		"os"
		"github.com/jmhodges/levigo"
		)

func CheckGet(db *levigo.DB, roptions *levigo.ReadOptions, key, expected []byte) {
	getValue, err := db.Get(roptions, key)

	if err != nil {
		fmt.Fprint(os.Stderr, "Get failed: %v", err)
	}
	if !bytes.Equal(getValue, expected) {
		fmt.Println(os.Stderr, "expected Get value %v, got %v", expected, getValue)
	}
}


func DbTest() {
	fmt.Println("yo starting")
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	fmt.Fprintln(os.Stdout, "creating filter")
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)
	ro := levigo.NewReadOptions()
	wo := levigo.NewWriteOptions()
	db, err := levigo.Open("./test.db", opts)
	putKey := []byte("foo")
	putValue := []byte("hello yo mama")
	fmt.Fprintln(os.Stdout, "putting value")
	err = db.Put(wo, putKey, putValue)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Put failed: %v", err)
	}
	fmt.Fprintln(os.Stdout, "checking get value")
	var getValueBytes []byte
	getValueBytes, err = db.Get(ro, putKey)
	//CheckGet(db, ro, putKey, putValue)
	fmt.Print("value=")
	fmt.Println(string(getValueBytes[:]))
}



