package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsLocalEnv(t *testing.T) {
	fmt.Println(os.Getenv("GIT_TERMINAL_PROMPT"))
	fmt.Println(IsLocalEnv())
	assert.Equal(t, true, IsLocalEnv())
}
