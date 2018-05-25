package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	//"github.com/nicksnyder/go-i18n/i18n"
	"github.com/autonomousdotai/handshake-exchange/integration/firebase_service"
	"github.com/autonomousdotai/handshake-exchange/service/cache"
	"github.com/autonomousdotai/handshake-exchange/url"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	// Load configuration env
	err := godotenv.Load(".env", "./credentials/.env")
	if err != nil {
		log.Fatal("OrgError loading .env file")
	}
	raven.SetEnvironment(os.Getenv("ENVIRONMENT"))
	raven.SetDSN(os.Getenv("RAVEN_DSN"))
	// End
}

func main() {
	log.Print("Start Crypto Exchange Service")

	// Logger
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/crypto_exchange.log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	})
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// end Logger
	/* Logger
	logFile, err := os.OpenFile("logs/_payment_service.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	 end Logger*/

	// Load configuration
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	sessionPrefix := os.Getenv("SESSION_PREFIX")
	cache.InitializeRedisClient(redisHost, redisPassword)
	// End

	// Load translation
	//i18n.MustLoadTranslationFile("./translations/en-US.flat.yaml")
	//i18n.MustLoadTranslationFile("./translations/zh-HK.flat.yaml")
	// End

	// Setting router
	router := gin.New()
	// Define session
	store, _ := sessions.NewRedisStore(10, "tcp", redisHost, redisPassword, []byte(""))
	router.Use(sessions.Sessions(sessionPrefix, store))

	router.Use(RouterMiddleware())
	router.Use(sentry.Recovery(raven.DefaultClient, false))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Uid",
			"Custom-Language",
			"Custom-Currency"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Router Index
	index := router.Group("/")
	{
		index.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"status": 1, "message": "Crypto Exchange Service works"})
		})
	}

	userUrl := url.UserUrl{}
	userUrl.Create(router)
	infoUrl := url.InfoUrl{}
	infoUrl.Create(router)
	offerUrl := url.OfferUrl{}
	offerUrl.Create(router)
	miscUrl := url.MiscUrl{}
	miscUrl.Create(router)
	cronJobUrl := url.CronJobUrl{}
	cronJobUrl.Create(router)
	creditCardUrl := url.CreditCardUrl{}
	creditCardUrl.Create(router)

	log.Printf(":%s", os.Getenv("SERVICE_PORT"))
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
}

func RouterMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestMethod := context.Request.Method
		requestURL := context.Request.URL.String()

		needToLog := false
		var body interface{}
		if requestMethod == "POST" || requestMethod == "PUT" || requestMethod == "PATCH" || requestMethod == "DELETE " {
			if !strings.Contains(requestURL, "/cron-job/") {
				if requestMethod == "POST" || requestMethod == "PUT" || requestMethod == "PATCH" {
					buf, _ := ioutil.ReadAll(context.Request.Body)
					rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
					rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

					body = readBody(rdr1)
					context.Request.Body = rdr2
				}
				needToLog = true
			}
		}

		requestRemoteAddress := context.Request.RemoteAddr
		userId := context.GetHeader("Custom-UserId")

		dbClient := firebase_service.FirestoreClient
		docId := time.Now().UTC().Format("2006-01-02 15:04:05.000000000")
		docRef := dbClient.Collection("logs").Doc(docId)
		// docId := docRef.ID

		if needToLog {
			log.Println(fmt.Sprintf("%s - %s - %s - %s %s - %s", docId, userId, requestRemoteAddress, requestMethod, requestURL, body))
		}

		context.Next()

		responseStatus := context.Writer.Status()
		responseData, _ := context.Get("ResponseData")
		if needToLog {
			log.Println(fmt.Sprintf("%s - %s - %s - %s", docId, userId, responseStatus, responseData))
			docRef.Set(context, map[string]interface{}{
				"uid":                    userId,
				"request_method":         requestMethod,
				"request_url":            requestURL,
				"request_remote_address": requestRemoteAddress,
				"request_data":           body,
				"response_status":        responseStatus,
				"response_data":          responseData,
				"create_at":              firestore.ServerTimestamp,
			})
		}
	}
}

func readBody(reader io.Reader) interface{} {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	var obj interface{}
	json.Unmarshal([]byte(s), &obj)
	return obj
}
