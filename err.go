package builder

import (
	"errors"
)

// Error definition
var (
	ErrEmptyCondition = errors.New("empty condition")

	ErrEmptySQLType = errors.New("empty sql type")

	ErrListIsNotEmpty = errors.New("There are some errors in sql, please check your sql")
)
