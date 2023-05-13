package config

type ServerCfg struct {
	Addr           string `toml:"address"`
	Port           int    `toml:"port"`
	Secret         string `toml:"secret"`
	EnableRegister bool   `toml:"enable-register"`
}

type ReCAPTCHACfg struct {
	Host   string `toml:"server"`
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
	ApiKey string `toml:"apikey"`
	Sender string `toml:"sender"`
}

type ChatGPTCfg struct {
	ReverseProxyUrl string `toml:"url"`
	TimeoutSec      uint   `toml:"timeout"` // unit: second
	ChatGPTToken    string `toml:"token"`
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

type LogCfg struct {
	Level  string `toml:"level"`  // panic, fatal, error, warn, warning, info, debug, trace
	Output string `toml:"output"` // stdout, stderr, [filename]
	Format string `toml:"format"` // json, text
}

type Cfg struct {
	Log        LogCfg        `toml:"debug"`
	Server     ServerCfg     `toml:"http"`
	Code       ReCAPTCHACfg  `toml:"code"`
	DB         DBCfg         `toml:"database"`
	Redis      RedisCfg      `toml:"redis"`
	Mail       MailCfg       `toml:"mail"`
	ChatGPT    ChatGPTCfg    `toml:"chatgpt"`
	ChatgptWeb ChatgptWebCfg `toml:"chatgpt-web"`
}

var GlobalCfg = Cfg{
	Log: LogCfg{
		Level:  "info",
		Output: "stderr",
	},
	Server: ServerCfg{
		Addr:           "127.0.0.1",
		Port:           8080,
		Secret:         "",
		EnableRegister: true,
	},
	Code: ReCAPTCHACfg{
		// https://developers.google.com/recaptcha/docs/faq?hl=zh-cn
		Host:   "www.recaptcha.net",
		Secret: "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe",
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
		ApiKey: "",
		Sender: "noreply@cheatppt.icu",
	},
	ChatGPT: ChatGPTCfg{
		ReverseProxyUrl: "",
		TimeoutSec:      60,
		ChatGPTToken:    "",
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
var LogOpts = &GlobalCfg.Log

var Server = &GlobalCfg.Server
var Code = &GlobalCfg.Code
var Mail = &GlobalCfg.Mail
var ChatGPT = &GlobalCfg.ChatGPT
