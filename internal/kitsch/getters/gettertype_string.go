// Code generated by "stringer -type=GetterType"; DO NOT EDIT.

package getters

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeCustom-0]
	_ = x[TypeFile-1]
	_ = x[TypeAncestorFile-2]
	_ = x[TypeEnv-3]
}

const _GetterType_name = "TypeCustomTypeFileTypeAncestorFileTypeEnv"

var _GetterType_index = [...]uint8{0, 10, 18, 34, 41}

func (i GetterType) String() string {
	if i < 0 || i >= GetterType(len(_GetterType_index)-1) {
		return "GetterType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GetterType_name[_GetterType_index[i]:_GetterType_index[i+1]]
}
