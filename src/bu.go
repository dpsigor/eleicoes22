package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

type bujson struct {
	EntidadeBoletimUrna struct {
		Cabecalho struct {
			DataGeracao string        `json:"dataGeracao"`
			IDEleitoral []interface{} `json:"idEleitoral"`
		} `json:"cabecalho"`
		ChaveAssinaturaVotosVotavel string        `json:"chaveAssinaturaVotosVotavel"`
		DadosSecaoSA                []interface{} `json:"dadosSecaoSA"`
		DataHoraEmissao             string        `json:"dataHoraEmissao"`
		Fase                        string        `json:"fase"`
		IdentificacaoSecao          struct {
			Local         int64 `json:"local"`
			MunicipioZona struct {
				Municipio int64 `json:"municipio"`
				Zona      int64 `json:"zona"`
			} `json:"municipioZona"`
			Secao int64 `json:"secao"`
		} `json:"identificacaoSecao"`
		QtdEleitoresCompBiometrico  int64 `json:"qtdEleitoresCompBiometrico"`
		QtdEleitoresLibCodigo       int64 `json:"qtdEleitoresLibCodigo"`
		ResultadosVotacaoPorEleicao []struct {
			IDEleicao         int64 `json:"idEleicao"`
			QtdEleitoresAptos int64 `json:"qtdEleitoresAptos"`
			ResultadosVotacao []struct {
				QtdComparecimento int64  `json:"qtdComparecimento"`
				TipoCargo         string `json:"tipoCargo"`
				TotaisVotosCargo  []struct {
					CodigoCargo    []string `json:"codigoCargo"`
					OrdemImpressao int64    `json:"ordemImpressao"`
					VotosVotaveis  []struct {
						Assinatura           string `json:"assinatura"`
						IdentificacaoVotavel struct {
							Codigo  int64 `json:"codigo"`
							Partido int64 `json:"partido"`
						} `json:"identificacaoVotavel"`
						QuantidadeVotos int64  `json:"quantidadeVotos"`
						TipoVoto        string `json:"tipoVoto"`
					} `json:"votosVotaveis"`
				} `json:"totaisVotosCargo"`
			} `json:"resultadosVotacao"`
		} `json:"resultadosVotacaoPorEleicao"`
		// Urna struct {
		// 	CorrespondenciaResultado struct {
		// 		Carga struct {
		// 			CodigoCarga       string `json:"codigoCarga"`
		// 			DataHoraCarga     string `json:"dataHoraCarga"`
		// 			NumeroInternoUrna int64  `json:"numeroInternoUrna"`
		// 			NumeroSerieFC     string `json:"numeroSerieFC"`
		// 		} `json:"carga"`
		// 		Identificacao []interface{} `json:"identificacao"`
		// 	} `json:"correspondenciaResultado"`
		// 	NumeroSerieFV string `json:"numeroSerieFV"`
		// 	TipoArquivo   string `json:"tipoArquivo"`
		// 	TipoUrna      string `json:"tipoUrna"`
		// 	VersaoVotacao string `json:"versaoVotacao"`
		// } `json:"urna"`
	} `json:"EntidadeBoletimUrna"`
	// EntidadeEnvelopeGenerico struct {
	// 	Cabecalho struct {
	// 		DataGeracao string        `json:"dataGeracao"`
	// 		IDEleitoral []interface{} `json:"idEleitoral"`
	// 	} `json:"cabecalho"`
	// 	Fase          string        `json:"fase"`
	// 	Identificacao []interface{} `json:"identificacao"`
	// 	TipoEnvelope  string        `json:"tipoEnvelope"`
	// } `json:"EntidadeEnvelopeGenerico"`
}

// readBU extrai dados do arquivo .bu
// Há dois passos. Segundo especificação do TSE,
// (https://www.tre-mt.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica)
// precisamos utilizar um script dado pelo TSE ("bu_dump.py") para obter o
// os dados em plain text. Entretanto, esse plain text tem um formato
// bastante particular. Por conveniência, damos parse nesse plain text
// utilizando um script próprio ("bu2json.py")
func readBU(fpath string, buDump, buSpec string) bujson {
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		log.Fatalf("err creating pipe: %s\n", err)
	}
	b := bytes.NewBuffer(nil)
	cmd := exec.Command("python3", buDump, "-a", buSpec, "-b", fpath)
	cmd2 := exec.Command("python3", bu2json)
	cmd.Stdout = pipeW
	cmd2.Stdin = pipeR
	cmd2.Stdout = b
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd start %s: %s\n", fpath, err)
	}
	if err := cmd2.Start(); err != nil {
		log.Fatalf("cmd2 start %s: %s\n", fpath, err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("cmd wait %s: %s\n", fpath, err)
	}
	pipeW.Close()
	if err := cmd2.Wait(); err != nil {
		log.Fatalf("cmd2 wait %s: %s\n", fpath, err)
	}
	buBytes, err := io.ReadAll(b)
	if err != nil {
		log.Fatalf("io read all %s: %s\n", fpath, err)
	}
	buBytes = bytes.ReplaceAll(buBytes, []byte("Infinity"), []byte("0"))
	bu := bujson{}
	if err := json.Unmarshal(buBytes, &bu); err != nil {
		log.Fatalf("decode json %s: %s\n", fpath, err)
	}
	return bu
}
