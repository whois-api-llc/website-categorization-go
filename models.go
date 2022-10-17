package websitecategorization

import (
	"fmt"
)

// CategoryItem is part of the category directory.
type CategoryItem struct {
	// ID is the unique category identifier.
	ID string `json:"id"`

	// Name is the readable name of the category.
	Name string `json:"name"`

	// Parent is the ID of parent category (if present).
	Parent *string `json:"parent"`
}

// Tier is the category object.
type Tier struct {
	// ID is the unique category identifier.
	ID string `json:"id"`

	// Confidence is the probability of how the category may be relevant for the website.
	Confidence float64 `json:"confidence"`

	// Name is the readable name of the category.
	Name string `json:"name"`
}

// Category is a part of the Website Categorization API response.
type Category struct {
	// Tier1 is the top level category object.
	Tier1 *Tier `json:"tier1"`

	// Tier2 is the 2nd level category object (if present).
	Tier2 *Tier `json:"tier2"`
}

// WCategorizationResponse is a response of Website Categorization API.
type WCategorizationResponse struct {
	// Result is the list of website's categories.
	Categories []Category `json:"categories"`

	// DomainName is a domain/website name.
	DomainName string `json:"domainName"`

	// WebsiteResponded Determines if the website was active during the crawling.
	WebsiteResponded bool `json:"websiteResponded"`
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
