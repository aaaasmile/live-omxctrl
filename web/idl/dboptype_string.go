// Code generated by "stringer -type DbOpType"; DO NOT EDIT.

package idl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DbOpHistoryInsert-0]
}

const _DbOpType_name = "DbOpHistoryInsert"

var _DbOpType_index = [...]uint8{0, 17}

func (i DbOpType) String() string {
	if i < 0 || i >= DbOpType(len(_DbOpType_index)-1) {
		return "DbOpType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DbOpType_name[_DbOpType_index[i]:_DbOpType_index[i+1]]
}
