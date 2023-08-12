package raptor

import (
	"errors"
	"reflect"
	"time"
)

// Record represents a Column/Value map for a row in the database.
type Record map[string]any

func (r Record) GetString(col string) string {
	return GetRecordValueLossy[string](r, col)
}

func (r Record) GetTime(col string) time.Time {
	return GetRecordValueLossy[time.Time](r, col)
}

func (r Record) GetBool(col string) bool {
	return GetRecordValueLossy[bool](r, col)
}

func (r Record) GetInt(col string) int64 {
	return GetRecordValueLossy[int64](r, col)
}

func GetRecordValue[T any](record Record, key string) (val T, found bool) {
	untypedVal, ok := record[key]
	if !ok {
		return
	}
	val, found = untypedVal.(T)
	return
}

func GetRecordValueLossy[T any](record Record, key string) (val T) {
	val, _ = GetRecordValue[T](record, key)
	return
}

func ScanRecord(s Scanner) (Record, error) {
	column, err := s.Columns()
	if err != nil {
		return nil, err
	}
	valPtr := make([]any, len(column))

	for i, v := range make([]any, len(column)) {
		v := v
		valPtr[i] = &v
	}

	err = s.Scan(valPtr...)
	if err != nil {
		return nil, err
	}

	r := make(Record, len(column))
	for i, v := range column {
		r[v] = *valPtr[i].(*any)
	}

	return r, nil
}

func ScanAllRecord(r *Rows) ([]Record, error) {
	var records []Record

	for r.Next() {
		r, err := ScanRecord(r)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, r.Close()
}

var (
	ErrRequirePointer = errors.New("raptor: unmarshal destination must be a pointer")
	ErrRequireStruct  = errors.New("raptor: unmarshal destination must be a struct")
)

func UnmarshalRow(s Scanner, dest any) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrRequirePointer
	}
	if rv.Elem().Kind() != reflect.Struct {
		return ErrRequireStruct
	}

	columns, err := s.Columns()
	if err != nil {
		return err
	}

	destE := rv.Elem()
	destT := destE.Type()

	fMap := make(map[string]reflect.Value, len(columns))

FieldLoop:
	for fi := 0; fi < destT.NumField(); fi++ {
		f := destT.Field(fi)

		name := f.Name
		switch f.Tag.Get("db") {
		case "-":
			continue FieldLoop
		case "":
			// no-op
		default:
			name = f.Tag.Get("db")
		}

		fMap[name] = destE.FieldByName(f.Name)
	}

	valPtr := make([]any, len(columns))

	for i, col := range columns {
		if f, ok := fMap[col]; ok {
			if f.IsValid() && f.CanSet() {
				valPtr[i] = f.Addr().Interface()
			} else {
				panic("TODO: Return Error for unable to set struct value")
			}
		} else {
			var v any
			valPtr[i] = &v
		}
	}

	return s.Scan(valPtr...)
}

// MarshalObject converts the given Struct into a Database Record map.
//
// It used the `db` tag to map struct fields names to database column names.
func MarshalObject(obj any) (Record, error) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, ErrRequireStruct
	}

	srcType := rv.Type()

	rec := make(Record)

FieldLoop:
	for fi := 0; fi < srcType.NumField(); fi++ {
		f := srcType.Field(fi)

		name := f.Name
		switch f.Tag.Get("db") {
		case "-":
			continue FieldLoop
		case "":
			// no-op
		default:
			name = f.Tag.Get("db")
		}

		rec[name] = rv.FieldByName(f.Name).Interface()
	}

	return rec, nil
}

type RecordMarshaler interface {
	MarshalRecord() (Record, error)
}

type RecordUnmarshaler interface {
	UnmarshalRecord(Record) error
}
