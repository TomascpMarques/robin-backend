package redishandle_test

import (
	"fmt"
	"testing"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/redishandle"
)

//TesteExtrairIDMaisRecente -
func TestExtrairIDMaisRecente(t *testing.T) {
	t.Run("Deve extrair o id mais recente ('ItemXX'>'XX')", func(t *testing.T) {
		inputLen := 100000
		input := make([]string, inputLen)
		for k := range input {
			input[k] = "Item" + fmt.Sprint(k)
		}
		expected := fmt.Sprint(inputLen - 1)

		actual := redishandle.ExtrairIDMaisRecente(&input)
		if actual != expected {
			t.Fail()
		}
	})
}
