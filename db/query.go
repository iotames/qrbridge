package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

// GetOneData 根据where条件查询单条数据，支持结构体指针或map接收结果
// querySQL SQL查询语句 例：select field1, field2 from table1 where name = $1 and status = $2
// dest: 用于接收结果的结构体指针或map[string]any
// args: SQL参数
func GetOneData(querySQL string, dest interface{}, args ...interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("dest必须是有效的非空指针")
	}

	stmt, err := GetDbOpen().Prepare(querySQL)
	if err != nil {
		return fmt.Errorf("预处理SQL语句失败: %v", err)
	}
	defer stmt.Close()

	// 改用Query获取sql.Rows（即使只查一行）
	rows, err := stmt.Query(args...)
	if err != nil {
		return fmt.Errorf("查询失败: %v", err)
	}
	defer rows.Close()

	// 直接读取首行（模拟QueryRow行为）
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return fmt.Errorf("查询错误: %v", err)
		}
		return nil // 无数据
	}

	switch d := dest.(type) {
	case *map[string]any:
		return scanRowToMap(rows, d)
	default:
		if val.Elem().Kind() == reflect.Struct {
			return scanRowToStruct(rows, dest)
		}
		return fmt.Errorf("不支持的dest类型(%T)", dest)
	}
}

// 动态扫描到map（使用sql.Rows）
func scanRowToMap(rows *sql.Rows, dest *map[string]any) error {
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取列失败: %v", err)
	}

	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(interface{})
	}

	if err := rows.Scan(values...); err != nil {
		return fmt.Errorf("扫描失败: %v", err)
	}

	result := make(map[string]any)
	for i, col := range cols {
		result[col] = *(values[i].(*interface{}))
	}
	*dest = result
	return nil
}

// 扫描到结构体（通用逻辑）
func scanRowToStruct(rows *sql.Rows, dest interface{}) error {
	// 使用rows.Columns()验证列与结构体标签匹配
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取列失败: %v", err)
	}

	destVal := reflect.ValueOf(dest).Elem()
	fields := make([]interface{}, len(cols))

	// 按列名映射结构体字段
	for i, col := range cols {
		fieldFound := false
		for j := 0; j < destVal.NumField(); j++ {
			tag := destVal.Type().Field(j).Tag.Get("db")
			if tag == col {
				fields[i] = destVal.Field(j).Addr().Interface()
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			return fmt.Errorf("列 %s 无对应的结构体字段", col)
		}
	}

	return rows.Scan(fields...)
}

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
