package converter

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/bignyap/go-admin/utils/misc"
)

type Converter[T any] interface {
	Convert(str string) (T, error)
}

type IntConverter struct{}

func (c IntConverter) Convert(str string) (int, error) {
	if str == "" {
		return 0, fmt.Errorf("empty string cannot be converted to int")
	}
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("not a valid integer: %v", err)
	}
	return intValue, nil
}

type FloatConverter struct{}

func (c FloatConverter) Convert(str string) (float64, error) {
	if str == "" {
		return 0.0, fmt.Errorf("empty string cannot be converted to float")
	}
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0, fmt.Errorf("not a valid float: %v", err)
	}
	return floatValue, nil
}

type BoolConverter struct{}

func (c BoolConverter) Convert(str string) (bool, error) {
	if str == "" {
		return false, fmt.Errorf("empty string cannot be converted to boolean")
	}
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return false, fmt.Errorf("not a valid boolean: %v", err)
	}
	return boolVal, nil
}

type DateConverter struct{}

func (c DateConverter) Convert(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, fmt.Errorf("empty string cannot be converted to date")
	}
	dateValue, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, fmt.Errorf("not a valid date: %v", err)
	}
	return dateValue, nil
}

type UnixTimeConverter struct{}

func (c UnixTimeConverter) Convert(str string) (int, error) {

	dateVal, err := ConvertString(str, DateConverter{})
	if err != nil {
		return -1, fmt.Errorf("not a valid integer: %v", err)
	}

	return int(misc.ToUnixTime(dateVal)), nil
}

type NullInt32Converter struct{}

func (c NullInt32Converter) Convert(str string) (sql.NullInt32, error) {
	if str == "" {
		return sql.NullInt32{Valid: false}, nil
	}
	intValue, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return sql.NullInt32{}, fmt.Errorf("not a valid integer: %v", err)
	}
	return sql.NullInt32{Int32: int32(intValue), Valid: true}, nil
}

type NullInt64Converter struct{}

func (c NullInt64Converter) Convert(str string) (sql.NullInt64, error) {
	if str == "" {
		return sql.NullInt64{Valid: false}, nil
	}
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return sql.NullInt64{}, fmt.Errorf("not a valid integer: %v", err)
	}
	return sql.NullInt64{Int64: intValue, Valid: true}, nil
}

type NullFloat64Converter struct{}

func (c NullFloat64Converter) Convert(str string) (sql.NullFloat64, error) {
	if str == "" {
		return sql.NullFloat64{Valid: false}, nil
	}
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return sql.NullFloat64{}, fmt.Errorf("not a valid float: %v", err)
	}
	return sql.NullFloat64{Float64: floatValue, Valid: true}, nil
}

type NullBoolConverter struct{}

func (c NullBoolConverter) Convert(str string) (sql.NullBool, error) {
	if str == "" {
		return sql.NullBool{Valid: false}, nil
	}
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return sql.NullBool{}, fmt.Errorf("not a valid boolean: %v", err)
	}
	return sql.NullBool{Bool: boolVal, Valid: true}, nil
}

type NullStringConverter struct{}

func (c NullStringConverter) Convert(str string) (sql.NullString, error) {
	return sql.NullString{String: str, Valid: str != ""}, nil
}

type NullTimeConverter struct{}

func (c NullTimeConverter) Convert(str string) (sql.NullTime, error) {
	if str == "" {
		return sql.NullTime{Valid: false}, nil
	}
	timeValue, err := time.Parse("2006-01-02", str)
	if err != nil {
		return sql.NullTime{}, fmt.Errorf("not a valid date: %v", err)
	}
	return sql.NullTime{Time: timeValue, Valid: true}, nil
}

type NullUnixTime64Converter struct{}

func (c NullUnixTime64Converter) Convert(str string) (sql.NullInt64, error) {

	dateVal, err := ConvertString(str, NullTimeConverter{})
	if err != nil {
		return sql.NullInt64{}, fmt.Errorf("not a valid integer: %v", err)
	}

	if !dateVal.Valid {
		return sql.NullInt64{Valid: false}, nil
	}

	return sql.NullInt64{Int64: misc.ToUnixTime(dateVal.Time), Valid: true}, nil
}

func ConvertString[T any](str string, converter Converter[T]) (T, error) {
	return converter.Convert(str)
}
