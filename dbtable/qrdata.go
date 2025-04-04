package dbtable

// 二维码数据表结构
type QRCodeData struct {
	BaseModel `xorm:"extends"`
	QRCode    string `xorm:"varchar(64) notnull unique 'qrcode'"`
	ToURL     string `xorm:"varchar(255) notnull 'to_url'"`
	Pv        uint   `xorm:"int unsigned default(0) 'pv' comment('总访问量')"`
	Status    int16  `xorm:"'status' smallint notnull default(10)"`
}

// 二维码查询记录表结构
type QRCodeQueryLog struct {
	BaseModel      `xorm:"extends"`
	QRCode         string `xorm:"varchar(64) notnull unique 'qrcode'"`
	UserAgent      string `xorm:"varchar(255) notnull 'user_agent'"`
	RequestHeaders string `xorm:"text notnull 'request_headers'"`
	RequestIP      string `xorm:"varchar(30) notnull 'request_ip'"`
}
