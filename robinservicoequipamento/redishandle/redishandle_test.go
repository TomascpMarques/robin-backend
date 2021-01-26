package redishandle_test

import (
	"go-graphql-equipamento/redishandle"
	"testing"
)

//TesteExtrairIDMaisRecente -
func TestExtrairIDMaisRecente(t *testing.T) {
	t.Run("Deve extrair o id mais recente ('Item5'>'5')", func(t *testing.T) {
		input := []string{
			"Item1",
			"Item2",
			"Item3",
			"Item4",
			"Item5",
		}

		expected := "5"

		actual := redishandle.ExtrairIDMaisRecente(&input)
		if actual != expected {
			t.Fail()
		}
	})
}
