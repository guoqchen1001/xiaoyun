package root

// MssqlConfig mssql数据库连接参数
type MssqlConfig struct {
	Dialects string `json:"dialects"`
	Parm     string `json:"parm"`
}

// MysqlConfig mysql数据库连接
type MysqlConfig struct {
	Dialects string `json:"dialects"`
	Parm     string `json:"parm"`
}

// HTTPConfig http服务设置
type HTTPConfig struct {
	Host string `json:"host"`
}

// AuthConfig 认证设置
type AuthConfig struct {
	Secret    string `json:"secret"`
	ExpiredAt int    `json:"expiredAt"`
}

//PathConfig 路径设置
type PathConfig struct {
	ConfigPath string `json:"configPath"`
	LogPath    string `json:"logPath"`
}

// Config 设置
type Config struct {
	Mssql *MssqlConfig `json:"mssql"`
	Mysql *MysqlConfig `json:"mysql"`
	HTTP  *HTTPConfig  `json:"http"`
	Auth  *AuthConfig  `json:"auth"`
	Path  *PathConfig  `json:"-"`
}

// Configer 获取系统配置信息
type Configer interface {
	GetConfig() (*Config, error)
}
