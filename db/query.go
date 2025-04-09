package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

// GetOne 根据where条件查询单条数据
// querySQL SQL查询语句
// args: SQL参数
// dest: 用于接收结果的结构体指针
func GetOne(querySQL string, dest []interface{}, args ...interface{}) error {
	// 使用预处理语句执行查询，防止SQL注入
	stmt, err := GetDbOpen().Prepare(querySQL)
	if err != nil {
		return fmt.Errorf("预处理SQL语句失败: %v", err)
	}
	defer stmt.Close()

	// 执行预处理查询
	row := stmt.QueryRow(args...)
	if err := row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			// return fmt.Errorf("未找到匹配的数据记录")
			return nil
		}
		return fmt.Errorf("查询数据失败: %v", err)
	}
	return nil
}

// GetMany 根据where条件查询多条数据
// querySQL SQL查询语句
// dest: 用于接收结果的切片指针
// args: SQL参数
func GetMany(querySQL string, dest interface{}, args ...interface{}) error {
	// 使用预处理语句执行查询，防止SQL注入
	stmt, err := GetDbOpen().Prepare(querySQL)
	if err != nil {
		return fmt.Errorf("预处理SQL语句失败: %v", err)
	}
	defer stmt.Close()

	// 执行预处理查询
	rows, err := stmt.Query(args...)
	if err != nil {
		return fmt.Errorf("查询数据失败: %v", err)
	}
	defer rows.Close()

	// 使用sql.Rows.Scan将结果扫描到目标切片
	if err := scanRows(rows, dest); err != nil {
		return fmt.Errorf("扫描数据失败: %v", err)
	}
	return nil
}

// scanRows 将sql.Rows的结果扫描到目标切片中
func scanRows(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("目标参数必须是切片指针")
	}
	sliceVal := v.Elem()
	elemType := sliceVal.Type().Elem()
	for rows.Next() {
		// 创建新的元素
		newElem := reflect.New(elemType).Interface()
		// 扫描当前行到新元素
		if err := rows.Scan(newElem); err != nil {
			return err
		}
		// 将新元素添加到切片
		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(newElem).Elem()))
	}
	return rows.Err()
}
