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
	err := godotenv.Load()
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
			"http://localhost:5000",
			"http://localhost:3000",
			"http://127.0.0.1:5000",
			"http://127.0.0.1:3000",
			"http://35.199.171.18",
			"http://35.199.176.142",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Custom-Username",
			"Custom-Token",
			"Custom-Tfa",
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

	log.Printf(":%s", os.Getenv("SERVICE_PORT"))
	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))
}

func RouterMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// begin before request
		// log.Print("RouterMiddleware: Before request")
		// end

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
		username := context.GetHeader("Custom-Username")

		dbClient := firebase_service.FirestoreClient
		docRef := dbClient.Collection("logs").NewDoc()
		docId := docRef.ID

		if needToLog {
			log.Println(fmt.Sprintf("%s - %s - %s - %s %s - %s", docId, username, requestRemoteAddress, requestMethod, requestURL, body))
		}

		context.Next()

		userId, _ := context.Get("UserId")
		responseStatus := context.Writer.Status()
		responseData, _ := context.Get("ResponseData")
		if needToLog {
			log.Println(fmt.Sprintf("%s - %s - %s - %s", docId, userId, responseStatus, responseData))
			docRef.Set(context, map[string]interface{}{
				"uid":                    userId,
				"username":               username,
				"request_method":         requestMethod,
				"request_url":            requestURL,
				"request_remote_address": requestRemoteAddress,
				"request_data":           body,
				"response_status":        responseStatus,
				"response_data":          responseData,
				"create_at":              firestore.ServerTimestamp,
			})
		}

		// after request
		// log.Print("RouterMiddleware: End request")
		// end
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
