package api

//import "github.com/dgrijalva/jwt-go"

type UserInformation struct {
	id       string
	username string
	password string
}

type LoginService interface {
	LoginUser(id string, username string, password string) bool
}

func (info *UserInformation) LoginUser(id string, username string, password string) bool {
	return info.id == id && info.username == username && info.password == password
}

func StaticLoginService() LoginService {
	LoginInfo := UserInformation{}
}
