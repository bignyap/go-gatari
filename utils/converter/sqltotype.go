package converter

import "database/sql"

type NullConverter interface {
	ConvertToPointer() interface{}
}

// NullStringToPointerConverter implements NullConverter for sql.NullString
type NullStringToPointerConverter struct {
	Input *sql.NullString
}

func (c NullStringToPointerConverter) ConvertToPointer() interface{} {
	if !(c.Input.Valid) {
		return nil
	}
	return &c.Input.String
}

// NullBoolToPointerConverter implements NullConverter for sql.NullBool
type NullBoolToPointerConverter struct {
	Input *sql.NullBool
}

func (c NullBoolToPointerConverter) ConvertToPointer() interface{} {
	if !(c.Input.Valid) {
		return nil
	}
	return &c.Input.Bool
}

// NullInt32ToPointerConverter implements NullConverter for sql.NullInt32
type NullInt32ToPointerConverter struct {
	Input *sql.NullInt32
}

func (c NullInt32ToPointerConverter) ConvertToPointer() interface{} {
	if !(c.Input.Valid) {
		return nil
	}
	intValue := int(c.Input.Int32)
	return &intValue
}

// NullInt64ToPointerConverter implements NullConverter for sql.NullInt64
type NullInt64ToPointerConverter struct {
	Input *sql.NullInt64
}

func (c NullInt64ToPointerConverter) ConvertToPointer() interface{} {
	if !(c.Input.Valid) {
		return nil
	}
	intValue := int(c.Input.Int64)
	return &intValue
}

type NullUnixTimeToPointerConverter struct {
	Input *sql.NullInt64
}

func (c NullUnixTimeToPointerConverter) ConvertToPointer() interface{} {
	if !(c.Input.Valid) {
		return nil
	}
	return &c.Input.Int64
}

func ConvertNullToPointer(converter NullConverter) interface{} {
	return converter.ConvertToPointer()
}
