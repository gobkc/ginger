package ginger

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//签名算法列表
const (
	HS256 = 1 << iota
	HS384
	HS512
	ES256
	ES384
	ES512
)

const SECRET = "AF9-C=AF,FJN+RVV(DDD(M1a" //混淆代码 盐
const EXPIRED = 3600                      //过期世间
const SIGNING = HS512                     //签名算法

//错误相关常量
const ErrBusy = "程序忙"
const ErrorReLogin = "登录过期或登录错误，请重新登录"

//检查JWT信息是否正确
func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if tmp, ok := c.Request.Header["Authorization"]; ok && len(tmp)>=7 {
			token = tmp[0][7:]
		}
		claim, err := Verify(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "token不正确或已过期,详细错误：" + err.Error(),
			})
			return
		}
		c.Set("user",claim)
		c.Set("token",token)
		/*处理请求*/
		c.Next()
	}
}

//生成签名，并携带用户信息
func SignWithInfo(claims *JWTClaims) (token string, err error) {
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(EXPIRED)).Unix()
	token, err = MakeToken(claims)

	return token, err
}

//常量转换为固定的算法
func ToSignMethod(methodInt int) (sm jwt.SigningMethod) {
	switch methodInt {
	case ES256:
		sm = jwt.SigningMethodES256
	case ES384:
		sm = jwt.SigningMethodES384
	case ES512:
		sm = jwt.SigningMethodES512
	case HS256:
		sm = jwt.SigningMethodHS256
	case HS384:
		sm = jwt.SigningMethodHS384
	case HS512:
		sm = jwt.SigningMethodHS512
	}
	return sm
}

//定义一个结构体，用于携带额外的信息
type JWTClaims struct {
	jwt.StandardClaims
	UserID     uint    `json:"user_id"`
	Username   string `json:"username"`
	IsRealName int    `json:"is_real_name"`
	Type       int    `json:"type"`
}

//生成token
func MakeToken(claims *JWTClaims) (signedToken string, err error) {
	token := jwt.NewWithClaims(ToSignMethod(SIGNING), claims)
	if signedToken, err = token.SignedString([]byte(SECRET)); err != nil {
		return "", errors.New(ErrBusy)
	}
	return signedToken, err
}

//验证TOKEN是否正确或则过期
func Verify(token string) (claims *JWTClaims, err error) {
	var parseToken *jwt.Token
	parseToken, err = jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, errors.New(ErrBusy)
	}

	claims, ok := parseToken.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReLogin)
	}
	if err := parseToken.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReLogin)
	}

	return claims, err
}

//刷新
func Refresh(token string) (newToken string,err error) {
	var claims *JWTClaims
	if claims, err = Verify(token); err != nil {
		return newToken,err
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	if newToken, err = MakeToken(claims); err != nil {
		return newToken,err
	}
	log.Println("用户:", claims.Username, " token:", newToken, " 已刷新 ", " 过期时间：", claims.ExpiresAt)
	return newToken,err
}
