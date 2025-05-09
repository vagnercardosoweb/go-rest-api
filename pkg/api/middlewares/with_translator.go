package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	ptbrtranslations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"golang.org/x/text/language"
)

func WithTranslator(c *gin.Context) {
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

		if acceptLanguage == language.BrazilianPortuguese.String() {
			_ = ptbrtranslations.RegisterDefaultTranslations(val, translator)
		} else {
			_ = entranslations.RegisterDefaultTranslations(val, translator)
		}
	}

	c.Set(utils.ValidateTranslatorCtxKey, &translator)

	c.Next()
}
