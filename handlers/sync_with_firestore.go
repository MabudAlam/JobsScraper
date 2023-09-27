package handlers

import (
	"context"
	"fmt"
	"goscraper/models"
	"goscraper/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// Defined a list of scraper functions
var jobScrapers = []struct {
	name    string
	scraper func() ([]models.Job, error)
}{
	{
		name:    "Amazon",
		scraper: services.AmazonScrapper,
	},
	{
		name:    "Atlassian",
		scraper: services.AtlassianScrapper,
	},
	{
		name:    "Coursera",
		scraper: services.CourseraScraper,
	},
	{
		name:    "CRED",
		scraper: services.CredScraper,
	}, {
		name:    "FreshWorks",
		scraper: services.FreshWorksScraper,
	},
	{
		name:    "Gojek",
		scraper: services.GojekScraper,
	},
	{
		name:    "MPL",
		scraper: services.MplScrapper,
	},
}

var firestoreClient *firestore.Client

func SyncWithFireBase(c *gin.Context, firestoreClient *firestore.Client) {
	// Create an empty slice to hold all the job postings.
	var allPostings []models.Job

	// Loop through the list of scrapers and call each one.
	for _, scraper := range jobScrapers {
		postings, err := scraper.scraper()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		allPostings = append(allPostings, postings...)
	}

	fmt.Println("Scraping completed. Total job postings:", len(allPostings))

	// Check if the "Jobs" collection exists in Firestore.
	collectionRef := firestoreClient.Collection("Jobs")
	iter := collectionRef.Limit(1).Documents(context.Background())
	_, err := iter.Next()
	if err == iterator.Done {
		// Collection is empty, which means it doesn't exist.
		fmt.Println("Firestore collection 'Jobs' does not exist. Creating...")
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		// Collection exists and contains at least one document.
		// Delete all documents from the collection.
		fmt.Println("Firestore collection 'Jobs' exists. Deleting existing documents...")
		deleteBatch := firestoreClient.Batch()
		iter = collectionRef.Documents(context.Background())
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			deleteBatch.Delete(doc.Ref)
		}
		_, err := deleteBatch.Commit(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Existing documents deleted.")
	}

	// Use a Firestore batch for more efficient batch writes.
	batch := firestoreClient.Batch()
	defer batch.Commit(context.Background()) // Commit the batch writes at the end.

	// Add the allPostings as new documents inside the "Jobs" collection using Batch.
	fmt.Println("Adding new job postings to Firestore...")
	for _, posting := range allPostings {
		ref := collectionRef.NewDoc()
		posting.ID = ref.ID
		batch.Set(ref, posting) // Use Set to add the document to the batch.
	}

	fmt.Println("Batch write completed.")
	c.JSON(http.StatusOK, allPostings)
}
