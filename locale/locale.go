package locale

import (
	"encoding/json"

	ginI18n "github.com/gin-contrib/i18n"
	"golang.org/x/text/language"
)

func NewI18nCfg() {
	ginI18n.Localize(ginI18n.WithBundle(&ginI18n.BundleCfg{
		RootPath:         "./locale/",
		AcceptLanguage:   []language.Tag{language.German, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}))
}
