package db

import (
	"e1/utils"
	"fmt"
)

/*
 * 数据库查询返回的行数据结果集
 */
type RowSet struct {
	Cols []string
	Datas map[string] []byte
}

func (rowSet *RowSet) GetValue(name string, value interface {}) error {
	err := utils.ConvertAssign(value,rowSet.Datas[name])
	if err != nil {
		return err
	}
	return nil
}

func (rowSet *RowSet) GetString(name string) string {
	var result string
	err := rowSet.GetValue(name, &result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (rowSet *RowSet) GetUint64(name string) uint64 {
	var result uint64
	err := rowSet.GetValue(name, &result)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return result
}

func (rowSet *RowSet) GetInt64(name string) int64 {
	var result int64
	err := rowSet.GetValue(name, &result)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return result
}

func (rowSet *RowSet) GetBoolean(name string) bool {
	var result bool
	err := rowSet.GetValue(name, &result)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result
}


