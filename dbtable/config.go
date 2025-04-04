package dbtable

type Config struct {
	BaseModel   `xorm:"extends"`
	Key         string `xorm:"varchar(255) notnull unique 'key' comment('配置key')"`
	Value       string `xorm:"varchar(255) notnull 'value' comment('配置value')"`
	Description string `xorm:"varchar(255) notnull 'description' comment('配置描述')"`
}

// Title       string `xorm:"varchar(255) notnull unique 'title' comment('站点标题')"`
// Description string `xorm:"varchar(255) notnull 'description' comment('站点描述')"`
// Keywords    string `xorm:"varchar(255) notnull 'keywords' comment('站点关键字')"`
// Logo        string `xorm:"varchar(255) notnull 'logo' comment('站点logo')"`
// Copyright   string `xorm:"varchar(255) notnull 'copyright' comment('版权信息')"`
// IcpCode     string `xorm:"varchar(255) notnull 'icpcode' comment('ICP备案号')"`
