package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-json"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniversalDomain         string `json:"universal_domain"`
}

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	projectID := viper.GetString("FIREBASE_PROJECT_ID")
	port := viper.GetString("PORT")
	universalDomain := viper.GetString("UNIVERSAL_DOMAIN")

	privateKeyID := viper.GetString("PRIVATE_KEY_ID")
	privateKey := viper.GetString("PRIVATE_KEY")
	stype := viper.GetString("TYPE")

	clientEmail := viper.GetString("CLIENT_EMAIL")
	clientID := viper.GetString("CLIENT_ID")
	authURI := viper.GetString("AUTH_URI")
	tokenURI := viper.GetString("TOKEN_URI")
	authProviderX509CertURL := viper.GetString("AUTH_PROVIDER_X509_CERT_URL")
	clientX509CertURL := viper.GetString("CLIENT_X509_CERT_URL")
	privateKey = strings.Replace(privateKey, "\\n", "\n", -1)

	fireBaseConfig := FirebaseConfig{
		Type:                    stype,
		ProjectID:               projectID,
		PrivateKeyID:            privateKeyID,
		PrivateKey:              privateKey,
		ClientEmail:             clientEmail,
		ClientID:                clientID,
		AuthURI:                 authURI,
		TokenURI:                tokenURI,
		AuthProviderX509CertURL: authProviderX509CertURL,
		ClientX509CertURL:       clientX509CertURL,
		UniversalDomain:         universalDomain,
	}

	// Convert the FirebaseConfig to a JSON string.
	fireBaseConfigJSON, err := json.MarshalIndent(fireBaseConfig, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling FirebaseConfig to JSON: %v", err)
	}

	filePath := "./firebase-config.json"
	err = os.WriteFile(filePath, fireBaseConfigJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	fmt.Println("JSON file written:", filePath)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	fmt.Println("JSON file content:")
	fmt.Println(string(fileContent))

	r := gin.Default()

	r.GET("/cred", handlers.GetPostingsHandler)
	r.GET("/atlassian", handlers.AtlassianHandler)
	r.GET("/amazon", handlers.Amazonhandler)
	r.GET("/coursera", handlers.CourseraHandler)
	r.GET("/freshworks", handlers.FreshWorksHandler)
	r.GET("/gojek", handlers.Gojekhandler)
	r.GET("/mpl", handlers.MplHandler)

	r.GET("/syncwithSql", handlers.AllScrapersHandler)
	r.GET("/sync", handlers.SyncAll)

	r.GET("/getallJobs", handlers.GetJobsFromDB)
	r.GET("/getallJobsFromSQL", handlers.GetAllJobsFromSqlite)

	opt := option.WithCredentialsFile(filePath)
	firestoreClient, err := firestore.NewClient(context.Background(), projectID, opt)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
	}

	r.GET("/syncFirestore", func(c *gin.Context) {
		handlers.SyncWithFireBase(c, firestoreClient)
	})

	r.GET("/getJobsFromFirestore", func(c *gin.Context) {
		handlers.GetJobsFromFirestore(c, firestoreClient)
	})

	r.GET("/", func(c *gin.Context) {
		c.File("static/base.html")
	})
	r.Use(CORSMiddleware())

	if err := r.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RootHandler handles the root route (/).
func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"hello": "Mabud"})
}
