package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	pt_br_translation "github.com/go-playground/validator/v10/translations/pt_BR"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"golang.org/x/text/language"
)

func Translator(c *gin.Context) {
	var translator ut.Translator

	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		acceptLanguage = "en;q=1.0"
	}

	parseAcceptLanguage, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	acceptLanguage = parseAcceptLanguage[0].String()

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		translator, _ = ut.New(en.New(), pt_BR.New()).
			GetTranslator(acceptLanguage)

		switch acceptLanguage {
		case language.BrazilianPortuguese.String():
			_ = pt_br_translation.RegisterDefaultTranslations(val, translator)
		default:
			_ = en_translation.RegisterDefaultTranslations(val, translator)
		}
	}

	c.Set(apicontext.TranslatorKey, &translator)

	c.Next()
}
