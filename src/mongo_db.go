package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	DBMNG string     = "parseprice"
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
		var p Product
		err = s.DB(DBMNG).C(COLPRICES).Find(bson.M{"articul" : prod.Articul, "manufactura" : prod.Manufactura}).One(&p)
		p.Price = prod.Price
		err = s.DB(DBMNG).C(COLPRICES).Update(bson.M{"articul":prod.Articul, "manufactura":prod.Manufactura},p)
		error_log(err)
	}
}

func addMenu(menu Punct) {
	s := createMongoSession()
	defer s.Close()
	err := s.DB(DBMNG).C(COLMENU).Insert(menu)
	error_log(err)
}
