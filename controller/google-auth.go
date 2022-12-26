package controller

import (
	"advanced-webapp-project/config"
	"advanced-webapp-project/model"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var appConfig, _ = config.LoadConfig(".")
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  appConfig.GoogleOauthRedirectUrl,
	ClientID:     appConfig.GoogleOauthClientId,
	ClientSecret: appConfig.GoogleOauthClientSecret,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type IOauthController interface {
	GoogleOauthLogin(c *gin.Context)
	GoogleOauthCallback(c *gin.Context)
}

type oauthController struct {
	logger      *utils.Logger
	jwtService  service.IJWTService
	authService service.IAuthService
}

func NewOauthController(logger *utils.Logger, jwtSvc service.IJWTService, authSvc service.IAuthService) *oauthController {
	return &oauthController{logger: logger, jwtService: jwtSvc, authService: authSvc}
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (oa *oauthController) GoogleOauthLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c.Writer)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
}

func (oa *oauthController) GoogleOauthCallback(c *gin.Context) {
	oauthState, _ := c.Request.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthState.Value {
		oa.logger.Warn("invalid oauth google state")
		http.Redirect(c.Writer, c.Request, appConfig.FrontendLoginUrl, http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		oa.logger.Error(err.Error())
		http.Redirect(c.Writer, c.Request, appConfig.FrontendLoginUrl, http.StatusTemporaryRedirect)
		return
	}

	type googleRes struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	var googleInfo googleRes
	_ = json.Unmarshal(data, &googleInfo)

	var newUser model.User

	userData, _ := oa.authService.GetUserByEmail(googleInfo.Email)
	if userData == nil {
		newUser.FullName = googleInfo.Name
		newUser.Email = googleInfo.Email
		newUser.ProfileImg = googleInfo.Picture
		newUser.IsVerified = true
		userId, _ := oa.authService.CreateUser(&newUser)
		newUser.Id = uint(userId)
	} else {
		newUser = *userData
	}
	oa.logger.Info(newUser)

	generatedToken := oa.jwtService.GenerateToken(strconv.Itoa(int(newUser.Id)), newUser.Email)
	c.JSON(http.StatusOK, map[string]any{
		"user":  userData,
		"token": generatedToken,
	})
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
