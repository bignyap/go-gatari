package formvalidator

import (
	"database/sql"

	"github.com/bignyap/go-admin/utils/converter"
)

type FormField interface {
	Validator(string) (interface{}, error)
}

type IntFormField struct{}

func (IntFormField) Validator(str string) (interface{}, error) {
	return converter.StrToInt(str)
}

type FloatFormField struct{}

func (FloatFormField) Validator(str string) (interface{}, error) {
	return converter.StrToFloat(str)
}

type BoolFormField struct{}

func (BoolFormField) Validator(str string) (interface{}, error) {
	return converter.StrToBool(str)
}

type DateFormField struct{}

func (DateFormField) Validator(str string) (interface{}, error) {
	return converter.StrToDate(str)
}

type UnixTimeFormField struct{}

func (UnixTimeFormField) Validator(str string) (interface{}, error) {
	return converter.StrToUnixTime(str)
}

type NullInt32FormField struct{}

func (c NullInt32FormField) Validator(str string) (interface{}, error) {
	return converter.StrToNullInt32(str)
}

type NullInt64FormField struct{}

func (c NullInt64FormField) Validator(str string) (interface{}, error) {
	return converter.StrToNullInt64(str)
}

type NullFloat64FormField struct{}

func (c NullFloat64FormField) Validator(str string) (interface{}, error) {
	return converter.StrToNullFloat64(str)
}

type NullBoolFormField struct{}

func (c NullBoolFormField) Validator(str string) (interface{}, error) {
	return converter.StrToNullBool(str)
}

type NullStringFormField struct{}

func (c NullStringFormField) Validator(str string) (interface{}, error) {
	return sql.NullString{String: str, Valid: str != ""}, nil
}

type NullTimeFormField struct{}

func (c NullTimeFormField) Validator(str string) (interface{}, error) {
	return converter.StrToNullTime(str)
}

type NullUnixTimeFormField struct{}

func (c NullUnixTimeFormField) Validator(str string) (interface{}, error) {
	return converter.StrToUnixNullTime(str)
}
