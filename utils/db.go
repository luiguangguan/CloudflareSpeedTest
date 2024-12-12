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

const createPasswordsTableSQL = `
CREATE TABLE IF NOT EXISTS Passwords (
	"pwd" TEXT(64) NOT NULL
);`

const createTableIpTraceInfosSQL = `
CREATE TABLE IF NOT EXISTS "IpTraceInfos" (
  "IP" TEXT(64) NOT NULL,
  "TraceInfo" TEXT(10000) NOT NULL,
  "Remark" TEXT(100),
  "CreateTime"  TEXT NOT NULL
);`

const creteTraceInfoUniqueView = `DROP VIEW IF EXISTS TraceInfoUnique;
CREATE VIEW "TraceInfoUnique" AS select IP,MAX(TraceInfo) TraceInfo from IpTraceInfos Group by IP;`

const creteMaxSpeedView = `
DROP VIEW IF EXISTS MaxSpeed;
CREATE VIEW MaxSpeed AS
SELECT 
    A.IP,
    A.Port,
    MAX(A.DownloadSpeed) AS MaxDownloadSpeed,
    MIN(A.DownloadSpeed) AS MinDownloadSpeed,
    AVG(A.DownloadSpeed) AS AvgDownloadSpeed,
		ROUND(MIN(A.Delay),2) AS MinDelay,
    ROUND(MAX(A.Delay),2) AS MaxDelay,
    ROUND(AVG(A.Delay),2) AS AvgDelay,
    SUM(A.LossRate) AS SumLossRate,
    AVG(A.LossRate) AS AVGLossRate,
    DATE(A.CreateTime) AS Date,
    COUNT(1) AS Count,
    A.Remark,
		B.TraceInfo
FROM 
    speedTestResult as A
		left join TraceInfoUnique As B
		on A.IP=B.IP
GROUP BY 
    A.IP, A.Port, DATE(A.CreateTime), A.Remark ORDER BY  MIN(A.DownloadSpeed) desc ,MAX(A.Delay) asc,MAX(A.LossRate) asc;`

const creteRecordView = `
DROP VIEW IF EXISTS Record;
CREATE VIEW Record AS
select IP||'#'||Port||'#'||Remark Record,
IP,
Port,
MAX(DownloadSpeed)MaxDownloadSpeed,
MIN(DownloadSpeed)MinDownloadSpeed,
ROUND(MIN(Delay),2) AS MinDelay,
ROUND(MAX(Delay),2) AS MaxDelay,
ROUND(AVG(Delay),2) AS AvgDelay,
SUM(LossRate)SumLossRate,
AVG(LossRate)AVGLossRate,
SUBSTR(CreateTime ,1,10) Date,
count(1)Count,
Remark from speedTestResult 
group by IP,
port,
SUBSTR(CreateTime,1,10),
Remark 
ORDER  by count(1) desc,
max(DownloadSpeed) desc`

const speedTestWithTrace = `DROP VIEW IF EXISTS speedTestWithTrace;
CREATE VIEW "speedTestWithTrace" AS SELECT
    s."IP",
    s."Port",
    s."Sended",
    s."Received",
    s."Delay",
    s."LossRate",
    s."DownloadSpeed",
    s."CreateTime",
    s."Remark" AS speedTestRemark,
    i."TraceInfo",
    i."Remark" AS traceInfoRemark
FROM
    speedTestResult s
LEFT JOIN
    IpTraceInfos i
ON
    s."IP" = i."IP"`

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

		// result, err := dbInstance.Exec("PRAGMA journal_mode=WAL;")
		// _, err := result.RowsAffected()

		// 检查连接是否可用
		if err = dbInstance.Ping(); err != nil {
			fmt.Println("数据库连接测试失败:", err)
			dbInstance = nil
		}

		// 如果是首次运行，初始化数据
		if _, err = dbInstance.Exec(createTableSQL); err != nil {
			fmt.Println("初始化表speedTestResult结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		// 如果是首次运行，初始化数据
		if _, err = dbInstance.Exec(createPasswordsTableSQL); err != nil {
			fmt.Println("初始化表Passwords结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		if _, err = dbInstance.Exec(createTableIpTraceInfosSQL); err != nil {
			fmt.Println("初始化表createTableIpTraceInfosSQL结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		if _, err = dbInstance.Exec(creteTraceInfoUniqueView); err != nil {
			fmt.Println("初始化视图TraceInfoUnique结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		if _, err = dbInstance.Exec(creteMaxSpeedView); err != nil {
			fmt.Println("初始化视图MaxSpeed结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		if _, err = dbInstance.Exec(creteRecordView); err != nil {
			fmt.Println("初始化视图Record结构失败:", err)
			dbInstance.Close()
			dbInstance = nil
			return
		}
		if _, err = dbInstance.Exec(speedTestWithTrace); err != nil {
			fmt.Println("初始化视图speedTestWithTrace结构失败:", err)
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

// 查询一行一列的结果，返回一个值
func Scalar(query string, args ...interface{}) (interface{}, error) {
	db, err := GetDBInstance()
	if err != nil {
		return nil, err
	}

	// 执行查询，获取结果
	row := db.QueryRow(query, args...)

	// 通过 Scan 将查询结果扫描到一个变量中
	var result interface{}
	if err := row.Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			// 如果没有返回行，返回 nil 或特定的错误
			return nil, nil
		}
		return nil, err
	}

	// 返回查询的结果
	return result, nil
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
