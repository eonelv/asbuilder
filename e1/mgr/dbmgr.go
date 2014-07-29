package mgr

import (
	SQL "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"e1/db"
	_ "e1/log"
	. "e1/log"
)

type DataBaseMgr struct {
	dbServer *SQL.DB
}

var DBMgr DataBaseMgr

func CreateDBMgr(path string) bool {
	db, err := SQL.Open("sqlite3", path)
	if err != nil {
		LogInfo("DataBase Connect Error %s \n", err.Error())
		return false
	}
	DBMgr = DataBaseMgr{}
	DBMgr.dbServer = db
	LogInfo("DataBase connect success.")
	return true
}

func (this *DataBaseMgr) Execute(sql string)(SQL.Result, error) {
	return this.dbServer.Exec(sql)
}

func (this *DataBaseMgr) PreExecute(sql string, args...interface {})(SQL.Result, error) {
	stmt,err := this.dbServer.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (this *DataBaseMgr) Query(sql string, args ...interface{}) ([]*db.RowSet, error) {
	rows, err := this.dbServer.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return this.buildRowSets(rows)
}

func (this *DataBaseMgr) PreQuery(sql string, args ...interface{}) ([]*db.RowSet, error) {
	stmt,err := this.dbServer.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, errq := stmt.Query(args...)

	if errq != nil {
		return nil, errq
	}
	defer  rows.Close()

	return this.buildRowSets(rows)
}

func (this *DataBaseMgr) buildRowSets(rows *SQL.Rows) ([]*db.RowSet, error) {
	colNames,err := rows.Columns()

	if err != nil {
		return nil, err
	}
	var results []*db.RowSet =[]*db.RowSet{}
	var rowSet *db.RowSet
	for rows.Next() {
		values := make([][]byte, len(colNames))
		scans := make([]interface{}, len(colNames))
		for i, _:= range values {
			scans[i] = &values[i]
		}
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		rowSet = &db.RowSet{}
		rowSet.Datas = make(map[string][]byte)
		rowSet.Cols = colNames
		for j, v := range values {
			key := colNames[j]
			rowSet.Datas[key] = v
		}
		results = append(results, rowSet)
	}
	return results, nil
}
