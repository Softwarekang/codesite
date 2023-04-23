package clickhouse

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	statement Statement
	clickhouse.Conn
	ctx context.Context
}

func (c *Client) Find(dest interface{}) error {
	if err := c.statement.BuildCondition(); err != nil {
		return err
	}
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("find must pass a pointer, not a value, to Select destination")
	}
	if value.IsNil() {
		return fmt.Errorf("find nil pointer passed to Select destination")
	}
	direct := reflect.Indirect(value)
	if direct.Kind() != reflect.Slice {
		return fmt.Errorf("must pass a slice to Select destination")
	}
	if direct.Len() != 0 {
		direct.Set(reflect.MakeSlice(direct.Type(), 0, direct.Cap()))
	}
	var (
		base      = direct.Type().Elem()
		rows, err = c.Query(c.ctx, c.statement.execSQL.String())
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	rowNums := 0
	for rows.Next() {
		elem := reflect.New(base)
		if err := rows.ScanStruct(elem.Interface()); err != nil {
			return err
		}
		direct.Set(reflect.Append(direct, elem.Elem()))
		rowNums++
	}
	if err := rows.Close(); err != nil {
		return err
	}
	fmt.Printf("[rows:%d] %s\n", rowNums, c.statement.execSQL.String())
	return rows.Err()
}
