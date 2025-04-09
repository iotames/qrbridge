package db

import (
	"fmt"
)

// ExecSqlText 执行SQL文本
func ExecSqlText(sqlText string) error {
	// 执行SQL文本
	_, err := GetDbOpen().Exec(sqlText)
	return err
}

// ExecSqlWithTransaction 在事务中执行多条SQL语句
func ExecSqlWithTransaction(sqlStatements []string) error {
	// 开始事务
	tx, err := GetDbOpen().Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 确保函数结束时要么提交要么回滚事务
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	// 执行每条SQL语句
	for _, sql := range sqlStatements {
		_, err := tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("执行SQL失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}
