package behold

import (
	"fmt"
)

func ErrNotExported(field string) error {
	return fmt.Errorf("Field %q not exported", field)
}

func ErrNotPositive(field string) error {
	return fmt.Errorf("Field %q must be positive", field)
}

func ErrAlreadySet(field string, value interface{}) error {
	return fmt.Errorf("Field %q already set: %v", field, value)
}
