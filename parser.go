package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type lataItems struct {
	XMLName xml.Name   `xml:"root"`
	Items   []lataItem `xml:"prefixdata"`
}

type lataItem struct {
	XMLName  xml.Name `xml:"prefixdata"`
	CityName string   `xml:"rc"`
	Npa      string   `xml:"npa"`
	Nxx      string   `xml:"nxx"`
	Lata     string   `xml:"lata"`
	Region   string   `xml:"region"`
}

// LastCityID is used for assign ID to city while parsing
var LastCityID int

func main() {
	// get list of files to process
	dataFiles, err := ioutil.ReadDir("./data")
	checkError("Cannot prepare file list", err)

	var lataData lataItems

	outDataCities := make(map[string]map[int]string)
	outDataNpaNxx := make(map[int]map[int][]int)
	outDataLata := make(map[int][]int)

	for _, file := range dataFiles {

		filename := "data/" + file.Name()

		// read file
		xmlFile, err := os.Open(filename)
		checkError("Cannot open file", err)

		fmt.Println("Successfully Opened: ", filename)

		defer xmlFile.Close()

		// parse data file
		byteValue, _ := ioutil.ReadAll(xmlFile)
		xml.Unmarshal(byteValue, &lataData)
		byteValue = nil

		for _, dataLataItem := range lataData.Items {

			var localCityID int

			// collect cities data
			localCityID = pushToCitiesData(dataLataItem, &outDataCities)

			// collect lata data
			pushToLataData(dataLataItem, localCityID, &outDataLata)

			// collect npa, nxx
			pushToNpaNxxData(dataLataItem, localCityID, &outDataNpaNxx)

		}
	}

	fmt.Println("Insert")

	preparedCitiesData, err := yaml.Marshal(&outDataCities)
	checkError("Cannot prepare yml for cities", err)

	writeCities := ioutil.WriteFile("output/cities.yml", preparedCitiesData, 0644)
	checkError("Failed to write to file cities data", writeCities)

	preparedLataData, err := yaml.Marshal(&outDataLata)
	checkError("Cannot prepare yml for cities", err)

	writeLata := ioutil.WriteFile("output/lata.yml", preparedLataData, 0644)
	checkError("Failed to write to file LATA data", writeLata)

	preparedNpaNxxData, err := yaml.Marshal(&outDataNpaNxx)
	checkError("Cannot prepare yml for cities", err)

	writeNpaNxx := ioutil.WriteFile("output/npa_nxx.yml", preparedNpaNxxData, 0644)
	checkError("Failed to write to file NPA, NXX data", writeNpaNxx)
}

func pushToLataData(item lataItem, cityID int, outData *map[int][]int) {

	// convert lata from string to int
	outLataInt, err := strconv.Atoi(item.Lata)
	checkError("Failed to parse Lata", err)

	// check if lata key already present in outData
	lataKey, lataKeyPresent := (*outData)[outLataInt]

	if !lataKeyPresent {
		(*outData)[outLataInt] = []int{cityID}
		return
	}

	if inSlice(lataKey, cityID) {
		return
	}

	(*outData)[outLataInt] = append((*outData)[outLataInt], cityID)
}

func pushToNpaNxxData(item lataItem, cityID int, outData *map[int]map[int][]int) {

	// convert NPA from string to int
	outNpaInt, err := strconv.Atoi(item.Npa)
	checkError("Failed to parse NPA", err)

	// convert NXX from string to int
	outNxxInt, err := strconv.Atoi(item.Nxx)
	checkError("Failed to parse NXX", err)

	// check if NPA key already present in outData
	_, npaKeyPresent := (*outData)[outNpaInt]

	if !npaKeyPresent {
		(*outData)[outNpaInt] = make(map[int][]int)
	}

	// check if city ID key already present in NPA data
	_, cityKeyPresent := (*outData)[outNpaInt][cityID]

	if !cityKeyPresent {
		(*outData)[outNpaInt][cityID] = []int{outNxxInt}
	}

	if inSlice((*outData)[outNpaInt][cityID], outNxxInt) {
		return
	}

	(*outData)[outNpaInt][cityID] = append((*outData)[outNpaInt][cityID], outNxxInt)
}

func pushToCitiesData(item lataItem, outData *map[string]map[int]string) int {

	var hasCity bool = false
	var cityID int
	var region string

	region = item.Region
	_, regionKeyPresent := (*outData)[region]

	if !regionKeyPresent {
		(*outData)[region] = make(map[int]string)
	}

	cityName := item.CityName

	for k, v := range (*outData)[region] {
		if v == cityName {
			hasCity = true
			cityID = k
			break
		}
	}

	if hasCity == true {
		return cityID
	}

	LastCityID = LastCityID + 1
	(*outData)[region][LastCityID] = cityName

	return LastCityID
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func inSlice(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
