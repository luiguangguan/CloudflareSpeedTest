package utils

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "modernc.org/sqlite" // 使用 modernc.org/sqlite 替代 go-sqlite3
)

const (
	defaultDbFile = "./data.db"
)

var (
	dbInstance *sql.DB         // 单例数据库实例
	dbOnce     sync.Once       // 确保单例的创建只执行一次
	dbErr      error           // 用于捕获数据库初始化错误
	isFirstRun bool            // 标记是否为首次运行
	DbFile     = defaultDbFile // 数据文件路径
)

// 初始化表结构
const createTableSQL = `
CREATE TABLE IF NOT EXISTS speedTestResult (
	"IP" TEXT(64) NOT NULL,
	"Port" integer NOT NULL,
	"Sended" integer NOT NULL,
	"Received" integer NOT NULL,
	"Delay" integer NOT NULL,
	"LossRate" real NOT NULL,
	"DownloadSpeed" real NOT NULL,
	"CreateTime" TEXT NOT NULL,
	"Remark" TEXT(100)
);`

// GetDBInstance 返回 SQLite 数据库的单例实例
func GetDBInstance() (*sql.DB, error) {
	var err error
	dbOnce.Do(func() {
		// 如果数据库文件不存在，则创建它
		if _, err = os.Stat(DbFile); os.IsNotExist(err) {
			isFirstRun = true
			fmt.Println("数据库文件不存在，正在创建...")
			file, err := os.Create(DbFile)
			if err != nil {
				fmt.Println("无法创建数据库文件:", err)
				return
			}
			file.Close() // 关闭文件指针
		}

		// 打开数据库连接
		dbInstance, err = sql.Open("sqlite", DbFile) // 使用 "sqlite" 替代 "sqlite3"
		if err != nil {
			fmt.Println("无法打开数据库连接:", err)
			return
		}

		// 检查连接是否可用
		if err = dbInstance.Ping(); err != nil {
			fmt.Println("数据库连接测试失败:", err)
			dbInstance = nil
		}

		// 如果是首次运行，初始化数据
		if _, err = dbInstance.Exec(createTableSQL); err != nil {
			fmt.Println("初始化表结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
	})

	return dbInstance, err
}

// 通用的非查询执行方法，用于插入、更新、删除操作
func ExecNonQuery(query string, args ...interface{}) (int64, error) {
	db, err := GetDBInstance()
	if err != nil {
		fmt.Println("获取数据库实例时出错:", err) // 输出获取数据库实例时的错误
		return 0, err
	}

	// 执行 SQL 语句
	result, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println("执行 SQL 语句时出错:", err) // 输出 SQL 执行时的错误
		return 0, err
	}

	// 返回受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("获取受影响行数时出错:", err) // 输出获取受影响行数时的错误
		return 0, err
	}

	return rowsAffected, nil
}

// 查询数据，返回结果集 []map[string]interface{}
func Select(query string, args ...interface{}) ([]map[string]interface{}, error) {
	db, err := GetDBInstance()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 存储结果的切片
	results := []map[string]interface{}{}

	for rows.Next() {
		// 构造一个键值对映射存储一行数据
		rowMap := make(map[string]interface{})
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = new(interface{})
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// 将每列的数据存入 map
		for i, colName := range columns {
			rowMap[colName] = *(columnPointers[i].(*interface{}))
		}
		results = append(results, rowMap)
	}
	return results, nil
}

// PathExists 检查文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}