package security

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	//ErrInvalidToken hen a token is invalid
	ErrInvalidToken = errors.New("Your token is invalid")

 	//ErrExpiredToken When a token is expired
 	ErrExpiredToken = errors.New("Your token is expired")
)

// UserToken stores details about the user token
type UserToken struct {
	User  			UserDetails
	Agent 			string
	Token				*jwt.Token
	SignedToken string
}

// UserTokenRepository keeps a repository of all tokens
type UserTokenRepository struct {
	key    string
	tokens map[string]*UserToken
	tokensByUser map[uint][]*UserToken
	mu		* sync.RWMutex
}

// NewUserTokenRepository creates a new repository
func NewUserTokenRepository(signKey string) *UserTokenRepository {
	return &UserTokenRepository{
		key:    signKey,
		tokens: make(map[string]*UserToken),
		tokensByUser: make(map[uint][]*UserToken),
		mu: &sync.RWMutex{},
	}
}

// GenerateToken generates a token based on the request
// Uses default of 5 hours
func (r *UserTokenRepository) GenerateToken(req *http.Request) (string, error) {
	return r.GenerateTokenWithExpiry(req, time.Now().Add(time.Hour * 5).Unix() )
}

// GenerateTokenWithExpiry generates a token based on the request and an expiry
func (r *UserTokenRepository) GenerateTokenWithExpiry(req *http.Request, expiry int64) (string, error) {

	allDetails := req.Context().Value(PrincipalKey).(map[string]interface{})
	usr := allDetails["user"].(UserDetails)

	t := &UserToken{
		User:  usr,
		Agent: req.Header.Get("User-Agent"),
	}

	// Check if there is an existing token
	if str := r.findExistingToken(t); str != "" {
		return str, nil
	}

	// Create a new token
	t.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiry" : expiry,
		"created": time.Now().Unix(),
	})

	str, err := t.Token.SignedString([]byte(r.key))
	if err != nil {
		return "", err
	}

	t.SignedToken = str

	r.mu.Lock()
	r.tokens[str] = t
	r.mu.Unlock()

	r.addToUserCache( usr, t)

	go func(){
		select {
		case <- time.After( time.Unix(expiry,0).Sub(time.Now()) ):
			r.RemoveToken(str)
		}
	}()

	return str, nil
}

// RemoveToken removes a token from the repository
func (r *UserTokenRepository) RemoveToken(str string) {

	r.mu.RLock()
	token := r.tokens[str]
	r.mu.RUnlock()

	if token == nil { return }

	r.removeTokenFromUserCache( token.User, token.Token)

	r.mu.Lock()
	delete(r.tokens, str)
	r.mu.Unlock()
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
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, ErrExpiredToken
		}
	}

	return nil, ErrInvalidToken
}

func (r * UserTokenRepository) addToUserCache( usr UserDetails, t * UserToken ){

	tokens := r.FindTokensByUser(usr)
	if tokens == nil {
		tokens = make([]*UserToken, 0, 10)
	}

	length := len(tokens)
	if len(tokens) == cap(tokens) {
		newToken := make([]*UserToken, length, length + 10)
		copy(newToken, tokens)
		tokens = newToken
	}

	tokens = append( tokens, t)

	r.mu.Lock()
	r.tokensByUser[ usr.GetID() ] = tokens
	r.mu.Unlock()
}

// FindTokensByUser -
func (r * UserTokenRepository) FindTokensByUser( usr UserDetails ) []*UserToken{
	id := usr.GetID()

	r.mu.RLock()
	tokens := r.tokensByUser[ id ]
	r.mu.RUnlock()

	return tokens
}

func (r * UserTokenRepository) removeTokenFromUserCache( usr UserDetails, t * jwt.Token ){

	tokens := r.FindTokensByUser(usr)
	for i, token := range tokens{
		if token.Token == t {
			tokens = append(tokens[:i],tokens[i+1:]...)
			break
		}
	}

	r.mu.RLock()
	r.tokensByUser[ usr.GetID() ] = tokens
	r.mu.RUnlock()
}

func (r * UserTokenRepository) findExistingToken( t * UserToken ) string{

	tokens := r.FindTokensByUser( t.User )
	for _, token := range tokens {
		if strings.Compare(token.Agent , t.Agent ) == 0 {
			return token.SignedToken
		}
	}

	return ""
}
