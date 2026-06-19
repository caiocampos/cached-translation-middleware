package util_test

import (
	"testing"

	"cached-translation-middleware/internal/util"
)

func TestGetTextWithoutLinks(t *testing.T) {
	original := "Middleware para o serviço de tradução da página de portfólio desenvolvido em Go com estratégia Cache-Aside e Background Refresh, utilizando Gin Gonic e Redis. Para o projeto da página acesse: https://github.com/caiocampos/caiocampos.github.io"
	expected := "Middleware para o serviço de tradução da página de portfólio desenvolvido em Go com estratégia Cache-Aside e Background Refresh, utilizando Gin Gonic e Redis. Para o projeto da página acesse: "

	result := util.GetTextWithoutLinks(original)

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}
