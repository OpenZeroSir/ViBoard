package Editor

import (
	"MinView/Lib"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// 编辑器数据库，用于保存数据和附件
type EditorSQL struct {
	db    *sql.DB
	Mutex *sync.Mutex
}

func NewSQL() (*EditorSQL, error) {
	var s EditorSQL
	var err = s.init()
	if err != nil {
		return nil, err
	}
	s.Mutex = &sync.Mutex{}
	return &s, nil
}

func (s *EditorSQL) init() error {
	db, err := sql.Open("sqlite3", ":memory:")
	// db, err := sql.Open("sqlite3", "a.db")
	if err != nil {
		return err
	}
	s.db = db
	//初始化媒体表，用于存贮图片、视频、音频和附件文件
	//@viewid 视图ID
	//@compid 组件ID
	sql := `CREATE TABLE IF NOT EXISTS "MEDIA" (
		"ViewID"	TEXT NOT NULL,
		"ComID"	TEXT NOT NULL,
		"Name"	TEXT NOT NULL,
		"Data"	BLOB NOT NULL,
		PRIMARY KEY("ViewID","ComID")
	);`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	//初始化映射表，因为一个Excel文件可能有很多表，一个文件对应一个key,一个key对应多个子表的key
	//csv就一个文件是一个key,一个key对应一个表的key
	//volume类似磁盘映射对应盘符
	//@PriKey文件key
	//@SecKey Sheet key
	sql = `CREATE TABLE IF NOT EXISTS "VOLUME" (
		"PriKey"	TEXT NOT NULL,
		"FileName"	TEXT NOT NULL,
		"SecKey"	TEXT NOT NULL,
		"SheetName"	BLOB NOT NULL,
		PRIMARY KEY("PriKey","SecKey")
	);`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
func (s *EditorSQL) GetDB() *sql.DB {
	s.Mutex.Lock()
	return s.db
}

func (s *EditorSQL) Close() error {
	return s.db.Close()
}

// 备份数据库
func (s *EditorSQL) Backup() (string, error) {
	file_name, err := Lib.BackupDatabase(s.GetDB())
	s.Mutex.Unlock()
	if err != nil {
		return "", err
	}
	// fmt.Println(file_name)
	return file_name, nil
}

// 还原数据库,从[]byte数据还原
func (s *EditorSQL) Restore(data []byte) error {
	tmpFile, err := ioutil.TempFile(``, `NNI-*`)
	if err != nil {
		return err
	}
	name := tmpFile.Name()
	_, err = tmpFile.Write(data)
	if err != nil {
		err1 := tmpFile.Close()
		if err1 != nil {
			return err1
		}
		return err
	}
	err1 := tmpFile.Close()
	if err1 != nil {
		return err1
	}
	err = Lib.RestoreDatabase(name, s.GetDB())
	s.Mutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

// 清空数据库
func (s *EditorSQL) Clear() error {
	sql := "SELECT PriKey,SecKey FROM 'VOLUME'"
	rows, err := s.GetDB().Query(sql)
	if err != nil {
		s.Mutex.Unlock()
		return err
	}
	// defer rows.Close()
	// s.Mutex.Unlock()
	//获取要删除的表名称
	var drop_list []string
	for rows.Next() {
		var (
			priKey string
			secKey string
		)
		err = rows.Scan(&priKey, &secKey)
		if err != nil {
			if e := rows.Close(); e != nil {
				s.Mutex.Unlock()
				return e
			}
			s.Mutex.Unlock()
			return err
		}
		drop_list = append(drop_list, secKey)
		//不能在这里删除是因为可能是bug,rows没读完又删除表，提示找不到表。
	}
	if err = rows.Err(); err != nil {
		s.Mutex.Unlock()
		return err
	}
	for _, name := range drop_list {
		sql = "DROP TABLE `" + name + "`;"
		_, err = s.db.Exec(sql)
		if err != nil {
			s.Mutex.Unlock()
			return err
		}
	}
	//清空MEDIA表和VOLUME表
	sql = "DELETE FROM 'MEDIA'"
	_, err = s.db.Exec(sql)
	if err != nil {
		s.Mutex.Unlock()
		return err
	}
	sql = "DELETE FROM 'VOLUME'"
	_, err = s.db.Exec(sql)
	if err != nil {
		s.Mutex.Unlock()
		return err
	}
	s.Mutex.Unlock()
	return nil
}

type VOLUMEResultStruct struct {
	//文件Key
	PriKey string `json:"PriKey"`
	//文件名称
	FileName string `json:"FileName"`
	//表key
	SecKey string `json:"SecKey"`
	//表名称
	SheetName string `json:"SheetName"`
	//错误信息
	// Error string `json:"Error"`
}

type VOLUMEResultMap map[string][]VOLUMEResultStruct

// 获取VOLUME表的所有内容,
// 结构 PriKey:FileName:SecKey:SheetName
func (s *EditorSQL) GetVOLUME() (VOLUMEResultMap, error) {
	sql := "SELECT PriKey,FileName,SecKey,SheetName FROM 'VOLUME'"
	rows, err := s.GetDB().Query(sql)
	if err != nil {
		s.Mutex.Unlock()
		return nil, err
	}
	// defer rows.Close()
	// s.Mutex.Unlock()
	vmap := make(VOLUMEResultMap)
	//获取要删除的表名称
	// var name_list []string
	for rows.Next() {
		var (
			priKey    string
			FileName  string
			secKey    string
			SheetName string
		)
		err = rows.Scan(&priKey, &FileName, &secKey, &SheetName)
		if err != nil {
			if e := rows.Close(); e != nil {
				s.Mutex.Unlock()
				return nil, e
			}
			s.Mutex.Unlock()
			return nil, err
		}
		// name_list = append(name_list, priKey+":"+FileName+":"+secKey+":"+SheetName)
		//不能在这里删除是因为可能是bug,rows没读完又删除表，提示找不到表。
		var vrs VOLUMEResultStruct
		vrs.PriKey = priKey
		vrs.FileName = FileName
		vrs.SecKey = secKey
		vrs.SheetName = SheetName
		vmap[priKey] = append(vmap[priKey], vrs)
	}
	if err = rows.Err(); err != nil {
		if e := rows.Close(); e != nil {
			s.Mutex.Unlock()
			return nil, e
		}
		s.Mutex.Unlock()
		return nil, err
	}
	s.Mutex.Unlock()
	return vmap, nil
}

// 通过表key或文件key获得，数据表数量，
// 一个表Key对一个表，但一个文件key可能对多个表，也就是说会有多个表key
func (x EditorSQL) GetVOLUMEInfoByKey(key string) (map[string]string, error) {
	result := make(map[string]string)
	sql := "SELECT PriKey,SecKey,SheetName FROM 'VOLUME' WHERE PriKey='" + key + "' OR SecKey='" + key + "'"
	rows, err := x.GetDB().Query(sql)
	if err != nil {
		x.Mutex.Unlock()
		return nil, err
	}
	// defer rows.Close()
	// x.Mutex.Unlock()
	for rows.Next() {
		var (
			priKey    string
			secKey    string
			sheetName string
		)
		err = rows.Scan(&priKey, &secKey, &sheetName)
		if err != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return nil, e
			}
			x.Mutex.Unlock()
			return nil, err
		}
		result[secKey] = sheetName
	}
	if e := rows.Close(); e != nil {
		x.Mutex.Unlock()
		return nil, e
	}
	x.Mutex.Unlock()
	return result, nil
}

// 获取数据库字段和字段类型,格式：字段名称：字段类型
func (x *EditorSQL) GetFieldAndTypeByKey(key string) ([]string, error) {
	var result []string
	query := fmt.Sprintf("PRAGMA table_info('%s')", key)
	rows, err := x.GetDB().Query(query)
	if err != nil {
		x.Mutex.Unlock()
		return nil, err
	}
	// defer rows.Close()
	// x.Mutex.Unlock()
	// 解析查询结果，获取字段和类型
	for rows.Next() {
		var cid int
		var name string
		var typename string
		var notnull int
		var dfltValue sql.NullString
		var pk int

		err := rows.Scan(&cid, &name, &typename, &notnull, &dfltValue, &pk)
		if err != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return nil, e
			}
			x.Mutex.Unlock()
			return nil, err
		}
		result = append(result, name+":"+typename)
	}
	if err = rows.Err(); err != nil {
		x.Mutex.Unlock()
		return nil, err
	}
	x.Mutex.Unlock()
	return result, nil
}

// 通过远程post数据更新显示器数据库
func (x *EditorSQL) RemotePostUpdateDB(key string, replace bool, data string) error {
	//事务的开始,这里不能用GetFieldAndTypeByKey来处理，因为事务前执行的查询可能导致rows关闭不了影响后面次操作。
	tx, err := x.GetDB().Begin()
	if err != nil {
		return err
	}
	var sql_field []string
	var sql_type []string
	query := fmt.Sprintf("PRAGMA table_info('%s')", key)
	rows, err := tx.Query(query)
	if err != nil {
		// rollerr := tx.Rollback() // 回滚事务
		// if rollerr != nil {
		// 	x.Mutex.Unlock()
		// 	return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
		// }
		x.Mutex.Unlock()
		return err
	}
	// defer rows.Close()
	// x.Mutex.Unlock()
	// 解析查询结果，获取字段和类型
	for rows.Next() {
		var cid int
		var name string
		var typename string
		var notnull int
		var dfltValue sql.NullString
		var pk int

		err := rows.Scan(&cid, &name, &typename, &notnull, &dfltValue, &pk)
		if err != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return e
			}
			x.Mutex.Unlock()
			return err
		}
		sql_field = append(sql_field, name)
		sql_type = append(sql_type, typename)
	}
	if err = rows.Err(); err != nil {
		rollerr := tx.Rollback() // 回滚事务
		if rollerr != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return e
			}
			x.Mutex.Unlock()
			return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
		}
		if e := rows.Close(); e != nil {
			x.Mutex.Unlock()
			return e
		}
		x.Mutex.Unlock()
		return err
	}
	if len(sql_field) == 0 {
		rollerr := tx.Rollback() // 回滚事务
		if rollerr != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return e
			}
			x.Mutex.Unlock()
			return errors.New("处理出错：获取到数据库表字段数量：0" + " 回滚出错：" + rollerr.Error())
		}
		if e := rows.Close(); e != nil {
			x.Mutex.Unlock()
			return e
		}
		x.Mutex.Unlock()
		return errors.New("获取到数据库表字段数量：0")
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		rollerr := tx.Rollback() // 回滚事务
		if rollerr != nil {
			if e := rows.Close(); e != nil {
				x.Mutex.Unlock()
				return e
			}
			x.Mutex.Unlock()
			return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
		}
		if e := rows.Close(); e != nil {
			x.Mutex.Unlock()
			return e
		}
		x.Mutex.Unlock()
		return err
	}
	if e := rows.Close(); e != nil {
		x.Mutex.Unlock()
		return e
	}
	x.Mutex.Unlock()
	//事务的开始
	tx, err = x.GetDB().Begin()
	if err != nil {
		x.Mutex.Unlock()
		return err
	}
	//如果是覆盖表，需要先尝试删除表，因为可能数据库里有表，也可能没有表。
	if replace {
		//删除表
		sql := Lib.GenerateDropSQL(key)
		_, err = tx.Exec(sql)
		if err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				x.Mutex.Unlock()
				return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
			}
			x.Mutex.Unlock()
			return err
		}
		// 创建表格
		sql = Lib.GenerateCreateTableSQL(key, sql_field, sql_type)
		_, err = tx.Exec(sql)
		if err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				x.Mutex.Unlock()
				return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
			}
			x.Mutex.Unlock()
			return err
		}
	}

	var new_data []interface{}
	// 解析收到的字符串为JSON格式
	err = json.Unmarshal([]byte(data), &new_data)
	if err != nil {
		x.Mutex.Unlock()
		return err
	}
	for _, row1 := range new_data {
		if arr, ok := row1.([]interface{}); ok {
			var row []string
			// 遍历数组
			for _, value := range arr {
				switch v := value.(type) {
				case float64:
					num_str := strconv.FormatFloat(v, 'f', 8, 64)
					row = append(row, num_str)
				case string:
					row = append(row, v)
				default:
					row = append(row, "")
				}
			}
			if len(row) != len(sql_field) {
				rollerr := tx.Rollback() // 回滚事务
				if rollerr != nil {
					x.Mutex.Unlock()
					return errors.New("处理出错：提交数据列数和数据库字段数量不一致 " + " 回滚出错：" + rollerr.Error())
				}
				x.Mutex.Unlock()
				return errors.New("提交数据列数和数据库字段数量不一致")
			}
			sql := Lib.GenerateInsertSQL(key, sql_field, sql_type, row)
			_, err = tx.Exec(sql)
			if err != nil {
				rollerr := tx.Rollback() // 回滚事务
				if rollerr != nil {
					x.Mutex.Unlock()
					return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
				}
				x.Mutex.Unlock()
				return err
			}
		} else {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				x.Mutex.Unlock()
				return errors.New("处理出错：行数据不是数组类型 " + " 回滚出错：" + rollerr.Error())
			}
			x.Mutex.Unlock()
			return errors.New("行数据不是数组类型")
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		rollerr := tx.Rollback() // 回滚事务
		if rollerr != nil {
			x.Mutex.Unlock()
			return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
		}
		x.Mutex.Unlock()
		return err
	}
	x.Mutex.Unlock()
	return nil
}
