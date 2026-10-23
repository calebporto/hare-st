# hare-st

Banco de dados em memória (estilo Redis), escrito em Go, com protocolo de
texto próprio via TCP. Projeto de estudo focado em redes, concorrência e
arquitetura em camadas (server → controller → aplicação).

## Requisitos

- Go 1.27+

## Executando o servidor

```bash
cd src
go run ./cmd/server
```

O servidor escuta por padrão na porta `8379`.

## Conectando

Qualquer cliente TCP de texto funciona, por exemplo `nc` (netcat) ou
`telnet`:

```bash
nc localhost 8379
```

Cada comando é enviado como uma linha de texto (terminada em `\n`), usando
espaço como delimitador entre os campos.

## Tipos de dados suportados

| Tipo      | Descrição                                   | Limite              |
|-----------|----------------------------------------------|----------------------|
| `char`    | Um único caractere                           | 1 byte               |
| `string`  | Texto                                         | até 1024 runes (1KB) |
| `int`     | Inteiro de 64 bits                            | -                     |
| `float`   | Ponto flutuante de 64 bits                    | -                     |
| `boolean` | `true` ou `false` (minúsculo)                  | -                     |

Regras gerais:
- Chave: até 256 caracteres, aceita letras, números, `_`, `:`, `-` e `ç`/`Ç`.
- Valor: até 1KB por objeto.
- Não há persistência: ao reiniciar o servidor, todos os dados são perdidos.

## Comandos

### `put` — grava uma chave

```
put {chave} {tipo} {valor}
```

Resposta:

```
status success
```

ou, em caso de erro, um JSON descrevendo o problema:

```json
{"error":"invalid_request","data":"...","message":"..."}
```

Exemplos:

```
put nome string Caleb
put idade int 30
put altura float 1.78
put ativo boolean true
put inicial char C
```

### `get` — lê uma chave

```
get {chave}
```

Resposta:

```
data: {tipo} {valor}
```

ou `null` se a chave não existir.

Exemplo:

```
get nome
```

### `clean` — remove uma chave

```
clean {chave}
```

Resposta:

```
status success
```

## Exemplo de sessão completa

```bash
$ nc localhost 8379
put user:1 string "Caleb Porto"
status success
get user:1
data: string "Caleb Porto"
clean user:1
status success
get user:1
null
```

## Mais detalhes

Para o detalhamento completo das regras de negócio e da arquitetura interna
(camadas server / controller / parser / aplicação), veja o
[`PROJECT.md`](./PROJECT.md).
