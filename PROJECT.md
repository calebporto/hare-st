# Projeto: Banco de Dados em Memória (estilo Redis)

## Propósito

Este é um projeto de estudo em Go com o objetivo de aprender conceitos de
redes, concorrência e arquitetura de sistemas, construindo um banco de dados
em memória inspirado no Redis.

O foco é o aprendizado prático: entender como funciona comunicação via
sockets, protocolos de aplicação, concorrência e sincronização em Go.

A longo prazo, o projeto pode evoluir para uma ferramenta real e útil, mas o
estágio atual é exploratório e educacional.

## Temas a estudar e aplicar

- Sockets TCP (pacote `net`)
- Modelo cliente-servidor
- Concorrência em Go (goroutines)
- Channels e sincronização (`sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`)
- Protocolos de aplicação (design de protocolo tipo RESP)
- Parsing de texto/bytes (`bufio`, `io.Reader`/`io.Writer`)
- Estruturas de dados em memória (map, listas, expiração/TTL)
- Context (cancelamento, timeouts)
- Tratamento de erros em Go
- Testes (unitários e de integração) em rede
- Persistência simples (snapshot / append-only file)
- Benchmarking e profiling (`pprof`)
- Graceful shutdown
- Connection pooling / multiplexação
- Serialização de dados
- Logging estruturado
- Design de arquitetura (camadas: transporte, protocolo, storage)
- Segurança básica (autenticação, timeouts, limites de conexão)

## Regras de negócio

- Os dados armazenados serão no formato chave:valor.
- O tamanho máximo de cada objeto será de 1Kb.
- O tamanho máximo de cada chave será de 256 caracteres.
- O dado deverá obrigatoriamente ser tipado.
- Os tipos permitidos serão `char`, `string`, `int`, `float` e `boolean`.
- Todo dado trafega e é armazenado internamente como `string`. O `{tipo}`
  informado é apenas um metadado que descreve como o dado deve ser
  interpretado, mas não altera a representação interna/no protocolo.
- O tamanho máximo do dado recebido será de 1024 caracteres,
  independentemente do tipo.
- O tipo `char` permitirá apenas um caractere de 1 byte.
- O tipo `string` permitirá strings de até 1024 caracteres.
- O tipo `float` deverá representar um valor compatível com `float64` nativo
  do Go, com as seguintes restrições de formato:
  - O dado terá no máximo 15 dígitos no total, sem contar o ponto `.` e um
    eventual sinal de negativo (`-`).
  - A parte decimal (à direita do ponto `.`) deverá ter até 6 dígitos.
  - O separador decimal deverá obrigatoriamente ser um ponto (`.`).
- O tipo `int` deverá representar um valor compatível com `int64`, com as
  seguintes restrições de formato:
  - O dado terá no máximo 19 dígitos, sem contar um eventual sinal de
    negativo (`-`), correspondendo ao limite de um `int64`
    (-9223372036854775808 a 9223372036854775807).
- O tipo `boolean` aceitará apenas `true` ou `false`, em minúsculo,
  utilizando o tipo `bool` nativo do Go.
- O servidor escutará na porta `8379`, que não é registrada pela IANA para
  nenhuma ferramenta/serviço conhecido e não exige privilégios especiais em
  nenhum sistema operacional (fora da faixa de portas privilegiadas 0-1023).
- Os comandos permitidos serão `get`, `put` e `clean`.
- O banco não terá persistência de dados. Se o servidor cair, perde-se tudo.

### Formato dos comandos

**Comando `put`**
- Formato: `put {chave} {tipo} {dado}`
- Resposta: `status {success ou error}`
- Erro: `{mensagem de erro ou null}`

**Comando `get`**
- Formato: `get {chave}`
- Resposta: `data: {tipo} {dado}` ou `null`
- Erro: `{mensagem de erro ou null}`

**Comando `clean`**
- Formato: `clean {chave}`
- Resposta: `status {success ou error}`
- Erro: `{mensagem de erro ou null}`

## Arquitetura do projeto

O projeto será organizado em 3 camadas: **server**, **controller** e
**aplicação**.

- **Server**: escuta a porta `8379` e aceita conexões dos clientes. Ao
  receber uma requisição, repassa a mensagem bruta ao controller e devolve ao
  cliente a resposta (sucesso ou erro) recebida de volta.
- **Controller**: recebe a mensagem do server e a submete a um **parser**
  antes de chamar a aplicação.
  - O **parser** é responsável por interpretar a mensagem, validando se ela
    segue o protocolo (próprio, inspirado no RESP do Redis). Se a mensagem
    não for válida, o parser rejeita; caso seja válida, retorna somente o
    comando já limpo/estruturado.
  - Se o parser rejeitar a mensagem, o controller monta uma mensagem de erro
    personalizada e a retorna diretamente ao server (sem chamar a
    aplicação).
  - Se o parser validar a mensagem, o controller chama a aplicação passando
    o comando retornado pelo parser.
- **Aplicação**: recebe o comando já validado/parseado do controller e
  processa a lógica de negócio (`put`, `get`, `clean`), retornando o
  resultado ao controller.

