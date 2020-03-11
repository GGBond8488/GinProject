package util
/*
	使用jwt做权限验证
*/
import (
	"My-gin-Project/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"time"
)


var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username,password string)(string,error)  {
	nowTime := time.Now()
	expireTime := nowTime.Add(3*time.Hour)

	claims:= Claim{
		Username:       username,
		Password:       password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:expireTime.Unix(),
			Issuer : "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token,err

}

func ParseToken(token string)(*Claim,error)  {
	tokenClaims,err := jwt.ParseWithClaims(token,&Claim{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret,err
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claim); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}