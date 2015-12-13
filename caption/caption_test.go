package caption

import (
	"testing"

	"github.com/martinlindhe/go-subber/testExtras"
	"github.com/stretchr/testify/assert"
)

func TestRenderTime(t *testing.T) {

	cap := Caption{
		Seq:   1,
		Text:  []string{"<i>Go ninja!</i>"},
		Start: testExtras.MakeTime(18, 40, 22, 110),
		End:   testExtras.MakeTime(18, 41, 20, 123)}

	assert.Equal(t, "18:40:22,110 --> 18:41:20,123", cap.SrtTime())
}