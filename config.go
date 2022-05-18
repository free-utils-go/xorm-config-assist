package xorm_config_assist

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Pwd          string
	Name         string
	MaxIdleCon   int
	MaxOpenCon   int
	Driver       string
	File         string
	Dsn          string
	ShowSQL      bool
	ShowExecTime bool
}

func InitDatabaseConfig(db_name string, tmp_dir string) {
	viper.SetDefault("database."+db_name+".host", "")
	viper.SetDefault("database."+db_name+".port", "")
	viper.SetDefault("database."+db_name+".user", "")
	viper.SetDefault("database."+db_name+".pwd", "")
	viper.SetDefault("database."+db_name+".name", db_name)
	viper.SetDefault("database."+db_name+".driver", "sqlite3")
	viper.SetDefault("database."+db_name+".file", tmp_dir+db_name+".db")
	viper.SetDefault("database."+db_name+".dsn", "")
	viper.SetDefault("database."+db_name+".max_idle_con", 1)
	viper.SetDefault("database."+db_name+".max_open_con", 1)
	viper.SetDefault("database."+db_name+".show_sql", false)
	viper.SetDefault("database."+db_name+".show_exec_time", true)
}

func LoadDatabaseConfig(db_name string) (db_config *DatabaseConfig) {
	db_config = &DatabaseConfig{
		Host:         viper.GetString("database." + db_name + ".host"),
		Port:         viper.GetInt("database." + db_name + ".port"),
		User:         viper.GetString("database." + db_name + ".user"),
		Pwd:          viper.GetString("database." + db_name + ".pwd"),
		Name:         viper.GetString("database." + db_name + ".name"),
		MaxIdleCon:   viper.GetInt("database." + db_name + ".max_idle_con"),
		MaxOpenCon:   viper.GetInt("database." + db_name + ".max_open_con"),
		Driver:       viper.GetString("database." + db_name + ".driver"),
		File:         viper.GetString("database." + db_name + ".file"),
		Dsn:          viper.GetString("database." + db_name + ".dsn"),
		ShowSQL:      viper.GetBool("database." + db_name + ".show_sql"),
		ShowExecTime: viper.GetBool("database." + db_name + ".show_exec_time"),
	}

	return db_config
}

// InitSqlite3 ...
func initSqlite3(sqfile string) *xorm.Engine {
	eng, err := xorm.NewEngine("sqlite3", sqfile)
	if err != nil {
		panic(err)
	}

	_, err = eng.Exec("PRAGMA journal_mode = OFF;")
	if err != nil {
		panic(err)
	}
	return eng
}

const sqlURL = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

func initMysql(addr, user_name, pass, db_name string) *xorm.Engine {
	data_source := fmt.Sprintf(sqlURL, user_name, pass, addr, db_name, url.QueryEscape("Asia/Shanghai"))
	before := time.Now()
	eng, err := xorm.NewEngine("mysql", data_source)
	if err != nil {
		panic(err)
	}

	fmt.Printf("took %v\n", time.Since(before))
	return eng
}

func InitXorm(cfg *DatabaseConfig) (eng *xorm.Engine) {
	cfg_driver := cfg.Driver
	if cfg_driver == "sqlite3" {
		eng = initSqlite3(cfg.File)
	} else if cfg_driver == "mysql" {
		eng = initMysql(cfg.Host, cfg.User, cfg.Pwd, cfg.Name)
	} else {
		panic("数据库类型不支持")
	}

	eng.ShowSQL(cfg.ShowSQL)
	return eng
}
