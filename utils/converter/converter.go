package converter

import (
	"database/sql"
	"time"

	"github.com/bignyap/go-admin/utils/misc"
)

func StrTo[T any](str string, converter Converter[T]) (T, error) {
	return ConvertString(str, converter)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//							Str to DataTypes -- Helpful with API form validation

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func StrToInt(str string) (int, error) {
	return StrTo[int](str, IntConverter{})
}

func StrToFloat(str string) (float64, error) {
	return StrTo[float64](str, FloatConverter{})
}

func StrToBool(str string) (bool, error) {
	return StrTo[bool](str, BoolConverter{})
}

func StrToDate(str string) (time.Time, error) {
	return StrTo[time.Time](str, DateConverter{})
}

func StrToUnixTime(str string) (int, error) {
	return StrTo[int](str, UnixTimeConverter{})
}

func StrToNullInt32(str string) (sql.NullInt32, error) {
	return StrTo[sql.NullInt32](str, NullInt32Converter{})
}

func StrToNullInt64(str string) (sql.NullInt64, error) {
	return StrTo[sql.NullInt64](str, NullInt64Converter{})
}

func StrToNullFloat64(str string) (sql.NullFloat64, error) {
	return StrTo[sql.NullFloat64](str, NullFloat64Converter{})
}

func StrToNullBool(str string) (sql.NullBool, error) {
	return StrTo[sql.NullBool](str, NullBoolConverter{})
}

func StrToNullStr(str string) sql.NullString {
	res, _ := StrTo[sql.NullString](str, NullStringConverter{})
	return res
}

func StrToNullTime(str string) (sql.NullTime, error) {
	return StrTo[sql.NullTime](str, NullTimeConverter{})
}

func StrToUnixNullTime(str string) (sql.NullInt64, error) {
	return StrTo[sql.NullInt64](str, NullUnixTime64Converter{})
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//						sql.Null to Type Specific Converter -- Helpful in SQL to Go Type

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type Nullable interface {
	IsValid() bool
}

type MyNullString struct {
	sql.NullString
}

type MyNullBool struct {
	sql.NullBool
}

type MyNullInt32 struct {
	sql.NullInt32
}

type MyNullInt64 struct {
	sql.NullInt64
}

type MyNullTime struct {
	sql.NullTime
}

func (n MyNullString) IsValid() bool {
	return n.Valid
}

func (n MyNullBool) IsValid() bool {
	return n.Valid
}

func (n MyNullInt32) IsValid() bool {
	return n.Valid
}

func (n MyNullInt64) IsValid() bool {
	return n.Valid
}

func (n MyNullTime) IsValid() bool {
	return n.Valid
}

func NullToPointer[T Nullable, P any](input *T, converter func(*T) P) *P {
	if input == nil || !(*input).IsValid() {
		return nil
	}
	val := converter(input)
	return &val
}

func FromNullString(input *MyNullString) string {
	return input.String
}

func FromNullBool(input *MyNullBool) bool {
	return input.Bool
}

func FromNullInt32(input *MyNullInt32) int {
	return int(input.Int32)
}

func FromNullInt64(input *MyNullInt64) int {
	return int(input.Int64)
}

func FromNullTime(input *MyNullTime) time.Time {
	return time.Time(input.Time)
}

func FromNullInt32ToTime(input *MyNullInt32) time.Time {
	return misc.FromUnixTime32(input.Int32)
}

func NullStrToStr(input *sql.NullString) *string {
	val := MyNullString{NullString: *input}
	return NullToPointer(&val, FromNullString)
}

func NullBoolToBool(input *sql.NullBool) *bool {
	val := MyNullBool{NullBool: *input}
	return NullToPointer(&val, FromNullBool)
}

func NullInt32ToInt(input *sql.NullInt32) *int {
	val := MyNullInt32{NullInt32: *input}
	return NullToPointer(&val, FromNullInt32)
}

func NullInt64ToInt(input *sql.NullInt64) *int {
	val := MyNullInt64{NullInt64: *input}
	return NullToPointer(&val, FromNullInt64)
}

func NullTimeToTime(input *sql.NullTime) *time.Time {
	val := MyNullTime{NullTime: *input}
	return NullToPointer(&val, FromNullTime)
}

func NullInt32ToTime(input *sql.NullInt32) *time.Time {
	val := MyNullInt32{NullInt32: *input}
	return NullToPointer(&val, FromNullInt32ToTime)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//										Pointer to SQLNull Types

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func IntPtrToNullInt32(ptr *int) sql.NullInt32 {
	if ptr != nil {
		return sql.NullInt32{Int32: int32(*ptr), Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

func TimePtrToUnixInt(t *time.Time) *int {
	if t == nil {
		return nil
	}
	unixTime := int(t.Unix())
	return &unixTime
}

func BoolPtrToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{Valid: false}
}
