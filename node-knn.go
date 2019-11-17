package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
)

var remotehostt string

type IrisInfos struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
}

type NodeInfos struct {
	Num           int
	NumNodes      int
	IrisU         IrisInfos
	RecordIrisSet []IrisInfos
	K             int
}

func main() {
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Enter port: ")
	port, _ := gin.ReadString('\n')
	port = strings.TrimSpace(port)
	hostname := fmt.Sprintf("localhost:%s", port)

	// Listener!
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handdle(conn)
	}
}

func handdle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	var info NodeInfos
	json.Unmarshal([]byte(str), &info)

	//fmt.Println(info.RecordIrisSet, info.IrisU)
	lengthR := len(info.RecordIrisSet)
	divisionRecordIris := lengthR / info.NumNodes
	var euclideanList []float64

	finalRecord := divisionRecordIris * (info.Num + 1)
	comienzoRecord := divisionRecordIris * (info.Num)

	for i := comienzoRecord; i < finalRecord; i++ {
		euclideanList = append(euclideanList, euclideanDistance(info.RecordIrisSet[i], info.IrisU))
	}

	k := info.K

	var min float64
	var indiceEncontrado, indiceEliminar = 0, 0
	var menores []float64
	//copia de euclideanList
	copiaEuLista := make([]float64, len(euclideanList))
	copy(copiaEuLista, euclideanList)
	for i := 0; i < k; i++ {
		min = math.MaxInt64
		for j := range euclideanList {
			if min > euclideanList[j] {
				min = euclideanList[j]
				indiceEncontrado = j
			}
		}
		indiceEliminar = indiceEncontrado
		menores = append(menores, euclideanList[indiceEliminar])
		//eliminar valor de euclideanList para volver a recorrer la lista y encontrar al siguiente menor
		euclideanList[indiceEliminar] = euclideanList[len(euclideanList)-1]
		euclideanList[len(euclideanList)-1] = 0.0
		euclideanList = euclideanList[:len(euclideanList)-1]
		//valor eliminado
	}

	//fmt.Println(euclideanList)
	//fmt.Println(menores)
	//fmt.Println(copiaEuLista)

	//encontrar los menores recorriendo todo el set de nuevo y
	//analizar si es versicolor, setosa o virginica
	var irisSpecies []string
	for j := 0; j < k; j++ {
		for i := 0; i < lengthR; i++ {
			if menores[j] == copiaEuLista[i] {
				irisSpecies = append(irisSpecies, info.RecordIrisSet[i].Species)
			}
		}
	}
	//K nearest neighbors
	//fmt.Println(irisSpecies)

	// aqui se envian las especies de iris al nodo principal para que muestre el resultado
	sendd(irisSpecies)

}

func sendd(irisSpecies []string) {
	remotehostt := "localhost:8000"
	conn, _ := net.Dial("tcp", remotehostt)
	defer conn.Close()
	fmt.Println(irisSpecies)
	bMsg, _ := json.Marshal(irisSpecies)

	fmt.Fprintln(conn, string(bMsg))
}

func euclideanDistance(iris IrisInfos, unknownIris IrisInfos) float64 {
	var distance float64
	//distance := 0.0

	distance += math.Pow(iris.SepalLength-unknownIris.SepalLength, 2)
	distance += math.Pow(iris.SepalWidth-unknownIris.SepalWidth, 2)
	distance += math.Pow(iris.PetalLength-unknownIris.PetalLength, 2)
	distance += math.Pow(iris.PetalWidth-unknownIris.PetalWidth, 2)

	distance = math.Sqrt(distance)
	return distance
}
