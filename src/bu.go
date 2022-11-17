package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"os/exec"

	ber "github.com/go-asn1-ber/asn1-ber"
)

type TipoVoto int

const (
	tipoVotoNominal           TipoVoto = 1
	tipoVotoBranco            TipoVoto = 2
	tipoVotoNulo              TipoVoto = 3
	tipoVotoLegenda           TipoVoto = 4
	tipoVotoCargoSemCandidato TipoVoto = 5
)

func (t TipoVoto) String() string {
	if t == tipoVotoNominal {
		return "nominal"
	}
	if t == tipoVotoBranco {
		return "branco"
	}
	if t == tipoVotoNulo {
		return "nulo"
	}
	return "outro"
}

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

type Presidente int

const (
	bolso Presidente = 22
	lula  Presidente = 13
)

type BUVoto struct {
	Qtd      int64
	Voto     Presidente
	TipoVoto TipoVoto
}

type BU struct {
	DataGeracao       time.Time
	DataHoraEmissao   time.Time
	Municipio         int64
	Zona              int64
	Local             int64
	Secao             int64
	IDEleicao         int64
	QtdEleitoresAptos int64
	QtdComparecimento int64
	Votos             []BUVoto
}

func berInt(b []byte) int64 {
	ret, err := ber.ParseInt64(b)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

func fmtTime(s string) time.Time {
	t, err := time.Parse("20060102T150405", s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// processaBU tem implementação própria do parse do arquivo .bu, que tem dados em ASN.1
// É tremendamente mais rápido que utilizando readBU, o qual precisa invocar dois scripts .py
func processaBU(bupath string) BU {
	bu := BU{}
	f, err := os.Open(bupath)
	if err != nil {
		log.Fatal(err)
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	packet := ber.DecodePacket(bs)
	if len(packet.Children) != 5 {
		log.Fatalf("%s: expected packet to have 5 children, got %d\n", bupath, len(packet.Children))
	}
	entidadeBU := ber.DecodePacket(packet.Children[4].ByteValue)
	cabecalho := entidadeBU.Children[0]
	if len(cabecalho.Children) != 2 {
		log.Fatalf("%s: expected cabecalho to have 2 children, got %d\n", bupath, len(packet.Children))
	}
	bu.DataGeracao = fmtTime(string(cabecalho.Children[0].ByteValue))
	bu.DataHoraEmissao = fmtTime(string(entidadeBU.Children[4].ByteValue))
	// x := entidadeBU.Children[3] // identificacaoSecao
	// x := entidadeBU.Children[5] // dois timestamps // dadosSecaoSA
	// x := entidadeBU.Children[6] // é um boolean // qtdEleitoresLibCodigo. Ignorar primeiros dois bytes
	// x := entidadeBU.Children[8] // resultadosVotacaoPorEleicao...nem sempre há índice 8
	// x := entidadeBU.Children[7] // é um inteiro // qtdEleitoresCompBiometrico. Ignorar primeiros dois bytes
	// x := entidadeBU.Children[9] // é um octet string

	var resultadosVotacaoPorEleicao *ber.Packet
	if len(entidadeBU.Children) > 8 {
		resultadosVotacaoPorEleicao = entidadeBU.Children[8] // contém um child, e este contem 3 children
	} else {
		for _, c := range entidadeBU.Children {
			if c.Tag == 3 && c.TagType == 32 && c.ClassType == 128 {
				resultadosVotacaoPorEleicao = c
				break
			}
		}
	}
	if resultadosVotacaoPorEleicao == nil {
		log.Fatalf("%s: didn't find resultadosVotacaoPorEleicao\n", bupath)
	}

	identificacaoSecao := entidadeBU.Children[3]
	if len(identificacaoSecao.Children) != 3 {
		log.Fatalf("%s: expected identificacaoSecao to have 3 children, got %d\n", bupath, len(identificacaoSecao.Children))
	}
	municipioZona := identificacaoSecao.Children[0]
	if len(municipioZona.Children) != 2 {
		log.Fatalf("%s: expected municipioZona to have 2 children, got %d\n", bupath, len(municipioZona.Children))
	}
	bu.Municipio = berInt(municipioZona.Children[0].ByteValue)
	bu.Zona = berInt(municipioZona.Children[1].ByteValue)
	bu.Local = berInt(identificacaoSecao.Children[1].ByteValue)
	bu.Secao = berInt(identificacaoSecao.Children[2].ByteValue)
	if len(resultadosVotacaoPorEleicao.Children) != 1 && len(resultadosVotacaoPorEleicao.Children) != 2 {
		log.Fatalf("%s: expected resultadosVotacaoPorEleicao.Children to have 1 ou 2 child, got %d\n", bupath, len(resultadosVotacaoPorEleicao.Children))
	}
	idx := 0
	bu.IDEleicao = berInt(resultadosVotacaoPorEleicao.Children[idx].Children[0].ByteValue)
	if bu.IDEleicao != 545 {
		idx = 1
		bu.IDEleicao = berInt(resultadosVotacaoPorEleicao.Children[idx].Children[0].ByteValue)
	}
	bu.QtdEleitoresAptos = berInt(resultadosVotacaoPorEleicao.Children[idx].Children[1].ByteValue)
	resultadoPkt := resultadosVotacaoPorEleicao.Children[idx].Children[2]
	if resultadoPkt.Tag != ber.TagSequence {
		log.Fatalf("%s: expected to have resultados da eleicao as tag sequence, got %d\n", bupath, resultadosVotacaoPorEleicao.Children[idx].Children[2].Tag)
	}
	if len(resultadoPkt.Children) != 1 {
		log.Fatalf("%s: expected resultadoPkt to have one child, got %d\n", bupath, len(resultadoPkt.Children))
	}
	resultado := resultadoPkt.Children[0]
	if len(resultado.Children) != 3 {
		log.Fatalf("%s: expected resultado to have 3 children, got %d\n", bupath, len(resultado.Children))
	}
	bu.QtdComparecimento = berInt(resultado.Children[1].ByteValue)
	totaisVotosCargoSeq := resultado.Children[2]
	if len(totaisVotosCargoSeq.Children) != 1 {
		log.Fatalf("%s: expected totaisVotosCargoSeq to have 1 child, got %d\n", bupath, len(totaisVotosCargoSeq.Children))
	}
	totaisVotosCargo := totaisVotosCargoSeq.Children[0]
	if len(totaisVotosCargo.Children) != 3 {
		log.Fatalf("%s: expected totaisVotosCargo to have 3 children, got %d\n", bupath, len(totaisVotosCargo.Children))
	}
	votosVotaveis := totaisVotosCargo.Children[2]
	if len(votosVotaveis.Children) < 1 {
		log.Fatalf("%s: expected votosVotaveis to have at least 1 child, got %d\n", bupath, len(votosVotaveis.Children))
	}
	for _, votoVotaveis := range votosVotaveis.Children {
		if votoVotaveis.Tag != ber.TagSequence {
			log.Fatalf("%s: expected votoVotaveis to be tag sequence, got %d\n", bupath, votoVotaveis.Tag)
		}
		if len(votoVotaveis.Children) < 2 {
			log.Fatalf("%s: expected votoVotaveis to have at least 2 children, got %d\n", bupath, len(votoVotaveis.Children))
		}
		tipoVoto := TipoVoto(berInt(votoVotaveis.Children[0].Bytes()[2:]))
		quantidadeVotos := berInt(votoVotaveis.Children[1].Bytes()[2:])
		if tipoVoto == tipoVotoNominal {
			identificacaoVotavel := votoVotaveis.Children[2]
			if len(identificacaoVotavel.Children) != 2 {
				log.Fatalf("%s: expected identificacaoVotavel to have 2 children, got %d\n", bupath, len(identificacaoVotavel.Children))
			}
			voto := berInt(identificacaoVotavel.Children[0].ByteValue)
			buVoto := BUVoto{
				Qtd:      quantidadeVotos,
				Voto:     Presidente(voto),
				TipoVoto: tipoVoto,
			}
			bu.Votos = append(bu.Votos, buVoto)
		} else {
			buVoto := BUVoto{
				Qtd:      quantidadeVotos,
				TipoVoto: tipoVoto,
			}
			bu.Votos = append(bu.Votos, buVoto)
		}
	}
	return bu
}

// readBU extrai dados do arquivo .bu
// NÃO É MAIS UTILIZADO
// Há dois passos. Segundo especificação do TSE,
// (https://www.tre-mt.jus.br/eleicoes/eleicoes-2022/documentacao-tecnica-do-software-da-urna-eletronica)
// precisamos utilizar um script dado pelo TSE ("bu_dump.py") para obter o
// os dados em plain text. Entretanto, esse plain text tem um formato
// bastante particular. Por conveniência, damos parse nesse plain text
// utilizando um script próprio ("bu2json.py")
func readBU(fpath string, buDump, buSpec string) BU {
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
	if len(bu.EntidadeBoletimUrna.ResultadosVotacaoPorEleicao) < 1 {
		log.Fatalf("bu sem resultado por eleição: %s\n", fpath)
	}
	var idEleicao int64 = 0
	var qtdEleitoresAptos int64 = 0
	var qtdComparecimento int64 = 0
	votos := make([]BUVoto, 0)
	for _, v := range bu.EntidadeBoletimUrna.ResultadosVotacaoPorEleicao {
		if v.IDEleicao != 545 {
			continue
		}
		if len(v.ResultadosVotacao) < 1 {
			log.Fatalf("bu sem resultadosVotacao len: %s\n", fpath)
		}
		idEleicao = v.IDEleicao
		qtdEleitoresAptos = v.QtdEleitoresAptos
		qtdComparecimento = v.ResultadosVotacao[0].QtdComparecimento
		if len(v.ResultadosVotacao[0].TotaisVotosCargo) != 1 {
			log.Fatalf("bujson deveria ter apenas um TotaisVotosCargo: %s\n", fpath)
		}
		totais := v.ResultadosVotacao[0].TotaisVotosCargo[0].VotosVotaveis
		for _, total := range totais {
			var tipoVoto TipoVoto
			if total.TipoVoto == "branco" {
				tipoVoto = tipoVotoBranco
			} else if total.TipoVoto == "nulo" {
				tipoVoto = tipoVotoNulo
			} else {
				tipoVoto = tipoVotoNominal
			}
			voto := BUVoto{
				Qtd:      total.QuantidadeVotos,
				TipoVoto: tipoVoto,
			}
			if tipoVoto == tipoVotoNominal {
				voto.Qtd = total.QuantidadeVotos
			}
			votos = append(votos, voto)
		}
		break
	}
	return BU{
		DataGeracao:       fmtTime(bu.EntidadeBoletimUrna.Cabecalho.DataGeracao),
		DataHoraEmissao:   fmtTime(bu.EntidadeBoletimUrna.DataHoraEmissao),
		Municipio:         bu.EntidadeBoletimUrna.IdentificacaoSecao.MunicipioZona.Municipio,
		Zona:              bu.EntidadeBoletimUrna.IdentificacaoSecao.MunicipioZona.Zona,
		Local:             bu.EntidadeBoletimUrna.IdentificacaoSecao.Local,
		Secao:             bu.EntidadeBoletimUrna.IdentificacaoSecao.Secao,
		IDEleicao:         idEleicao,
		QtdEleitoresAptos: qtdEleitoresAptos,
		QtdComparecimento: qtdComparecimento,
		Votos:             votos,
	}
}
