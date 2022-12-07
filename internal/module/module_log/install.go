package module_log

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建登录表
		{
			Version: "1.0",
			Module:  ModuleLog,
			Stage:   `创建表[` + TableLog + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableLog + ` (
	logId bigint(20) NOT NULL COMMENT '登录ID',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	action varchar(100) DEFAULT NULL COMMENT '操作',
	method varchar(20) DEFAULT NULL COMMENT '方法',
	param text DEFAULT NULL COMMENT '参数',
	data text DEFAULT NULL COMMENT '数据',
	userAgent text DEFAULT NULL COMMENT 'User-Agent',
	status int(2) NOT NULL DEFAULT 0 COMMENT '状态',
	error varchar(200) DEFAULT NULL COMMENT '异常',
	createTime datetime NOT NULL COMMENT '创建时间',
	startTime datetime DEFAULT NULL COMMENT '开始时间',
	endTime datetime DEFAULT NULL COMMENT '结束时间',
	useTime int(10) DEFAULT 0 COMMENT '使用时长',
	PRIMARY KEY (logId),
	KEY index_userId (userId),
	KEY index_ip (ip),
	KEY index_status (status),
	KEY index_action (action),
	KEY index_useTime (useTime),
	KEY index_startTime (startTime),
	KEY index_endTime (endTime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableLogComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableLog + ` (
	logId bigint(20) NOT NULL,
	userId bigint(20) DEFAULT NULL,
	ip varchar(50) DEFAULT NULL,
	action varchar(100) DEFAULT NULL,
	method varchar(20) DEFAULT NULL,
	param text DEFAULT NULL,
	data text DEFAULT NULL,
	userAgent text DEFAULT NULL,
	status int(2) NOT NULL DEFAULT 0,
	error varchar(200) DEFAULT NULL,
	createTime datetime NOT NULL,
	startTime datetime DEFAULT NULL,
	endTime datetime DEFAULT NULL,
	useTime int(10) DEFAULT 0,
	PRIMARY KEY (logId)
);
`,
					`CREATE INDEX ` + TableLog + `_index_userId on ` + TableLog + ` (userId);`,
					`CREATE INDEX ` + TableLog + `_index_ip on ` + TableLog + ` (ip);`,
					`CREATE INDEX ` + TableLog + `_index_status on ` + TableLog + ` (status);`,
					`CREATE INDEX ` + TableLog + `_index_action on ` + TableLog + ` (action);`,
					`CREATE INDEX ` + TableLog + `_index_useTime on ` + TableLog + ` (useTime);`,
					`CREATE INDEX ` + TableLog + `_index_startTime on ` + TableLog + ` (startTime);`,
					`CREATE INDEX ` + TableLog + `_index_endTime on ` + TableLog + ` (endTime);`,
				},
			},
		},
	}
}
