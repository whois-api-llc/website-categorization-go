package websitecategorization

import (
	"fmt"
)

// CategoryItem is part of the category directory.
type CategoryItem struct {
	// ID is the unique category identifier.
	ID int `json:"id"`

	// Name is the readable name of the category.
	Name string `json:"name"`
}

// WCategorizationResponse is a response of Website Categorization API.
type WCategorizationResponse struct {
	// AS Autonomous System.
	AS *AS `json:"as,omitempty"`

	// DomainName is a domain/website name.
	DomainName string `json:"domainName"`

	// Categories is the list of website's categories.
	Categories []Category `json:"categories"`

	// CreatedDate is date of initial creation of the WHOIS record for the domain in ISO8601 format. Omitted if the record is not found.
	CreatedDate *string `json:"createdDate,omitempty"`

	// WebsiteResponded Determines if the website was active during the crawling.
	WebsiteResponded bool `json:"websiteResponded"`
}

// Category is a part of the Website Categorization API v3 response.
type Category struct {
	// Confidence The probability of how the category may be relevant for the website.
	Confidence float64 `json:"confidence"`

	// ID The unique category identifier.
	ID int `json:"id"`

	// Name The readable name of the category.
	Name string `json:"name"`
}

// AS is a part of the Website Categorization API response.
type AS struct {
	// ASN Autonomous System Number.
	ASN int `json:"asn"`

	// Domain Autonomous System Website's URL.
	Domain string `json:"domain"`

	// Name Autonomous System Name.
	Name string `json:"name"`

	// Route Autonomous System Route.
	Route string `json:"route"`

	// Type Autonomous System Type.
	Type string `json:"type"`
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
