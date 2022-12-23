package controller

import (
	"advanced-webapp-project/config"
	"advanced-webapp-project/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
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
	logger *utils.Logger
}

func NewOauthController(logger *utils.Logger) *oauthController {
	return &oauthController{logger: logger}
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

	oa.logger.Info(string(data))
	fmt.Fprintf(c.Writer, "UserInfo: %s\n", data)
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
