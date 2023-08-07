package db

import (
	"fmt"
)

var (
	ErrDuplicate = fmt.Errorf("duplicate value error")
	ErrNotFound  = fmt.Errorf("not found error")
	ErrModify    = fmt.Errorf("failed to modify")
)
