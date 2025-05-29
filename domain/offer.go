package domain

// Retailer represents a supported retailer.
type Retailer string

const (
    Amazon  Retailer = "Amazon"
    Walmart Retailer = "Walmart"
)

type Offer struct {
    Title    string
    Price    float64
    URL      string
    Retailer Retailer
}
