package utils

import (
	"github.com/xormplus/xorm"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-log"
	"encoding/json"
)

var Config *config

type data struct {
	Micro	*micro		`json:"micro"`
}

type micro struct {
	Conf		*conf	`json:"config"`
}

type conf struct {
	Config 		*config	`json:"sso"`
}

type config struct {
	Name        string 	`json:"name"`
	HttpPort    string 	`json:"http-port"`
	TablePrefix string 	`json:"table-prefix"`
	Db          []Db   	`json:"db"`
	SecretKey	string 	`json:"secret-key"`
	Expire		int64 	`json:"expire"`
}

type Db struct {
	Name         string `json:"name"`
	Driver       string `json:"driver"`
	Dsn          string `json:"dsn"`
	Log          string `json:"log"`
	MaxIdleConns int    `json:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns"`
	ShowSql      bool   `json:"show-sql"`
}

func NewConfig() *config {
	return &config{}
}

func init() {
	//Config = NewConfig()
	//file, err := ioutil.ReadFile("config/app.yaml")
	//if err != nil {
	//	log.Fatalf("config.Get err   #%v ", err)
	//}
	//err = yaml.Unmarshal(file, &Config)
	//if err != nil {
	//	log.Fatalf("Unmarshal: %v", err)
	//}
	consulSource := consul.NewSource(
		// optionally specify consul address; default to localhost:8500
		consul.WithAddress("localhost:8500"),
		// optionally specify prefix; defaults to /micro/config
		consul.WithPrefix("/micro/config/sso"),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
	)
	ChangeSet, err := consulSource.Read()
	if err != nil {
		log.Log(err)
	}
	var d data
	json.Unmarshal(ChangeSet.Data, &d)
	Config = d.Micro.Conf.Config
}

func (db *Db) GetEngin() (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine(db.Driver, db.Dsn)
	if err != nil {
		return
	}
	if db.Log != "" {
		var f *os.File
		f, err = os.Create(db.Log)
		if err != nil {
			return
		}
		engine.SetLogger(xorm.NewSimpleLogger(f))
		engine.ShowSQL(true)
	}
	if db.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(db.MaxIdleConns)
	}
	if db.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(db.MaxOpenConns)
	}
	return
}
