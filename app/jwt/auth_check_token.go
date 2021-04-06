package jwt

import (
	"errors"
	sa "github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

// AuthCheckToken проверка авторизации по токену jwt
func (j *JWT) AuthCheckToken(ctx *sa.RequestCtx) error {
	// пропускаем точку авторизации.
	if string(ctx.Path()) == "/login" {
		return ctx.Next()
	}

	jwtCookie := ctx.Request.Header.Cookie("atreugo_jwt")

	if len(jwtCookie) == 0 {
		return ctx.ErrorResponse(errors.New("login required"), fasthttp.StatusForbidden)
	}

	token, err := j.ValidateToken(string(jwtCookie))
	if err != nil {
		return err
	}

	if !token.Valid {
		return ctx.ErrorResponse(errors.New("your session is expired, login again please"), fasthttp.StatusForbidden)
	}

	return ctx.Next()
}
