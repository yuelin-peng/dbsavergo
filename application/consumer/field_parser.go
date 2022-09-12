package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-collections/collections/stack"
)

type FieldParser struct {
	operatorList []Calculator
}

type Calculator interface {
	Calc(v map[string]interface{}, operandList *stack.Stack) (interface{}, error)
	String() string
}

const (
	FuncUnix            = "Unix"
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

func createCalculator(config string) (Calculator, error) {
	switch config {
	case FuncUnix:
		return &UnixFunc{}, nil
	default:
		return NewFieldValueFunc(config)
	}
}

type fieldValueFunc struct {
	fieldName string
}

func NewFieldValueFunc(fieldName string) (*fieldValueFunc, error) {
	if len(fieldName) == 0 {
		return nil, fmt.Errorf("field name can't be empty")
	}
	return &fieldValueFunc{
		fieldName: fieldName,
	}, nil
}

func (f *fieldValueFunc) Calc(v map[string]interface{}, operandList *stack.Stack) (interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("value is nil")
	}
	if result, ok := v[f.fieldName]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("missing field=%v, v=%v", f.fieldName, v)
}

func (f *fieldValueFunc) String() string {
	return f.fieldName
}

type UnixFunc struct {
}

func (f *UnixFunc) Calc(v map[string]interface{}, operandList *stack.Stack) (interface{}, error) {
	value := operandList.Pop()
	if value == nil {
		return nil, fmt.Errorf("%s missing operand", f.String())
	}
	switch t := value.(type) {
	case time.Time:
		return t.Unix(), nil
	case string:
		tmp, err := time.ParseInLocation(YYYY_MM_DD_HH_MM_SS, t, time.Local)
		if err != nil {
			return nil, fmt.Errorf("[%s] time is invalid, value=%v, err=%v", f.String(), t, err)
		}
		return tmp.Unix(), nil
	default:
		return nil, fmt.Errorf("[%s] operand type is invalid, value=%T,%v",
			f.String(), t, t)
	}
}

func (f *UnixFunc) String() string {
	return "Unix"
}

func NewFieldParser(fieldConfig string) (*FieldParser, error) {
	operandList := strings.Split(fieldConfig, ".")
	operatorList := make([]Calculator, 0, len(operandList))
	for _, operand := range operandList {
		c, err := createCalculator(operand)
		if err != nil {
			return nil, err
		}
		operatorList = append(operatorList, c)
	}
	return &FieldParser{
		operatorList: operatorList,
	}, nil
}

func (p *FieldParser) GetFieldValue(v map[string]interface{}) (interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("input value is nil")
	}
	operandList := stack.New()
	for _, operator := range p.operatorList {
		result, err := operator.Calc(v, operandList)
		if err != nil {
			return nil, err
		}
		operandList.Push(result)
	}
	if operandList.Len() != 1 {
		return nil, fmt.Errorf("calculate result size is not 1, result=%v", operandList)
	}
	return operandList.Pop(), nil
}
