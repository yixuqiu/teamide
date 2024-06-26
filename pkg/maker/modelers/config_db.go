package modelers

type ConfigDbModel struct {
	ElementNode
	Comment         string `json:"comment,omitempty"` // 说明
	Note            string `json:"note,omitempty"`    // 注释
	Type            string `json:"type,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	Database        string `json:"database,omitempty"`
	DbName          string `json:"dbName,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	OdbcDsn         string `json:"odbcDsn,omitempty"`
	OdbcDialectName string `json:"odbcDialectName,omitempty"`

	Schema       string `json:"schema,omitempty"`
	Sid          string `json:"sid,omitempty"`
	MaxIdleConn  int    `json:"maxIdleConn,omitempty"`
	MaxOpenConn  int    `json:"maxOpenConn,omitempty"`
	DatabasePath string `json:"databasePath,omitempty"`
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigDbName,
		Comment: "Redis配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "type", Comment: ""},
			{Name: "host", Comment: ""},
			{Name: "port", Comment: ""},
			{Name: "database", Comment: ""},
			{Name: "dbName", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
			{Name: "odbcDsn", Comment: ""},
			{Name: "odbcDialectName", Comment: ""},
			{Name: "schema", Comment: ""},
			{Name: "sid", Comment: ""},
			{Name: "maxIdleConn", Comment: ""},
			{Name: "maxOpenConn", Comment: ""},
			{Name: "databasePath", Comment: ""},
		},
	})
}
