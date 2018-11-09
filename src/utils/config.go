package utils

import (
	"github.com/xormplus/xorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var Config *config

type config struct {
	Name        string 	`yaml:"name"`
	HttpPort    string 	`yaml:"http-port"`
	TablePrefix string 	`yaml:"table-prefix"`
	Db          []Db   	`yaml:"db"`
	SecretKey	string 	`yaml:"secret-key"`
	Expire		int64 	`yaml:"expire"`
}

type Db struct {
	Name         string `yaml:"name"`
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	Log          string `yaml:"log"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	ShowSql      bool   `yaml:"show-sql"`
}

func NewConfig() *config {
	return &config{}
}

func init() {
	Config = NewConfig()
	file, err := ioutil.ReadFile("config/app.yaml")
	if err != nil {
		log.Fatalf("config.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
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
