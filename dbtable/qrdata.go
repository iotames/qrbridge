package dbtable

// 二维码数据表结构
type Qrcode struct {
	BaseModel
	Code   string `xorm:"varchar(64) notnull unique 'code'"`
	ToURL  string `json:"to_url" xorm:"varchar(255) notnull 'to_url'"`
	Pv     int    `xorm:"int notnull default(0) 'pv'"`
	Status int16  `xorm:"smallint notnull default(10) 'status'"`
}

func (t *Qrcode) TableName() string {
	return "st_qrcode_list"
}

// 二维码查询记录表结构
type QrcodeQueryLog struct {
	BaseModel
	Code           string `xorm:"varchar(64) notnull 'code'"`
	UserAgent      string `json:"user_agent" xorm:"varchar(255) default('') 'user_agent'"`
	RequestHeaders string `json:"request_headers" xorm:"text notnull default('') 'request_headers'"`
	RequestIP      string `json:"request_ip" xorm:"varchar(45) default('') 'request_ip'"`
}

func (t *QrcodeQueryLog) TableName() string {
	return "st_qrcode_query_log"
}
