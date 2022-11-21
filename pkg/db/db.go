package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-driver/db_dm"
	"github.com/team-ide/go-driver/db_kingbase_v8r6"
	"github.com/team-ide/go-driver/db_mysql"
	"github.com/team-ide/go-driver/db_sqlite3"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/util"
)

type DatabaseType struct {
	DialectName string `json:"dialectName"`
	newDb       func(config *DatabaseConfig) (db *sql.DB, err error)
	dia         dialect.Dialect
	matches     []string
}

func (this_ *DatabaseType) init() {
	var err error
	this_.dia, err = dialect.NewDialect(this_.DialectName)
	if err != nil {
		panic(err)
	}
	return
}

var (
	DatabaseTypes []*DatabaseType
)

func init() {
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_mysql.GetDSN(config.Username, config.Password, config.Host, config.Port, config.Database)
			db, err = db_mysql.Open(dsn)
			return
		},
		DialectName: db_mysql.GetDialect(),
		matches:     []string{"mysql"},
	})

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_sqlite3.GetDSN(config.DatabasePath)
			db, err = db_sqlite3.Open(dsn)
			return
		},
		DialectName: db_sqlite3.GetDialect(),
		matches:     []string{"sqlite", "sqlite3"},
	})

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_dm.GetDSN(config.Username, config.Password, config.Host, config.Port)
			db, err = db_dm.Open(dsn)
			return
		},
		DialectName: db_dm.GetDialect(),
		matches:     []string{"DaMeng", "dm"},
	})
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_kingbase_v8r6.GetDSN(config.Username, config.Password, config.Host, config.Port, config.DbName)
			db, err = db_kingbase_v8r6.Open(dsn)
			return
		},
		DialectName: db_kingbase_v8r6.GetDialect(),
		matches:     []string{"KingBase", "kb"},
	})

	initOracleDatabase()
	initShenTongDatabase()
}

func addDatabaseType(databaseType *DatabaseType) *DatabaseType {
	databaseType.init()
	DatabaseTypes = append(DatabaseTypes, databaseType)
	return databaseType
}

func GetDatabaseType(databaseType string) *DatabaseType {
	for _, one := range DatabaseTypes {
		if strings.EqualFold(databaseType, one.DialectName) {
			return one
		}
		for _, match := range one.matches {
			if strings.EqualFold(databaseType, match) {
				return one
			}
		}
	}
	return nil
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string `json:"type,omitempty"`
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	Database     string `json:"database,omitempty"`
	DbName       string `json:"dbName,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Sid          string `json:"sid,omitempty"`
	MaxIdleConns int    `json:"maxIdleConns,omitempty"`
	MaxOpenConns int    `json:"maxOpenConns,omitempty"`
	DatabasePath string `json:"databasePath,omitempty"`
}

// NewDatabaseWorker 根据数据库配置创建DatabaseWorker
func NewDatabaseWorker(config *DatabaseConfig) (databaseWorker *DatabaseWorker, err error) {
	databaseWorker = &DatabaseWorker{config: config}
	err = databaseWorker.init()
	if err != nil {
		return nil, err
	}
	return databaseWorker, nil
}

// DatabaseWorker 基础操作
type DatabaseWorker struct {
	config       *DatabaseConfig
	databaseType *DatabaseType
	db           *sql.DB
	dialect.Dialect
}

func (this_ *DatabaseWorker) GetDialectName() string {
	return this_.databaseType.DialectName
}

func (this_ *DatabaseWorker) init() (err error) {
	this_.databaseType = GetDatabaseType(this_.config.Type)
	if this_.databaseType == nil {
		err = errors.New("数据库类型[" + this_.config.Type + "]暂不支持")
		return
	}

	this_.Dialect = this_.databaseType.dia
	this_.db, err = this_.databaseType.newDb(this_.config)
	if err != nil {
		return
	}

	if this_.config.MaxIdleConns > 0 {
		this_.db.SetMaxIdleConns(this_.config.MaxIdleConns)
	}
	if this_.config.MaxOpenConns > 0 {
		this_.db.SetMaxOpenConns(this_.config.MaxOpenConns)
	}

	err = this_.db.Ping()
	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseWorker) GetConfig() (config *DatabaseConfig) {
	config = this_.config
	return
}

func (this_ *DatabaseWorker) Open() (err error) {
	err = this_.db.Ping()
	return
}

func (this_ *DatabaseWorker) Close() (err error) {
	err = this_.db.Close()
	return
}

func (this_ *DatabaseWorker) Exec(sql string, args []interface{}) (rowsAffected int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			err = errors.New("Exec error sql:" + sql + ",error:" + err.Error())
		}
	}()

	rowsAffected, err = this_.Execs([]string{sql}, [][]interface{}{args})

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseWorker) Execs(sqlList []string, argsList [][]interface{}) (rowsAffected int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("Execs Error", zap.Any("sqlList", sqlList), zap.Any("argsList", argsList), zap.Error(err))
			err = errors.New("Execs error sql:" + strings.Join(sqlList, ";") + ",error:" + err.Error())
		}
	}()
	res, errSql, errArgs, err := worker.DoExecs(this_.db, sqlList, argsList)
	if err != nil {
		util.Logger.Error("Execs Error", zap.Any("sql", errSql), zap.Any("args", errArgs), zap.Error(err))
		return
	}
	for _, one := range res {
		rowsAffected_, _ := one.RowsAffected()
		rowsAffected += rowsAffected_
	}
	return
}

func (this_ *DatabaseWorker) Count(sql string, args []interface{}) (count int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("Count Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("Count error sql:" + sql + ",error:" + err.Error())
		}
	}()
	count_, err := worker.DoQueryCount(this_.db, sql, args)
	if err != nil {
		util.Logger.Error("Count Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}
	count = int64(count_)
	return
}

func (this_ *DatabaseWorker) QueryOne(sql string, args []interface{}, one interface{}) (find bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("QueryOne Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("QueryOne error sql:" + sql + ",error:" + err.Error())
		}
	}()
	find, err = worker.DoQueryStruct(this_.db, sql, args, one)

	if err != nil {
		util.Logger.Error("QueryOne Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) Query(sql string, args []interface{}, list interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("Query error sql:" + sql + ",error:" + err.Error())
		}
	}()
	err = worker.DoQueryStructs(this_.db, sql, args, list)

	if err != nil {
		util.Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryMap(sql string, args []interface{}) (list []map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("QueryMap error sql:" + sql + ",error:" + err.Error())
		}
	}()
	list, err = worker.DoQuery(this_.db, sql, args)

	if err != nil {
		util.Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryPage(sql string, args []interface{}, list interface{}, page *worker.Page) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("QueryPage Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("QueryPage error sql:" + sql + ",error:" + err.Error())
		}
	}()
	err = worker.DoQueryPageStructs(this_.db, this_.Dialect, sql, args, page, list)

	if err != nil {
		util.Logger.Error("QueryPage Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryMapPage(sql string, args []interface{}, page *worker.Page) (list []map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("QueryMapPage Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
			err = errors.New("QueryMapPage error sql:" + sql + ",error:" + err.Error())
		}
	}()
	list, err = worker.DoQueryPage(this_.db, this_.Dialect, sql, args, page)

	if err != nil {
		util.Logger.Error("QueryMapPage Error", zap.Any("sql", sql), zap.Any("args", args), zap.Error(err))
		return
	}

	return
}
