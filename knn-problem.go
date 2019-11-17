package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var species []string

var remotehost string

var irisUnknown IrisInfo

type IrisInfo struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
}

type NodeInfo struct {
	Num           int
	NumNodes      int
	IrisU         IrisInfo
	RecordIrisSet []IrisInfo
	K             int
}

func main() {

	irisDF, err := os.Open("iris9.data")
	if err != nil {
		log.Fatal(err)
	}
	defer irisDF.Close()

	reader := csv.NewReader(irisDF)
	reader.Comma = ','

	var recordSet []IrisInfo

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//todas las filas en una lista
		recordSet = append(recordSet, parseIrisRecord(record))
	}

	// tipo de flor desconocida

	irisUnknown.SepalLength = 5.0
	irisUnknown.SepalWidth = 3.3
	irisUnknown.PetalLength = 1.5
	irisUnknown.PetalWidth = 1.1
	irisUnknown.Species = "?"

	//Comienza
	var n int
	var K int

	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese su puerto: ")
	port, _ := gin.ReadString('\n')
	port = strings.TrimSpace(port)
	//hostname := fmt.Sprintf("localhost:%s", port)

	fmt.Print("Ingrese cantidad de nodos (puertos): ")
	fmt.Scanf("%d\n", &n)
	var remotehosts = make([]string, n)

	for i := range remotehosts {
		fmt.Printf("Puerto %d: ", i+1)
		fmt.Scanf("%s\n", &(remotehosts[i]))
	}
	//fmt.Println(remotehosts)

	fmt.Print("Ingrese cantidad K : ")
	fmt.Scanf("%d\n", &K)

	go func() {
		fmt.Println("Presione enter para comenzar...")
		r := bufio.NewReader(os.Stdin)
		r.ReadString('\n')

		for i := range remotehosts {
			var sendInfo NodeInfo
			sendInfo.Num = i
			sendInfo.NumNodes = n
			sendInfo.IrisU = irisUnknown
			sendInfo.RecordIrisSet = recordSet
			sendInfo.K = K
			remotehost := fmt.Sprintf("localhost:%s", remotehosts[i])
			//fmt.Println(remotehost)
			go send(sendInfo, remotehost)
		}
	}()

	server()

}

func server() {
	host := fmt.Sprintf("localhost:8000")
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}

func send(sendInfo NodeInfo, rh string) {
	conn, _ := net.Dial("tcp", rh)
	defer conn.Close()
	//fmt.Println(sendInfo)
	bMsg, _ := json.Marshal(sendInfo)

	fmt.Fprintln(conn, string(bMsg))
}

func handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	var species []string
	json.Unmarshal([]byte(msg), &species)

	fmt.Println(species)

	response(species)
}
func parseIrisRecord(record []string) IrisInfo {
	var iris IrisInfo

	iris.SepalLength, _ = strconv.ParseFloat(record[0], 64)
	iris.SepalWidth, _ = strconv.ParseFloat(record[1], 64)
	iris.PetalLength, _ = strconv.ParseFloat(record[2], 64)
	iris.PetalWidth, _ = strconv.ParseFloat(record[3], 64)
	iris.Species = record[4]

	return iris
}

func response(irisSpecies []string) {
	typeSetosa := 0
	typeVersicolor := 0
	typeVirginica := 0

	for i := range irisSpecies {
		switch irisSpecies[i] {
		case "Iris-setosa":
			typeSetosa++
		case "Iris-versicolor":
			typeVersicolor++
		case "Iris-virginica":
			typeVirginica++

		}
	}

	fmt.Println(typeSetosa, typeVirginica, typeVersicolor)
	if typeVirginica > typeSetosa && typeVirginica > typeVersicolor {
		fmt.Println("El tipo de flor para el set de datos ", irisUnknown, " es Virginica")
	}
	if typeVersicolor > typeSetosa && typeVersicolor > typeVirginica {
		fmt.Println("El tipo de flor para el set de datos ", irisUnknown, " es Versicolor")
	}
	if typeSetosa > typeVirginica && typeSetosa > typeVersicolor {
		fmt.Println("El tipo de flor para el set de datos ", irisUnknown, " es Setosa")
	}
}
