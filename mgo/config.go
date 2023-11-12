package mgo

var (
	C = new(Config)
)

// Config 配置参数
type Config struct {
	Mongodb MongoConfig
}

type MongoConfig struct {
	// mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority"
	Host       string `json:"host"`                                              // 地址
	Port       string `json:"port" default:"27017"`                              // 端口
	Database   string `json:"database"`                                          // 数据库
	Username   string `json:"username"`                                          // 用户
	Password   string `json:"password"`                                          // 密码
	RawOptions string `json:"raw_options" default:"w=majority&authSource=admin"` // w=majority { w: <value>, j: <boolean>, wtimeout: <number> }
}
