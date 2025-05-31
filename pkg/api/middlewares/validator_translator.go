package middlewares

import (
	"strings"

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

var universalTranslator *ut.UniversalTranslator

func ValidatorTranslator(c *gin.Context) {
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		acceptLanguage = "en;q=1.0"
	}

	parseAcceptLanguage, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	acceptLanguage = strings.Replace(parseAcceptLanguage[0].String(), "-", "_", -1)

	translator, _ := universalTranslator.GetTranslator(acceptLanguage)
	c.Set(apicontext.ValidatorTranslatorKey, translator)

	c.Next()
}

func init() {
	enLocale := en.New()
	ptBRLocale := pt_BR.New()

	universalTranslator = ut.New(
		enLocale, // Fallback
		enLocale,
		ptBRLocale,
	)

	enTranslator, _ := universalTranslator.GetTranslator(enLocale.Locale())
	ptBRTranslator, _ := universalTranslator.GetTranslator(ptBRLocale.Locale())

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = en_translation.RegisterDefaultTranslations(val, enTranslator)
		_ = pt_br_translation.RegisterDefaultTranslations(val, ptBRTranslator)
	}
}
