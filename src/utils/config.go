package utils

import (
	"github.com/xormplus/xorm"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-log"
	"encoding/json"
	"github.com/go-redis/redis"
)

var Config *config

type config struct {
	Name        string 	`json:"name"`
	HttpPort    string 	`json:"http-port"`
	TablePrefix string 	`json:"table-prefix"`
	Db          []Db   	`json:"db"`
	Redis 		[]Redis `json:"redis"`
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

type Redis struct {
	Name 		string	`json:"name"`
	Addr 		string	`json:"addr"`
	Password 	string	`json:"password"`
	Db 			int		`json:"db"`
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
		log.Fatal(err)
	}
	if ChangeSet == nil {
		log.Fatal("/micro/config/sso is nil")
	}
	var data map[string]map[string]map[string]*config
	if err := json.Unmarshal(ChangeSet.Data, &data); err != nil {
		log.Fatal(err)
	}
	if data["micro"]["config"]["sso"] != nil {
		Config = data["micro"]["config"]["sso"]
	}

	// 加载mysql数据库连接
	initOrm()

	// 加载redis数据库连接
	initRedis()

	// 加载验证码配置
	initCaptcha()
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

func (r *Redis) GetClient() (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr: r.Addr,
		Password: r.Password, // no password set
		DB: r.Db,  // use default DB
	})
	return client, nil
}
