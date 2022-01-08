package users

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"task-scheduler/internal/platform/logger"

	"github.com/bnkamalesh/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID       int64      `json:"uid,omitempty"`
	Name      string     `json:"name,omitempty"`
	Password  string     `json:"password,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type Claims struct {
	Fullname string `json:"fullname"`
	jwt.StandardClaims
}

type JWT struct {
	Access_token string `json:"access_token,omitempty"`
	Token_type   string `json:"token_type,omitempty"`
	Expires_in   int64  `json:"expires_in,omitempty"`
}

func (u *User) init() {
	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &now
	}
}

func (u *User) Sanitize() {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
}

func (u *User) Validate() error {
	if u.Email == "" {
		return nil
	}

	err := validateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.Validation("invalid email address provided")
	}

	return nil
}

type Users struct {
	logHandler logger.Logger
	store      store
}

func (us *Users) Register(ctx context.Context, u *User) (*User, error) {
	u.init()
	u.Sanitize()
	err := u.Validate()
	if err != nil {
		return nil, err
	}

	hash, _ := HashPassword(u.Password)
	u.Password = hash
	err = us.store.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	u, err = us.GetUserByEmail(ctx, u.Email)
	u.Password = ""
	return u, nil
}

func CreateToken(u *User) (JWT, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Fullname: u.Name,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.Email,
			Id:        strconv.FormatInt(u.UID, 10),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret_key := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secret_key)
	jwt := JWT{}
	if err != nil {
		return jwt, err
	}
	jwt.Access_token = tokenString
	jwt.Token_type = "Bearer"
	jwt.Expires_in = expirationTime.Unix()
	return jwt, nil
}

func (us *Users) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.TrimSpace(email)
	err := validateEmail(email)
	if err != nil {
		return nil, err
	}

	u, err := us.store.GetUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *Users) Login(ctx context.Context, email string, password string) (JWT, error) {
	emptyJWT := JWT{}

	u, err := us.GetUserByEmail(ctx, email)
	if err != nil || !CheckPasswordHash(password, u.Password) {
		return emptyJWT, errors.Unauthorized("Wrong username or password")
	}
	token, err := CreateToken(u)
	if err != nil {
		return emptyJWT, errors.InternalErr(err, err.Error())
	}

	return token, nil
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(l logger.Logger, pqdriver *pgxpool.Pool) (*Users, error) {
	ustore, err := newStore(pqdriver)
	if err != nil {
		return nil, err
	}

	return &Users{
		logHandler: l,
		store:      ustore,
	}, nil
}
