package db

import (
	"fmt"
)

// ошибки уровня сторадж, с которыми будем сравнивать поступающие ошибки в сервис
var (
	ErrDuplicate = fmt.Errorf("duplicate value error")
	ErrNotFound  = fmt.Errorf("not found error")
	ErrModify    = fmt.Errorf("failed to modify")
)
