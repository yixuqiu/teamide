package modelers

type ConfigRedisModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Auth     string `json:"auth,omitempty"`     //
	CertPath string `json:"certPath,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigRedisName,
		Comment: "Redis配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "auth", Comment: ""},
			{Name: "certPath", Comment: ""},
		},
	})
}
