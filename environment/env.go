package environment

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	// Versioning
	CurrentVersion string

	// Token related
	APISecretKey             string
	AuthTokenValidityDays    string
	RefreshTokenValidityDays string

	// Router & path related
	RouterSecretPath  string
	RouterSecretPath2 string

	// URLs
	BaseURLForScraping   string
	EpisodesURLFormat    string
	TitleSearchURLFormat string
	FCMURL               string

	// CSS selectors
	HomeSelector             string
	StreamSelector           string
	EpisodesSelector         string
	DetailSelector           string
	TitlesSelector           string
	IDAtDetailSelector       string
	ImageURLAtDetailSelector string
	VideoURLSelector         string

	// :)
	SuperSecretKey1 string
	SuperSecretKey2 string

	// Anime detail
	Keyword1 string
	Keyword2 string
	Keyword3 string
	Keyword4 string
	Keyword5 string

	// DB connection string
	RDBMSConnString string
	DBMSConnString  string

	// Token
	GHAuthToken string

	// FCM
	FCMKey string

	// Static file path
	InitSQLPath string
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		// Is in production
		err = getAllRequiredFiles()
		if err != nil {
			log.Fatal(err)
		}

		err = godotenv.Load("/app/.env")
		if err != nil {
			log.Fatal(err)
		}
	}

	CurrentVersion = os.Getenv("CURRENT_VERSION")
	APISecretKey = os.Getenv("API_SECRET_KEY")
	AuthTokenValidityDays = os.Getenv("AUTH_TOKEN_VALIDITY_DAYS")
	RefreshTokenValidityDays = os.Getenv("REFRESH_TOKEN_VALIDITY_DAYS")
	RouterSecretPath = os.Getenv("ROUTER_SECRET_PATH")
	RouterSecretPath2 = os.Getenv("ROUTER_SECRET_PATH2")
	BaseURLForScraping = os.Getenv("BASE_URL_FOR_SCRAPING")
	EpisodesURLFormat = os.Getenv("EPISODES_URL_FORMAT")
	TitleSearchURLFormat = os.Getenv("TITLE_SEARCH_URL_FORMAT")
	FCMURL = os.Getenv("FCM_URL")
	HomeSelector = os.Getenv("HOME_SELECTOR")
	StreamSelector = os.Getenv("STREAM_SELECTOR")
	EpisodesSelector = os.Getenv("EPISODES_SELECTOR")
	DetailSelector = os.Getenv("DETAIL_SELECTOR")
	TitlesSelector = os.Getenv("TITLES_SELECTOR")
	IDAtDetailSelector = os.Getenv("ID_AT_DETAIL_SELECTOR")
	ImageURLAtDetailSelector = os.Getenv("IMAGE_URL_AT_DETAIL_SELECTOR")
	VideoURLSelector = os.Getenv("VIDEO_URL_SELECTOR")
	SuperSecretKey1 = os.Getenv("SUPER_SECRET_KEY1")
	SuperSecretKey2 = os.Getenv("SUPER_SECRET_KEY2")
	Keyword1 = os.Getenv("KEYWORD1")
	Keyword2 = os.Getenv("KEYWORD2")
	Keyword3 = os.Getenv("KEYWORD3")
	Keyword4 = os.Getenv("KEYWORD4")
	Keyword5 = os.Getenv("KEYWORD5")
	RDBMSConnString = os.Getenv("RDBMS_CONN_STRING")
	DBMSConnString = os.Getenv("DBMS_CONN_STRING")
	GHAuthToken = os.Getenv("GH_AUTH_TOKEN")
	FCMKey = os.Getenv("FCM_KEY")
	InitSQLPath = os.Getenv("INIT_SQL_PATH")
}

func getAllRequiredFiles() error {
	fileBytes, err := os.ReadFile("/app/urls.txt")
	if err != nil {
		return err
	}

	sliceFileBytes := strings.Split(string(fileBytes), "\n")
	envURL := sliceFileBytes[0]
	certsURL := sliceFileBytes[1]

	err = downloadFile(envURL, "/app/.env")
	if err != nil {
		return err
	}

	err = downloadFile(certsURL, "/app/prod-ca-2021.crt")
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(downloadURL, dstPath string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f.Write(fileBytes)

	return nil
}
