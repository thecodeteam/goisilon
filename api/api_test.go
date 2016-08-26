package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAuthHeaderVal(t *testing.T) {
	assert.Equal(
		t, "Basic YWt1dHo6cGFzc3dvcmQ=",
		fmtAuthHeaderVal("akutz", "password"))
}

func BenchmarkFormatAuthHeaderVal(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if "Basic YWt1dHo6cGFzc3dvcmQ=" !=
				fmtAuthHeaderVal("akutz", "password") {
				b.FailNow()
			}
		}
	})
}
