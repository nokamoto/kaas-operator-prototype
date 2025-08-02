package namegen

import (
	"fmt"

	"github.com/google/uuid"
)

type Namegen struct{}

// New generates a unique name for a resource.
// It formats the name using the provided format string and appends a random suffix.
func (*Namegen) New(format string, v ...any) string {
	return fmt.Sprintf("%s-%s", fmt.Sprintf(format, v...), uuid.NewString())
}
