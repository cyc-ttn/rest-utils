package security

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("Your token is invalid") //When a token is invalid
	ErrExpiredToken = errors.New("Your token is expired") //When a token is expired
)

// UserToken stores details about the user token
type UserToken struct {
	User  UserDetails
	Agent string
}

// UserTokenRepository keeps a repository of all tokens
type UserTokenRepository struct {
	key    string
	tokens map[string]*UserToken
}

// NewUserTokenRepository creates a new repository
func NewUserTokenRepository(signKey string) *UserTokenRepository {
	return &UserTokenRepository{
		key:    signKey,
		tokens: make(map[string]*UserToken),
	}
}

// GenerateToken generates a token based on the request
func (r *UserTokenRepository) GenerateToken(req *http.Request) (string, error) {

	usr := req.Context().Value(PrincipalKey).(UserDetails)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiry":  time.Now().Add(time.Hour * 5).Unix(),
		"created": time.Now().Unix(),
	})

	str, err := token.SignedString([]byte(r.key))
	if err != nil {
		return "", err
	}

	t := &UserToken{
		User:  usr,
		Agent: req.Header.Get("User-Agent"),
	}

	r.tokens[str] = t

	return str, nil
}

// FindToken returns the token matching the string
func (r *UserTokenRepository) FindToken(str string) *UserToken {
	token, ok := r.tokens[str]
	if !ok {
		return nil
	}
	return token
}

// FindAndVerifyToken verify token and string
func (r *UserTokenRepository) FindAndVerifyToken(str string) (*UserToken, error) {
	token := r.FindToken(str)
	if token == nil {
		return nil, ErrInvalidToken
	}

	_, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.key), nil
	})

	if err == nil {
		return token, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		// if ve.Errors & jwt.ValidationErrorMalformed != 0 {
		//   return nil, ErrInvalidToken
		// }
		//
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrExpiredToken
		}
	}

	return nil, ErrInvalidToken
}
