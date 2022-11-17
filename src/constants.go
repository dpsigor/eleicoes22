package main

var LF = []byte("\n")
var PR = []byte("Voto confirmado para [Presidente]")
var TS = []byte("30/10/2022")
var TSMinus1 = []byte("29/10/2022") // Há urnas que foram ligadas no dia anterior
var TSMinus2 = []byte("28/10/2022") // Há urnas que foram ligadas dois dias antes: ba_o00407-3413401710232.logjez
var MU = []byte{'M', 'u', 'n', 'i', 'c', 'í', 'p', 'i', 'o', ':', ' '}
var ZN = []byte{'Z', 'o', 'n', 'a', ' ', 'E', 'l', 'e', 'i', 't', 'o', 'r', 'a', 'l', ':', ' '}
var SE = []byte{'S', 'e', 'ç', 'ã', 'o', ' ', 'E', 'l', 'e', 'i', 't', 'o', 'r', 'a', 'l'}
var VRS = []byte{'V', 'e', 'r', 's', 'ã', 'o'}
var ModUrna = []byte("Modelo de Urna")
var TCLIND = []byte("Tecla indevida")
var AL = []byte("ALERTA")
