package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EMmails(t *testing.T) {
	t.Parallel()

	assert.Equal(t, false, ValidateEmail("dgmccallum.gmail.com"), "Should be rejected for no @")
	assert.Equal(t, true, ValidateEmail("dgmccallum@gmail.com"), "Should be allowed")
	assert.Equal(t, false, ValidateEmail("dgmccallum@gmail.commm"), "Should be rejected for too long tld")
	assert.Equal(t, false, ValidateEmail("dgmccallum@gmail"), "Should be rejected for not enough domains")
	assert.Equal(t, true, ValidateEmail("allen.paton@ggc.scot.nhs.co.uk"), "Should be allowed for many domains")
}
