package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type irisInfo struct {
	sepalLength float64
	sepalWidth  float64
	petalLength float64
	petalWidth  float64
	species     string
}

func main() {
	irisDF, err := os.Open("iris9.data")
	if err != nil {
		log.Fatal(err)
	}
	defer irisDF.Close()

	reader := csv.NewReader(irisDF)
	reader.Comma = ','

	var recordSet []irisInfo

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recordSet = append(recordSet, parseIrisRecord(record))
	}

	//K Nearest Neighbors

	//fmt.Println(recordSet[0].sepalLength + 0.4)

	// tipo de flor desconocida
	var irisUnknown irisInfo
	irisUnknown.sepalLength = 5.0
	irisUnknown.sepalWidth = 3.3
	irisUnknown.petalLength = 1.5
	irisUnknown.petalWidth = 1.1
	irisUnknown.species = "?"

	recordSet = append(recordSet, irisUnknown)

	var euclideanList []float64

	//se deshace el irisInfo añadido por el usuario
	lengthR := len(recordSet) - 1
	//Calcula la distancia euclidiana para todos los objetos del recordSet
	for i := 0; i < lengthR; i++ {
		euclideanList = append(euclideanList, euclideanDist(recordSet[i], recordSet[len(recordSet)-1]))
	}

	fmt.Println(recordSet[5])
	fmt.Println(recordSet[len(recordSet)-1])
	//fmt.Println(euclideanDist(recordSet[5], recordSet[len(recordSet)-1]))
	fmt.Println(euclideanList)

	//Encontrar las K distancias menores
	//DEFINICION DE K
	k := 5

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

	}
	//fmt.Println(euclideanList)
	fmt.Println(menores)
	//fmt.Println(copiaEuLista)
	fmt.Println(recordSet)

	//encontrar los menores recorriendo todo el set de nuevo y
	//analizar si es versicolor, setosa o virginica
	var irisSpecies []string
	for j := 0; j < k; j++ {
		for i := 0; i < lengthR; i++ {
			if menores[j] == copiaEuLista[i] {
				irisSpecies = append(irisSpecies, recordSet[i].species)
			}
		}
	}
	//K nearest neighbors
	fmt.Println(irisSpecies)

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

func parseIrisRecord(record []string) irisInfo {
	var iris irisInfo

	iris.sepalLength, _ = strconv.ParseFloat(record[0], 64)
	iris.sepalWidth, _ = strconv.ParseFloat(record[1], 64)
	iris.petalLength, _ = strconv.ParseFloat(record[2], 64)
	iris.petalWidth, _ = strconv.ParseFloat(record[3], 64)
	iris.species = record[4]

	return iris
}

func euclideanDist(iris irisInfo, unknownIris irisInfo) float64 {
	var distance float64
	//distance := 0.0

	distance += math.Pow(iris.sepalLength-unknownIris.sepalLength, 2)
	distance += math.Pow(iris.sepalWidth-unknownIris.sepalWidth, 2)
	distance += math.Pow(iris.petalLength-unknownIris.petalLength, 2)
	distance += math.Pow(iris.petalWidth-unknownIris.petalWidth, 2)

	distance = math.Sqrt(distance)
	return distance
}
