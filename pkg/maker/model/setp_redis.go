package model

type StepRedisModel struct {
	*StepModel `json:",inline"`

	Redis     string `json:"redis,omitempty"`     // redis操作
	Key       string `json:"key,omitempty"`       // key，用于redis、cache、MQ等操作
	KeyType   string `json:"keyType,omitempty"`   // 值key类型
	Value     string `json:"value,omitempty"`     // value，用于redis、cache、MQ等操作
	ValueType string `json:"valueType,omitempty"` // 值类型
	Expire    string `json:"expire,omitempty"`    // 过期时间
}

var (
	docTemplateStepRedisName = "step_redis"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepRedisName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "redis",
				Comment: "redis操作",
			},
			{
				Name:    "key",
				Comment: "操作的Key",
			},
			{
				Name:    "value",
				Comment: "操作的Value",
			},
		},
	})
}
