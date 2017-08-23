package github

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/github/hub/ui"
	"github.com/github/hub/utils"
	"github.com/gizak/termui"
	go_github "github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v1"
	// "github.com/k0kubun/pp"
)

// AuthService manages GitHub user authentication
type AuthService struct {
	defaultConfigFile string
	host              string
	user              string
	oauthToken        string
	protocol          string
	client            *go_github.Client
}

type yamlHost struct {
	User       string `yaml:"user"`
	OAuthToken string `yaml:"oauth_token"`
	Protocol   string `yaml:"protocol"`
}
type yamlConfig map[string][]yamlHost

// NewAuthService constructs a new AuthService
func NewAuthService(host string) (a *AuthService, err error) {
	a = new(AuthService)
	homeDir, err := homedir.Dir()
	if err != nil {
		return
	}
	a.defaultConfigFile = filepath.Join(homeDir, ".config", "mcc")
	a.host = host
	a.protocol = "https"
	return
}

// InitClient initializes an AuthService instance
// by loading ~/.config/mcc or basic authentication
func (a *AuthService) InitClient() (client *go_github.Client, err error) {

	if err = a.loadConfig(); err != nil {
		return
	}

	if a.oauthToken == "" {
		if err = a.authorizeClient(a.host); err != nil {
			return
		}
	}

	if err = a.login(); err != nil {
		if err = a.authorizeClient(a.host); err != nil {
			return
		}
	}
	client = a.client

	return
}

func (a *AuthService) login() (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.oauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	a.client = go_github.NewClient(tc)

	// make sure a.client(token) is available
	opt := &go_github.RepositoryListByOrgOptions{Type: "public"}
	_, _, err = a.client.Repositories.ListByOrg(ctx, "github", opt)

	return
}

func (a *AuthService) authorizeClient(host string) (err error) {
	termui.Close()
	ui.Println("-----------------------------------------------------------------------")
	ui.Println("GitHub widget is configured in the yaml")
	ui.Println("Please authenticate your GitHub account first (generating access token)")
	ui.Println("-----------------------------------------------------------------------")

	user := a.promptForUser(host)
	pass := a.promptForPassword(host, user)

	var code, token string
	for {
		token, err = a.createToken(user, pass, code)
		if err == nil {
			break
		}
		if _, ok := err.(*go_github.TwoFactorAuthError); ok {
			if code != "" {
				ui.Errorln("warning: invalid two-factor code")
			}
			code = a.promptForOTP()
		} else {
			break
		}
	}
	if err != nil {
		return
	}
	if token == "" {
		ui.Println("-----------------------------------------------------------------------")
		ui.Errorln("error: invalid username or password")
		os.Exit(1)
	} else {
		a.saveConfig()
		ui.Println("-----------------------------------------------------------------------")
		ui.Println("Authentication succeeded! (access token has been stored ~/.config/mcc)")
		ui.Println("Please restart mcc again")
		os.Exit(0)
	}

	return
}

func (a *AuthService) createToken(user, password, twoFactorCode string) (token string, err error) {
	tp := go_github.BasicAuthTransport{
		Username: user,
		Password: password,
		OTP:      twoFactorCode,
	}
	a.client = go_github.NewClient(tp.Client())

	cnt := 1
	var result *go_github.Authorization
	for {
		desc, err := a.authTokenNote(cnt)
		if err != nil {
			break
		}
		input := &go_github.AuthorizationRequest{
			Note:   go_github.String(desc),
			Scopes: []go_github.Scope{go_github.ScopeRepo},
		}
		result, _, err = a.client.Authorizations.Create(context.Background(), input)

		if err == nil {
			token = result.GetToken()
			a.oauthToken = token
			a.user = result.User.GetName()
			break
		}
		if cnt >= 9 {
			break
		} else {
			cnt++
			continue
		}
	}

	return
}

func (a *AuthService) authTokenNote(num int) (string, error) {
	n := os.Getenv("USER")
	if n == "" {
		n = os.Getenv("USERNAME")
	}
	if n == "" {
		whoami := exec.Command("whoami")
		whoamiOut, err := whoami.Output()
		if err != nil {
			return "", err
		}
		n = strings.TrimSpace(string(whoamiOut))
	}
	h, err := os.Hostname()
	if err != nil {
		return "", err
	}

	if num > 1 {
		return fmt.Sprintf("hub for %s@%s %d", n, h, num), nil
	}

	return fmt.Sprintf("mcc for %s@%s", n, h), nil
}

func (a *AuthService) saveConfig() (err error) {
	err = os.MkdirAll(filepath.Dir(a.defaultConfigFile), 0771)
	if err != nil {
		return
	}

	var w *os.File
	w, err = os.OpenFile(a.defaultConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer w.Close()

	yc := make(yamlConfig)
	yc[a.host] = []yamlHost{
		{
			User:       a.user,
			OAuthToken: a.oauthToken,
			Protocol:   a.protocol,
		},
	}

	d, err := yaml.Marshal(yc)
	if err != nil {
		return
	}

	n, err := w.Write(d)
	if err == nil && n < len(d) {
		err = io.ErrShortWrite
	}
	return
}

func (a *AuthService) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (a *AuthService) loadConfig() (err error) {
	if !a.fileExists(a.defaultConfigFile) {
		return
	}
	r, err := os.Open(a.defaultConfigFile)
	if err != nil {
		return
	}
	defer r.Close()

	d, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	yc := make(yamlConfig)
	err = yaml.Unmarshal(d, &yc)

	if err != nil {
		return
	}

	a.user = yc[a.host][0].User
	a.oauthToken = yc[a.host][0].OAuthToken
	a.protocol = yc[a.host][0].Protocol

	return
}

func (a *AuthService) promptForUser(host string) (user string) {
	user = os.Getenv("GITHUB_USER")
	if user != "" {
		return
	}

	ui.Printf("%s username: ", host)
	user = a.scanLine()

	return
}

func (a *AuthService) promptForPassword(host, user string) (pass string) {
	pass = os.Getenv("GITHUB_PASSWORD")
	if pass != "" {
		return
	}

	ui.Printf("%s password for %s (never stored): ", host, user)
	if ui.IsTerminal(os.Stdin) {
		if password, err := getPassword(); err == nil {
			pass = password
		}
	} else {
		pass = a.scanLine()
	}

	return
}

func (a *AuthService) promptForOTP() string {
	fmt.Print("two-factor authentication code: ")
	return a.scanLine()
}

func (a *AuthService) scanLine() string {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	utils.Check(scanner.Err())
	return line
}

func getPassword() (string, error) {
	stdin := int(syscall.Stdin)
	initialTermState, err := terminal.GetState(stdin)
	if err != nil {
		return "", err
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		terminal.Restore(stdin, initialTermState)
		switch sig := s.(type) {
		case syscall.Signal:
			if int(sig) == 2 {
				fmt.Println("^C")
			}
		}
		os.Exit(1)
	}()

	passBytes, err := terminal.ReadPassword(stdin)
	if err != nil {
		return "", err
	}

	signal.Stop(c)
	fmt.Print("\n")
	return string(passBytes), nil
}
