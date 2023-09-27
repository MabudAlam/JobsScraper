// package handlers

// import (
// 	"context"
// 	"net/http"
// "goscraper/models"
// 	"cloud.google.com/go/firestore"
// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/api/iterator"
// )

// func GetJobsFromFirestore(c *gin.Context, firestoreClient *firestore.Client) {
// 	collectionRef := firestoreClient.Collection("Jobs")

// 	iter := collectionRef.Documents(context.Background())

// 	var jobDocuments []models.Job

// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Create a Job model and decode the document data into it.
// 		var job models.Job
// 		if err := doc.DataTo(&job); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		jobDocuments = append(jobDocuments, job)
// 	}

// 	c.JSON(http.StatusOK, jobDocuments)
// }
