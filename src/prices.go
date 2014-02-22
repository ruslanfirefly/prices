package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"os"
	"bufio"
	"strings"
	"fmt"
	"io"
)

type Menu struct {
	Menu []Punct `json:"menu"`
}

type Punct struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Parentid int     `json:"parentid"`
}

type Product struct {
	Articul     string
	Name        string
	Price       string
	Manufactura string
	Menu        string
}
type Tovar struct {
	Site        string
	Link        string
	NameProduct string
	Action_acia string
	Action_new  string
	Price       string
	Art         string
	Sex         string
	Age         string
	Descrip     string
	Pic         []string
}
type TovarBitrix struct {
	Art         string
	NameProduct string
	Action_acia string
	Action_new  string
	Price       string
	Sex         string
	Age         string
	Descrip     string
	Pic         []string
	Provider    string
	Menu        string
}

var (
	menu bool
	menuFile string
	mongoURL string
	upPrice bool
	newPrice string
)

func init() {
	flag.BoolVar(&menu, "menu", false, "Update menu")
	flag.BoolVar(&upPrice, "upprice", false, "Update Price")
	flag.StringVar(&newPrice, "filePrice", "", "New Price")
	flag.StringVar(&menuFile, "menuFile", "", "JSONFile menu")
	flag.StringVar(&mongoURL, "mdb", "127.0.0.1:27017", "Connect MongoDB")
}

func saveMenu() {
	var menuFromFile Menu
	dropOldMenu()
	text, err := ioutil.ReadFile(menuFile)
	err = json.Unmarshal(text, &menuFromFile)
	error_log(err)
	for _, v := range menuFromFile.Menu {
		addMenu(v)
	}
}

func savePrice() {
	file, err := os.Open(newPrice)
	error_log(err)
	reader := bufio.NewReader(file)
	for {
		var prod Product
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			error_log(err)
		}
		str := strings.Split(string(line), ";")
		prod.Articul = prepareStrings(str[0])
		prod.Name = prepareStrings(str[1])
		prod.Price = prepareStrings(str[2])
		prod.Manufactura = prepareStrings(str[3])
		prod.Menu = prepareStrings(str[4])
		saveOrUpdateProduct(prod)
	}
}
func main() {
	flag.Parse()
	if (menu) {
		saveMenu()
	}
	if (upPrice) {
		savePrice()
	}
	updatePrice()
	fmt.Println("the end")
}
