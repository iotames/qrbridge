package dbtable

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type GetIDer interface {
	ParseID() ID
	GetID() int64
}

type IModel interface {
	GetIDer
	ToMap(m IModel) map[string]interface{}
}

type BaseModel struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b BaseModel) ParseID() ID {
	return ParseInt64(b.ID)
}

func (b BaseModel) GetID() int64 {
	return b.ID
}

func (b BaseModel) ToMap(m IModel) map[string]interface{} {
	typeof := reflect.TypeOf(m).Elem()
	typevalue := reflect.ValueOf(m).Elem()
	fieldLen := typeof.NumField()
	fieldsMap := make(map[string]interface{}, fieldLen+2)
	for i := 0; i < fieldLen; i++ {
		field := typeof.Field(i)
		fvalue := typevalue.Field(i)
		value := fvalue.Interface()
		if field.Name == "BaseModel" {
			for j := 0; j < field.Type.NumField(); j++ {
				fieldj := field.Type.Field(j)
				fvaluej := fvalue.Field(j)
				valuej := fvaluej.Interface()
				if fieldj.Name == "ID" {
					valuej = fmt.Sprintf("%d", valuej)
				}
				fieldsMap[fieldj.Name] = valuej
			}
		} else {
			if strings.Contains(field.Name, "ID") {
				value = fmt.Sprintf("%d", value)
			}
			fieldsMap[field.Name] = value
		}
	}
	return fieldsMap
}
