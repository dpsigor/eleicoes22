# eleicoes22

Utilidades para obter e extrair dados de arquivos das eleições 2022

- Processa arquivos .bu (BOLETIM DE URNA) e .logjez (LOG DE URNA) do TSE para extrair alguns dados:
    ```go
    type urna struct {
      buname            string
      logname           string
      municipio         string
      zona              string
      secao             string
      bolso             int64
      lula              int64
      brancos           int64
      nulos             int64
      qtdComparecimento int64
      qtdVotosPR        int64
      qtdTeclaIndevida  int64
      qtdAlertas        int64
      versao            string
      modelo            string
    }

    ```

## TODO

Modo de download dos arquivos do TSE. Note que alguns arquivos (na ordem de 10, dentre os 940mil) vieram corrompidos da primeira vez. Verificar.

## Como utilizar

Primeiro deve-se realizar o download dos arquivos .bu e .logjez no formato: `[uf]_o00407-MMMMMZZZZSSSS.(bu|logjez)`, sendo `uf` ac, al, am, etc. Para processar, o caminho do diretório com esses arquivos deve ser passado para o flag `-dir`

Requerimentos:

- go versão 1.18+

Para obter informações: `go run ./src -h`

Exemplo no script `run.bash`

## Testes

Para garantir que utiliza-se os arquivos da pasta "test", utilizar: `go test -v ./src/`

## ASN.1

A função `readBU` não é mais utilizada, pois implementamos leitura de ans1 por conta própria. Para testar:

- python3 no PATH, com esse nome
- [Documentação técnica do software da urna eletrônica](https://www.tre-mt.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica): realizar download de "Formato dos arquivos de BU, RDV e assinatura digital (formato ZIP)" e extrair. Para processar, o caminho do diretório resultante deve ser passado para o flag `-fa`

## Arquivos com problemas

- Jundiá, AL, Zona 0014 Seção 0036 -> Log da Urna devolve um log identificação incorreta: município 27855 ao invés de 27774, zona 0001 ao invés de 0014 e seção 0000 ao invés de 0036
- Manaus, AM, Zona 68 Seção 0512 -> Log da Urna tem quase nada no segundo turno.
