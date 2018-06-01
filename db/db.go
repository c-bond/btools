package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type SwitchTableInfo struct {
	name, revision                                               string
	clientUUID, contractorUUID                                   string
	clientProjno, clientDocno, contractorProjno, contractorDocno int
}

type myError struct {
	When time.Time
	What string
}

const dbName string = "user=l37 password=l37 dbname=btools"
const timeoutAttempts = 2

var (
	db  *sql.DB
	err error
)

func (e myError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func openConnection() error {
	if db, err = sql.Open("postgres", dbName); err != nil {
		return err
	}

	return db.Ping()
}

func checkConnectionNTimes(trys int) error {
	if db != nil && db.Ping() != nil {
		return nil
	} else if openConnection() != nil {
		time.Sleep(1 * time.Second)
		if trys > 0 {
			checkConnectionNTimes(trys - 1)
			fmt.Println(trys)
		} else {

			return myError{time.Now(), `Failed to connect to database, 
				service timed out after ` + strconv.Itoa(timeoutAttempts) + ` times`}
		}
	}
	return nil
}

func checkConnection() {
	if err = checkConnectionNTimes(timeoutAttempts); err != nil {
		fmt.Println(err)
	}
}

func countRows() {
	checkConnection()
	count := 0
	err = db.QueryRow(`select count(*) from doc_switch`).Scan(&count)
	fmt.Println(count)
}

func SelectDocDDS(guidIn string) (SwitchTableInfo, error) {
	checkConnection()
	var d SwitchTableInfo
	query := `select * from doc_switch where (contractor_uuid = '` + guidIn + `'
						OR client_uuid = '` + guidIn + `')`
	err = db.QueryRow(query).Scan(&d.name, &d.revision, &d.clientUUID, &d.clientProjno, &d.clientDocno,
		&d.contractorUUID, &d.contractorProjno, &d.contractorDocno)
	return d, err
}

func DdsDoesClientDocExist(projno int32, docno int32) bool {
	checkConnection()
	op := 0
	stmt := `SELECT client_docno FROM doc_switch WHERE client_projno = $1
			 AND client_docno = $2`
	db.QueryRow(stmt, projno, docno).Scan(&op)
	if op == 0 {
		return false
	}
	return true
}

func DdsDoesContractorDocExist(projno int, docno int) bool {
	checkConnection()
	op := 0
	stmt := `SELECT contractor_docno FROM doc_switch WHERE contractor_projno = $1
			 AND contractor_docno = $2`
	db.QueryRow(stmt, projno, docno).Scan(&op)
	if op == 0 {
		return false
	}
	return true
}

func insertTestRows(start int, count int, wg *sync.WaitGroup) (rowsInserted int64) {
	defer wg.Done()
	dbs, errs := sql.Open("postgres", dbName)
	if err != nil {
		return
	}
	dbs.SetMaxOpenConns(1)
	dbs.Exec("BEGIN TRANSACTION")
	defer dbs.Exec("END TRANSACTION")

	fstr := `INSERT INTO doc_switch (name, revision,
		client_uuid, client_projno,  client_docno, contractor_uuid,
		contractor_projno, contractor_docno)
		VALUES('doc0%v','P01', uuid_generate_v4(),
		12345, %v, uuid_generate_v4(), 543321, %v)`
	for i := start; i <= count; i++ {
		stmt := fmt.Sprintf(fstr, i, i, i)
		if _, errs = dbs.Exec(stmt); errs != nil {
			if errs.Error() == "database is locked" {
				i--
			} else {
				//log it
			}
		} else {
			rowsInserted++
		}
	}
	return
}

func InsertConcurrent(threads int, count int) {
	deleteAllRecords()
	split := count / threads
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go insertTestRows(i*split, ((i*split)+split)-1, &wg)
	}
	wg.Wait()
	fmt.Println("done")
}

func deleteAllRecords() (rowsDeleted int64) {
	checkConnection()
	query := `delete from doc_switch`
	res, err := db.Exec(query)
	if err != nil {

		log.Fatal(err)
	}
	rowsDeleted, _ = res.RowsAffected()
	return
}

func initDb() {
	checkConnection()
	stmt := `create table doc_switch (
		name text not null,
		revision text not null,
		client_uuid UUID not null, 
		client_projno INTEGER not null, 
		client_docno INTEGER not null,
		contractor_uuid UUID not null, 
		contractor_projno INTEGER not null, 
		contractor_docno INTEGER not null,		
		CONSTRAINT name_rev UNIQUE(name, revision));`
	if _, err = db.Exec(stmt); err != nil {
		fmt.Println(err)
	}
	//logit
	// }
	// stmt = `CREATE INDEX clientguid_index
	// 				ON doc_switch (client_guid)`
	// if _, err = db.Exec(stmt); err != nil {
	// 	//logit
	// }
	// stmt = `CREATE INDEX contractorguid_index
	// 				ON doc_switch (contractor_guid)`
	// if _, err = db.Exec(stmt); err != nil {
	//logit
	//}
}

func ResetDb() {
	checkConnection()
	stmt := `drop table doc_switch`
	db.Exec(stmt)
	initDb()
}
