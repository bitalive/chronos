package mem

import (
	"unsafe"
)

// StringToBytes converts a string to a byte slice WITHOUT allocation.
// The caller must NOT modify the resulting byte slice as strings are immutable.
// This is safe to use when the resulting slice is used for read-only operations
// or temporary operations like map lookups.
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts a byte slice to a string WITHOUT allocation.
// The caller must ensure that the underlying byte slice is NOT modified
// while the string is in use, as strings are expected to be immutable.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// UnsafeString returns a string from a byte slice. This is an alias for BytesToString.
func UnsafeString(b []byte) string {
	return BytesToString(b)
}

// UnsafeBytes returns a byte slice from a string. This is an alias for StringToBytes.
func UnsafeBytes(s string) []byte {
	return StringToBytes(s)
}
