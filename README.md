# cached-translation-middleware

[![License](https://img.shields.io/github/license/caiocampos/go.db.restapi.svg)](LICENSE)

Middleware para o serviço de tradução da página de portfólio desenvolvido em Go com estratégia Cache-Aside e Background Refresh, utilizando Gin Gonic e Redis

## Executando:

Para executar o projeto é necessário o Go instalado e configurado, siga as instruções do site a seguir para configurar:

http://www.golangbr.org/doc/instalacao

Antes de executar modifique o arquivo config.toml para apontar para o MongoDB instalado.

Após instalar o Go e configurar o arquivo config.toml compile o código, utilize o seguinte comando para isso:

> go build cmd/main.go

E depois, para executar:

> ./main

Para configurar o Redis localmente siga o tutorial:

https://medium.com/@toticavalcanti/redis-no-windows-10-26fcbfce78ae
