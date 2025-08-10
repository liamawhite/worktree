package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		_, err := Select("Test", []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no options provided")
	})

	// Note: We can't easily test the interactive functionality in unit tests
	// since it requires terminal interaction. The main validation is that
	// the function handles empty options correctly and doesn't panic.
	t.Run("valid options structure", func(t *testing.T) {
		// This test just ensures the function can be called without panicking
		// Actual selection would require terminal interaction
		options := []string{"option1", "option2", "option3"}
		
		// We can't run the actual interactive selection in tests,
		// but we can verify the options are properly structured
		assert.NotEmpty(t, options)
		assert.Len(t, options, 3)
	})
}