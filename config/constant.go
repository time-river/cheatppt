package config

type ServerCfg struct {
	Addr   string `toml:"address"`
	Port   int    `toml:"port"`
	Secret string `toml:"secret"`
}

type ReCAPTCHACfg struct {
	Enable bool   `toml:"enable"`
	Server string `toml:"server"`
	Secret string `toml:"secret"`
}

type DBCfg struct {
	Driver string `toml:"driver"`
	Addr   string `toml:"address"`
	Passwd string `toml:"password"`
	DbName string `toml:"database"`
}

type RedisCfg struct {
	Addr   string `toml:"address"`
	Passwd string `toml:"password"`
	Db     int    `toml:"db"`
	Lease  int64  `toml:"lease"`
}

type MailCfg struct {
	Driver string `toml:"driver"`
	ApiKey string `toml:"apikey"`
}

type ChatgptWebCfg struct {
	ApiModel        string `toml:"api-model"`
	TimeoutMs       uint64 `toml:"timeout"`
	ReverseProxyUrl string `toml:"reverse-proxy-url"`
	SocksProxy      string `toml:"socks-proxy"`
	HttpsProxy      string `toml:"https-proxy"`
	OpenaiApiBase   string `toml:"openai-api-base"`
	OpenaiApiToken  string `toml:"openai-api-token"`
	AuthSecretKey   string `toml:"auth-secret-key"`
}

type Cfg struct {
	Debug      bool          `toml:"debug"`
	Server     ServerCfg     `toml:"http"`
	ReCAPTCHA  ReCAPTCHACfg  `toml:"reCAPTCHA"`
	DB         DBCfg         `toml:"database"`
	Redis      RedisCfg      `toml:"redis"`
	Mail       MailCfg       `toml:"mail"`
	ChatgptWeb ChatgptWebCfg `toml:"chatgpt-web"`
}

var GlobalCfg = Cfg{
	Debug: true,
	Server: ServerCfg{
		Addr:   "127.0.0.1",
		Port:   8080,
		Secret: "",
	},
	ReCAPTCHA: ReCAPTCHACfg{
		Enable: false,
		Server: "https://www.google.com/recaptcha/api/siteverify",
		Secret: "",
	},
	DB: DBCfg{
		Driver: "sqlite3",
		Addr:   "./cheatppt.db",
		Passwd: "",
		DbName: "cheatppt",
	},
	Redis: RedisCfg{
		Addr:   "127.0.0.1:6379",
		Passwd: "",
		Db:     0,
		Lease:  3600 * 24 * 3, /* 3 days */
	},
	Mail: MailCfg{
		Driver: "sendgrid",
		ApiKey: "",
	},
	ChatgptWeb: ChatgptWebCfg{
		ApiModel:        "text-davinci-002-render-sha",
		TimeoutMs:       0,
		ReverseProxyUrl: "https://openai.vvl.me:8443/chat/api/conversation",
		SocksProxy:      "-",
		HttpsProxy:      "-",
		OpenaiApiBase:   "https://api.openai.com",
		OpenaiApiToken:  "",
		AuthSecretKey:   "",
	},
}

var GlobalKey [32]byte
