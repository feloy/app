package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}

func formatWhereClause(where []string) string {
	if len(where) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(where, " AND ")
}

func findMany(ctx context.Context, tx *sqlx.Tx, ss interface{}, query string, args ...interface{}) error {
	rows, err := tx.QueryxContext(ctx, query, args...)

	if err != nil {
		return err
	}

	defer rows.Close()

	sPtrVal, err := asSlicePtrValue(ss) // get the reflect.Value of the ptr to slice

	if err != nil {
		return err
	}

	sVal := sPtrVal.Elem()                           // get the relfect.Value of the slice pointed to by ss
	newSlice := reflect.MakeSlice(sVal.Type(), 0, 0) // new slice
	elemType := sliceElemType(sVal)                  // get the slice element's type

	for rows.Next() {
		newVal := reflect.New(elemType) // create a new value of this type
		if err := rows.StructScan(newVal.Interface()); err != nil {
			return nil
		}
		newSlice = reflect.Append(newSlice, newVal)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	sPtrVal.Elem().Set(newSlice) // change the value pointed to be the ptr to slice to our new slice

	return nil
}

// func findOne(ctx context.Context, tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
// 	s, err := asStructPtr(dest)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return tx.QueryRowxContext(ctx, query, args...).StructScan(s.Interface())
// }

// sliceElemType takes a reflect.Value which is a ptr to slice or a slice,
// and returns the reflect.Type of the elements the slice holds.
// If the slice holds a pointer type, it returns the type pointed to.
func sliceElemType(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	vv := v.Type().Elem() // get the reflect.Type of the elements of the slice

	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem() // if it is a pointer, get the type it points to
	}

	return vv
}

func isSlicePtr(v interface{}) bool {
	typ := reflect.TypeOf(v)

	return typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice
}

func asSlicePtrValue(v interface{}) (reflect.Value, error) {
	if !isSlicePtr(v) {
		return reflect.Value{}, errors.New("expecting a pointer to slice")
	}
	return reflect.ValueOf(v), nil
}

//
// func asStructPtr(v interface{}) (reflect.Value, error) {
// 	vType := reflect.TypeOf(v)

// 	if vType.Kind() != reflect.Ptr && vType.Elem().Kind() != reflect.Struct {
// 		return reflect.Value{}, errors.New("expecting a pointer to struct type")
// 	}

// 	return reflect.ValueOf(v), nil
// }

func execQuery(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) error {
	_, err := tx.ExecContext(ctx, query, args...)

	return err
}
