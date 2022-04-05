package eco

import "errors"

var (
	ErrRequiresNonNilPtr = errors.New("Unmarshal requires non-nil pointer")
)
