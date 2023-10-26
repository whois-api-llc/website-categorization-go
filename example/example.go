package example

import (
	"context"
	"errors"
	websitecategorization "github.com/whois-api-llc/website-categorization-go"
	"log"
)

func GetData(apikey string) {
	client := websitecategorization.NewBasicClient(apikey)

	// Get parsed Website Categorization API response by a domain name as a model instance.
	wCategorizationResp, resp, err := client.Get(context.Background(),
		"whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only.
		websitecategorization.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *websitecategorization.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print all returned categories for the domain name.
	for _, obj := range wCategorizationResp.Categories {
		log.Printf("ID: %d, Name: %s, Confidence: %f ",
			obj.ID, obj.Name, obj.Confidence)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := websitecategorization.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.GetRaw(context.Background(),
		"whoisxmlapi.com",
		// this option causes those only categories having a relevance greater than 0.8 to be returned.
		websitecategorization.OptionMinConfidence(0.8))

	if err != nil {
		// Handle error message returned by server.
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}

func GetCategories(apikey string) {
	client := websitecategorization.NewBasicClient(apikey)

	// Get all possible categories as an array.
	wCategorizationResp, _, err := client.WCategorizationService.GetAllCategories(context.Background(),
		// this option causes the categories to be ordered alphabetically by the name field.
		websitecategorization.OptionOrder("ABC"))

	if err != nil {
		// Handle error message returned by server.
		log.Fatal(err)
	}

	// Then print IDs and names for all possible categories.
	for _, obj := range wCategorizationResp {
		log.Println(obj.ID, obj.Name)
	}
}

func GetCategoriesRaw(apikey string) {
	client := websitecategorization.NewBasicClient(apikey)

	// Get all possible categories as a raw API response.
	resp, err := client.WCategorizationService.GetAllCategoriesRaw(context.Background())

	if err != nil {
		// Handle error message returned by server.
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
