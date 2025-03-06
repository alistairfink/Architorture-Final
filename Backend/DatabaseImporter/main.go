package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Card struct {
	Id              string
	CardTypeId      string
	Name            string
	Description     string
	PlayImmediately string
	Quantity        string
	ExpansionId     string
	Archivable      string
}

type CardType struct {
	Id   string
	Name string
}

type Expansion struct {
	Id   string
	Name string
}

func main() {
	args := os.Args
	if len(args) != 3 {
		println("ERROR: Input CSV and Table required.")
		println("0: Card")
		println("1: CardType")
		println("2: Expansion")
		return
	}

	inputCSV := args[1]
	procedure := args[2]

	csvFile, err := os.Open(inputCSV)
	if err != nil {
		println(err.Error())
		return
	}

	reader := csv.NewReader(csvFile)
	if procedure == "0" {
		handleCardImport(reader)
	} else if procedure == "1" {
		handleCardTypeImport(reader)
	} else if procedure == "2" {
		handleExpansionImport(reader)
	}
}

func handleCardImport(reader *csv.Reader) {
	records := []Card{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			println(err.Error())
			return
		}

		records = append(records, Card{
			Id:              record[0],
			CardTypeId:      record[1],
			Name:            strings.ReplaceAll(record[2], "'", "''"),
			Description:     strings.ReplaceAll(record[3], "'", "''"),
			PlayImmediately: record[4],
			Quantity:        record[6],
			ExpansionId:     record[7],
			Archivable:      record[5],
		})
	}

	sqlBuilder := "INSERT INTO \"card\" (\"id\", \"card_type_id\", \"name\", \"description\", \"play_immediately\", \"quantity\", \"expansion_id\", \"archivable\") VALUES "
	for i, r := range records {
		sqlBuilder += "(" + r.Id + ", " + r.CardTypeId + ", '" + r.Name + "', '" + r.Description + "', " + r.PlayImmediately + ", " + r.Quantity + ", " + r.ExpansionId + ", " + r.Archivable + ")"
		if i == len(records)-1 {
			sqlBuilder += ";"
		} else {
			sqlBuilder += ", "
		}
	}

	println(sqlBuilder)
}

func handleCardTypeImport(reader *csv.Reader) {
	records := []CardType{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			println(err.Error())
			return
		}

		records = append(records, CardType{
			Id:   record[0],
			Name: strings.ReplaceAll(record[1], "'", "''"),
		})
	}

	sqlBuilder := "INSERT INTO \"card_type\" (\"id\", \"name\") VALUES "
	for i, r := range records {
		sqlBuilder += "(" + r.Id + ", '" + r.Name + "')"
		if i == len(records)-1 {
			sqlBuilder += ";"
		} else {
			sqlBuilder += ", "
		}
	}

	println(sqlBuilder)
}

func handleExpansionImport(reader *csv.Reader) {
	records := []Expansion{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			println(err.Error())
			return
		}

		records = append(records, Expansion{
			Id:   record[0],
			Name: strings.ReplaceAll(record[1], "'", "''"),
		})
	}

	sqlBuilder := "INSERT INTO \"expansion\" (\"id\", \"name\") VALUES "
	for i, r := range records {
		sqlBuilder += "(" + r.Id + ", '" + r.Name + "')"
		if i == len(records)-1 {
			sqlBuilder += ";"
		} else {
			sqlBuilder += ", "
		}
	}

	println(sqlBuilder)
}
