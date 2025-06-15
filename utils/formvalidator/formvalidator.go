package formvalidator

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func ParseFormData(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("error parsing form data: %s", err)
	}
	return nil
}

func FormValidator(r *http.Request, fields []string, formField FormField) (map[string]interface{}, error) {

	parsed := make(map[string]interface{})

	for _, field := range fields {
		value, err := formField.Validator(r.FormValue(field))
		if err != nil {
			return nil, err
			// fmt.Errorf("%s must be a valid value", field)
		}

		parsed[field] = value
	}

	return parsed, nil
}

func ParseFromForm[T any](r *http.Request, fields []string, converter FormField) (map[string]T, error) {
	parsed, err := FormValidator(r, fields, converter)
	if err != nil {
		return nil, err
	}

	result := make(map[string]T)
	for k, v := range parsed {
		result[k] = v.(T)
	}

	return result, nil
}

func ParseStringFromForm(r *http.Request, fields []string) (map[string]string, error) {

	parsed := make(map[string]string)

	for _, field := range fields {
		if r.FormValue(field) == "" {
			return nil, fmt.Errorf("%s parameter needs to be a valid string", field)
		}
		parsed[field] = r.FormValue(field)
	}

	return parsed, nil
}

func ParseIntFromForm(r *http.Request, fields []string) (map[string]int, error) {
	return ParseFromForm[int](r, fields, IntFormField{})
}

func ParseFloatFromForm(r *http.Request, fields []string) (map[string]float64, error) {
	return ParseFromForm[float64](r, fields, FloatFormField{})
}

func ParseBoolFromForm(r *http.Request, fields []string) (map[string]bool, error) {
	return ParseFromForm[bool](r, fields, BoolFormField{})
}

func ParseDateFormForm(r *http.Request, fields []string) (map[string]time.Time, error) {
	return ParseFromForm[time.Time](r, fields, DateFormField{})
}

func ParseUnixTimeFromForm(r *http.Request, fields []string) (map[string]int, error) {
	return ParseFromForm[int](r, fields, UnixTimeFormField{})
}

func ParseNullInt32FromForm(r *http.Request, fields []string) (map[string]sql.NullInt32, error) {
	return ParseFromForm[sql.NullInt32](r, fields, NullInt32FormField{})
}

func ParseNullInt64FromForm(r *http.Request, fields []string) (map[string]sql.NullInt64, error) {
	return ParseFromForm[sql.NullInt64](r, fields, NullInt64FormField{})
}

func ParseNullFloat64FromForm(r *http.Request, fields []string) (map[string]sql.NullFloat64, error) {
	return ParseFromForm[sql.NullFloat64](r, fields, NullFloat64FormField{})
}

func ParseNullBoolFromForm(r *http.Request, fields []string) (map[string]sql.NullBool, error) {
	return ParseFromForm[sql.NullBool](r, fields, NullBoolFormField{})
}

func ParseNullStringFromForm(r *http.Request, fields []string) (map[string]sql.NullString, error) {
	return ParseFromForm[sql.NullString](r, fields, NullStringFormField{})
}

func ParseNullTimeFromForm(r *http.Request, fields []string) (map[string]sql.NullTime, error) {
	return ParseFromForm[sql.NullTime](r, fields, NullTimeFormField{})
}

func ParseNullUnixTimeFromForm(r *http.Request, fields []string) (map[string]sql.NullInt64, error) {
	return ParseFromForm[sql.NullInt64](r, fields, NullUnixTimeFormField{})
}

func ParseNullUnixTime32FromForm(r *http.Request, fields []string) (map[string]sql.NullInt32, error) {
	parsed, err := ParseFromForm[sql.NullInt64](r, fields, NullUnixTimeFormField{})
	if err != nil {
		return nil, err
	}

	result := make(map[string]sql.NullInt32)
	for k, v := range parsed {
		if !v.Valid {
			result[k] = sql.NullInt32{
				Int32: 0,
				Valid: false,
			}
		} else {
			result[k] = sql.NullInt32{
				Int32: int32(v.Int64),
				Valid: true,
			}
		}
	}

	return result, nil
}
