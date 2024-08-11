package Editor

import (
	"MinView/Lib"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
)

// 处理xlsx和csv文件
// 需要把xlsx和csv文件读取到数据库中，
// csv文件需要支持：
//
//		1、UTF-8
//		2、UTF-8 BOM
//	    3、GBK
type EditorXlsx struct {
	ExcelFile *excelize.File
	// CsvFile   *os.File
	//需要保存文件名称，方便映射到数据库
	FileName string
	//文件类型，xlsx csv,方便读取的时候是用哪个函数读取
	Ext string
}

// 新建EditorXlsx对象，用于读取Excel xlsx文件
func NewEditorXlsx() *EditorXlsx {
	return &EditorXlsx{ExcelFile: nil, FileName: "", Ext: ""}
}

// 打开文件
func (x *EditorXlsx) Open(file_path string) error {
	ext := filepath.Ext(file_path)
	if ext != ".xlsx" {
		return errors.New("不是.xlsx格式文件")
	}
	err := x.openSheet(file_path)
	if err != nil {
		return err
	}
	// } else if ext == ".csv" {
	// 	err := x.openCsv(file_path)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// 通过路径加载xlsx文件到EditorXlsx对象中
func (x *EditorXlsx) openSheet(file_path string) error {
	excel, err := excelize.OpenFile(file_path)
	if err != nil {
		return err
	}
	x.ExcelFile = excel
	fileName := filepath.Base(file_path)
	x.FileName = strings.Replace(fileName, ".xlsx", "", -1)
	x.Ext = "xlsx"
	return nil
}

// 通过路径加载csv
// func (x *EditorXlsx) openCsv(file_path string) error {
// 	iCsv, err := os.Open(file_path)
// 	if err != nil {
// 		return err
// 	}
// 	x.CsvFile = iCsv
// 	fileName := filepath.Base(file_path)
// 	x.FileName = strings.Replace(fileName, ".csv", "", -1)
// 	x.Ext = "csv"
// 	return nil
// }

type SheetInfo struct {
	//文件key
	PriKey string
	//文件名称
	FileName string
	// Sheet key
	SecKey string
	//Sheet名称
	SheetName string
}

// 读取所有Sheet数据到数据库，空表需要跳过，
// 导入到数据库表中，所有表都会在数据库中生成key表，
// @items数据库表名和excel表名对应关系,
// 如果items为空那就是新建数据库表，
// @replace 用于表示如果数据库中存在对应表，是增量更新数据还是覆盖数据
func (x *EditorXlsx) CopyDataToDB(items []SheetInfo, db *sql.DB, replace bool) error {
	defer func() {
		if err := x.ExcelFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	fileNames := x.ExcelFile.GetSheetList()
	//如何items为空，那就直接创建
	if items == nil {
		for _, sheet := range fileNames {
			pkey := Lib.ReplaceFirstNChars(uuid.NewV4().String(), "e", 1)
			skey := Lib.ReplaceFirstNChars(uuid.NewV4().String(), "e", 1)
			// fmt.Println(pkey, skey)
			item := SheetInfo{
				PriKey:    pkey,
				FileName:  x.FileName,
				SecKey:    skey,
				SheetName: sheet,
			}
			//这里直接是true,就是确保，不管数据库里有没有这个表明都会尝试删除后新建表
			err := x.createCopy(item, db, true)
			if err != nil {
				return err
			}
		}
	} else {
		for _, item := range items {
			//这里需要由replace来决定是否替换表，不管有没有表，应当尝试创建表。
			err := x.createCopy(item, db, replace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 创建表并复制数据
func (x *EditorXlsx) createCopy(item SheetInfo, db *sql.DB, replace bool) error {

	rows, err := x.ExcelFile.Rows(item.SheetName)
	if err != nil {
		return err
	}
	// defer rows.Close()
	if !rows.Next() {
		return errors.New(item.SheetName + " 没有检测到可用标头")
	}
	headName, err := rows.Columns()
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New(item.SheetName + " 没有检测到可用数据")
	}
	row1, err := rows.Columns()
	if err != nil {
		return err
	}
	for {
		if len(row1) < len(headName) {
			row1 = append(row1, "")
		} else {
			break
		}
	}
	var sheetHead []string
	var headType []string
	for index, value := range row1 {
		currentType := "REAL"
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			currentType = "TEXT"
		}
		sheetHead = append(sheetHead, headName[index])
		headType = append(headType, currentType)
	}
	//如果是覆盖表，需要先尝试删除表，因为可能数据库里有表，也可能没有表。
	if replace {
		//删除表
		sql := Lib.GenerateDropSQL(item.SecKey)
		_, err = db.Exec(sql)
		if err != nil {
			return err
		}
		err = x.updateVOLUME(item, db, true)
		if err != nil {
			return err
		}
	}
	// 创建表格
	sql := Lib.GenerateCreateTableSQL(item.SecKey, sheetHead, headType)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	//插入数据
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sql = Lib.GenerateInsertSQL(item.SecKey, sheetHead, headType, row1)
	_, err = tx.Exec(sql)
	if err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	for rows.Next() {
		row1, err = rows.Columns()
		if err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
		for {
			if len(row1) != len(headName) {
				row1 = append(row1, "")
			} else {
				break
			}
		}
		sql = Lib.GenerateInsertSQL(item.SecKey, sheetHead, headType, row1)
		_, err = tx.Exec(sql)
		if err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	err = x.updateVOLUME(item, db, false)
	if err != nil {
		return err
	}
	return nil
}

// 更新VOLUME表的name和key
// @delete 有可能是更新，有可能是删除
func (x *EditorXlsx) updateVOLUME(item SheetInfo, db *sql.DB, delete bool) error {
	sql := ""
	if delete {
		sql = "DELETE FROM 'VOLUME' WHERE PriKey = '" + item.PriKey + "' and SecKey = '" + item.SecKey + "'"
	} else {
		sql = "REPLACE INTO 'VOLUME' VALUES ('" + item.PriKey + "', '" + item.FileName + "', '" + item.SecKey + "','" + item.SheetName + "')"
	}
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// 更新数据库表，数据库已经有对应key的表，但有时候是增量更形，
// 有时候是覆盖更新表，和CopyDataToDB不一样，CopyDataToDB会初始化所有数据库表，生成表key,
// 但这个是更新数据库表，表明key不会有变化。
// @item类型是[sheetkey]sheetname,如果长度是1 直接忽略sheetname,用sheet0的数据,
// 如果是多个表，需要比较表名
func (x *EditorXlsx) UpdateDataToDB(items map[string]string, SQL *EditorSQL, replace bool) error {
	defer func() {
		if err := x.ExcelFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for key := range items {
		sheet := ""
		if len(items) == 1 {
			sheet = x.ExcelFile.GetSheetList()[0]
		} else {
			sheet = items[key]
			//检测数据库的表名是否在表格中，因为多个表的表名需要和数据库中表明对应
			sheet_has_db_name := false
			for _, name := range x.ExcelFile.GetSheetList() {
				if sheet == name {
					sheet_has_db_name = true
					break
				}
			}
			if !sheet_has_db_name {
				return errors.New("表名 " + sheet + " 在工作簿中找不到")
			}
		}
		rows, err := x.ExcelFile.Rows(sheet)
		if err != nil {
			SQL.Mutex.Unlock()
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			SQL.Mutex.Unlock()
			return errors.New(sheet + " 没有检测到可用表头")
		}
		headName, err := rows.Columns()
		if err != nil {
			return err
		}
		//事务的开始,这里不能用GetFieldAndTypeByKey来处理，因为事务前执行的查询可能导致rows关闭不了影响后面次操作。
		tx, err := SQL.GetDB().Begin()
		if err != nil {
			SQL.Mutex.Unlock()
			return err
		}
		var sql_field []string
		var sql_type []string
		query := fmt.Sprintf("PRAGMA table_info('%s')", key)
		sqlrows, err := tx.Query(query)
		if err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
			}
			if e := sqlrows.Close(); e != nil {
				SQL.Mutex.Unlock()
				return e
			}
			SQL.Mutex.Unlock()
			return errors.New("获取表字段失败" + err.Error())
		}
		// defer sqlrows.Close()

		// 解析查询结果，获取字段和类型
		for sqlrows.Next() {
			var cid int
			var name string
			var typename string
			var notnull int
			var dfltValue sql.NullString
			var pk int
			err := sqlrows.Scan(&cid, &name, &typename, &notnull, &dfltValue, &pk)
			if err != nil {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return err
			}
			sql_field = append(sql_field, name)
			sql_type = append(sql_type, typename)
		}
		if err = sqlrows.Err(); err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
			}
			if e := sqlrows.Close(); e != nil {
				SQL.Mutex.Unlock()
				return e
			}
			SQL.Mutex.Unlock()
			return errors.New("sqlrows读取失败" + err.Error())
		}
		if len(sql_field) == 0 {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return errors.New("处理出错：获取到数据库表字段数量：0" + " 回滚出错：" + rollerr.Error())
			}
			if e := sqlrows.Close(); e != nil {
				SQL.Mutex.Unlock()
				return e
			}
			SQL.Mutex.Unlock()
			return errors.New("获取到数据库表字段数量：0")
		}
		// 提交事务
		err = tx.Commit()
		if err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return errors.New("处理出错：" + err.Error() + " 回滚出错：" + rollerr.Error())
			}
			if e := sqlrows.Close(); e != nil {
				SQL.Mutex.Unlock()
				return e
			}
			SQL.Mutex.Unlock()
			return err
		}
		//检测表头是否在数据库表的表头中
		if len(sql_field) != len(headName) {
			if e := sqlrows.Close(); e != nil {
				SQL.Mutex.Unlock()
				return e
			}
			SQL.Mutex.Unlock()
			return errors.New("数据字段数量和数据库表字段数量不一致")
		}

		for _, field := range sql_field {

			has_field := false
			for _, head := range headName {
				if head == field {
					has_field = true
					break
				}
			}
			if !has_field {
				if e := sqlrows.Close(); e != nil {
					SQL.Mutex.Unlock()
					return e
				}
				SQL.Mutex.Unlock()
				return errors.New(field + " 字段在表 " + sheet + " 中没有找到对应的列")
			}
		}
		SQL.Mutex.Unlock()
		//事务的开始
		tx, err = SQL.GetDB().Begin()
		if err != nil {
			SQL.Mutex.Unlock()
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
					SQL.Mutex.Unlock()
					return errors.New(err.Error() + " " + rollerr.Error())
				}
				SQL.Mutex.Unlock()
				return err
			}
			// 创建表格
			sql = Lib.GenerateCreateTableSQL(key, sql_field, sql_type)
			_, err = tx.Exec(sql)
			if err != nil {
				rollerr := tx.Rollback() // 回滚事务
				if rollerr != nil {
					SQL.Mutex.Unlock()
					return errors.New(err.Error() + " " + rollerr.Error())
				}
				SQL.Mutex.Unlock()
				return err
			}
		}

		//插入数据
		for rows.Next() {
			row1, err := rows.Columns()
			if err != nil {
				rollerr := tx.Rollback() // 回滚事务
				if rollerr != nil {
					SQL.Mutex.Unlock()
					return errors.New(err.Error() + " " + rollerr.Error())
				}
				SQL.Mutex.Unlock()
				return err
			}
			for {
				if len(row1) != len(sql_field) {
					row1 = append(row1, "")
				} else {
					break
				}
			}
			sql := Lib.GenerateInsertSQL(key, sql_field, sql_type, row1)
			_, err = tx.Exec(sql)
			if err != nil {
				rollerr := tx.Rollback() // 回滚事务
				if rollerr != nil {
					SQL.Mutex.Unlock()
					return errors.New(err.Error() + " " + rollerr.Error())
				}
				SQL.Mutex.Unlock()
				return err
			}
		}
		// 提交事务
		err = tx.Commit()
		if err != nil {
			rollerr := tx.Rollback() // 回滚事务
			if rollerr != nil {
				SQL.Mutex.Unlock()
				return errors.New(err.Error() + " " + rollerr.Error())
			}
			SQL.Mutex.Unlock()
			return err
		}
		SQL.Mutex.Unlock()

	}

	return nil
}
