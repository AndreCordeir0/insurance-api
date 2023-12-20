package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

func AbstractInsert[T comparable](table string, columns []string, args []T, transaction *sql.Tx) (int64, error) {
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(table)
	sb.WriteString("(")
	sb.WriteString(strings.Join(columns, ","))
	sb.WriteString(")")
	sb.WriteString("VALUE")
	sb.WriteString("(")
	sb.WriteString(ReturnInterrogations(args))
	sb.WriteString(")")
	fmt.Println("query string: ", sb.String())
	fmt.Println("args: ", args)
	var interfaceArgs []interface{}
	for _, v := range args {
		interfaceArgs = append(interfaceArgs, v)
	}
	result, err := transaction.Exec(sb.String(), interfaceArgs...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ReturnInterrogations[T comparable](args []T) string {
	size := len(args)
	stringArray := make([]string, size)
	for i := 0; i < size; i++ {
		stringArray[i] = "?"
	}
	return strings.Join(stringArray, ",")
}
