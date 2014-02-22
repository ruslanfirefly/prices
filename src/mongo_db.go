package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"fmt"

)

const (
	DBMNG string     = "parseprice"
	DBMNG_PARSE string = "parsetoys"
	DBMNG_PRICE_BITRIX string = "price_bitrix"
	COLPRICES_BITRIX string = "price"
	COLNOTNEED string = "not_need"
	COLTOYS string = "toysinfo"
	COLMENU string   = "menu"
	COLPRICES string = "prices"
)

func createMongoSession() *mgo.Session {
	session, err := mgo.Dial(mongoURL)
	error_log(err)
	return session
}

func dropOldMenu() {
	s := createMongoSession()
	defer s.Close()
	listCol, err := s.DB(DBMNG).C(COLMENU).Count()
	error_log(err)
	if (listCol != 0) {
		err = s.DB(DBMNG).C(COLMENU).DropCollection()
		error_log(err)
	}
}

func saveOrUpdateProduct(prod Product) {
	s := createMongoSession()
	defer s.Close()
	listRes, err := s.DB(DBMNG).C(COLPRICES).Find(bson.M{"articul" : prod.Articul, "manufactura" : prod.Manufactura}).Count()
	error_log(err)
	if (listRes == 0) {
		err = s.DB(DBMNG).C(COLPRICES).Insert(prod)
		error_log(err)
	}else {
		var productFromBase Product
		err = s.DB(DBMNG).C(COLPRICES).Find(bson.M{"articul" : prod.Articul, "manufactura" : prod.Manufactura}).One(&productFromBase)
		productFromBase.Price = prod.Price
		err = s.DB(DBMNG).C(COLPRICES).Update(bson.M{"articul":prod.Articul, "manufactura":prod.Manufactura}, productFromBase)
		error_log(err)
	}
}

func addMenu(menu Punct) {
	s := createMongoSession()
	defer s.Close()
	err := s.DB(DBMNG).C(COLMENU).Insert(menu)
	error_log(err)
}

func isNotNeedsTovar(tovar Tovar, tovarFromPrice Product) bool {
	s := createMongoSession()
	defer s.Close()
	cnt, err := s.DB(DBMNG).C(COLNOTNEED).Find(bson.M{"articul":prepareStrings(tovar.Art), "provider" : tovarFromPrice.Manufactura}).Count()
	error_log(err)
	if(cnt > 0){
		return false
	}else{
		return true
	}
}

func findFromParser(tovarFromPrice Product) []Tovar {
	var arrayTovarsToSave []Tovar
	s := createMongoSession()
	defer s.Close()
	var currentTovar Tovar
	iterator := s.DB(DBMNG_PARSE).C(COLTOYS).Find(bson.M{"articul":bson.M{"$regex":".*"+prepareStrings(tovarFromPrice.Articul)+".*"}}).Iter()
	for iterator.Next(&currentTovar){
		if(isNotNeedsTovar(currentTovar, tovarFromPrice)){
			arrayTovarsToSave = append(arrayTovarsToSave, currentTovar)
		}
	}
	return arrayTovarsToSave
}

func updatePrice(){
	var (tovars []Tovar
		value Tovar
		curTovar TovarBitrix
		)
	s := createMongoSession()
	defer s.Close()
	var currentProduct Product
	iterator := s.DB(DBMNG).C(COLPRICES).Find(bson.M{}).Iter()
	for iterator.Next(&currentProduct){
		cnt, err := s.DB(DBMNG_PRICE_BITRIX).C(COLPRICES_BITRIX).Find(bson.M{"art" : currentProduct.Articul, "provider" : currentProduct.Manufactura}).Count()
		if(cnt > 0){
			err = s.DB(DBMNG_PRICE_BITRIX).C(COLPRICES_BITRIX).Find(bson.M{"art" : currentProduct.Articul, "provider" : currentProduct.Manufactura}).One(&curTovar)
			error_log(err)
			curTovar.Price = prepareStrings(currentProduct.Price)
			err = s.DB(DBMNG_PRICE_BITRIX).C(COLPRICES_BITRIX).Update(bson.M{"art" : currentProduct.Articul, "provider" : currentProduct.Manufactura},curTovar)
			error_log(err)
		}else{
			tovars = findFromParser(currentProduct)
			for  _, value = range tovars {
				curTovar.Art = prepareStrings(value.Art)
				curTovar.NameProduct = prepareStrings(value.NameProduct)
				curTovar.Action_acia = prepareStrings(value.Action_acia)
				curTovar.Action_new  = prepareStrings(value.Action_new)
				curTovar.Price       = prepareStrings(currentProduct.Price)
				curTovar.Sex         = prepareStrings(value.Sex)
				curTovar.Age         = prepareStrings(value.Age)
				curTovar.Descrip     = prepareStrings(value.Descrip)
				curTovar.Pic         = value.Pic
				curTovar.Provider    = prepareStrings(currentProduct.Manufactura)
				curTovar.Menu        = prepareStrings(currentProduct.Menu)
				err = s.DB(DBMNG_PRICE_BITRIX).C(COLPRICES_BITRIX).Insert(curTovar)
				error_log(err)
			}
		}
		fmt.Printf("%+v\n",currentProduct)
	}
}




