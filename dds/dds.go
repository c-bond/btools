package dds

//const str string = "Hello"

//var flagTest = flag.String("dbname", "db", "a string")
import (
	db "btools/Db"
)

func dds() {

	//resetDb()

	db.InsertConcurrent(8, 1000000)

	// if db.DdsDoesContractorDocExist(543321, 999999) {
	// 	fmt.Println("Doc exists")
	// } else {
	// 	fmt.Println("Doc does not exist")
	// }
}
