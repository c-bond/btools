package main

//const str string = "Hello"

//var flagTest = flag.String("dbname", "db", "a string")
import (
	db "btools/db"
	q "mytypes/queue"
)

type docInfo struct {
	name, revision                                               string
	clientUUID, contractorUUID                                   string
	clientProjno, clientDocno, contractorProjno, contractorDocno int
}

var queue q.Queue

func taskManager() {
	var doc docInfo
	queue.Enqueue(doc)
}

func scanner(docguid string) {
	//guid := "doc4436587clientguidx" //get from pwise
	// if doc, err := selectDoc(docguid); err != nil {
	// 	fmt.Println("not found")
	// } else {
	// 	fmt.Println(doc.name)
	// }
}

func dds() {

	//resetDb()

	db.InsertConcurrent(8, 1000000)

	// if db.DdsDoesContractorDocExist(543321, 999999) {
	// 	fmt.Println("Doc exists")
	// } else {
	// 	fmt.Println("Doc does not exist")
	// }
}
