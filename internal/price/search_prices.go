package price

import (
	"context"
	"sort"

	"savvyshopper/domain"
)

// SearchPrices queries both Amazon and Walmart, merges, sorts, and enforces invariants.
// If searchers is nil, uses the default real searchers.
func SearchPrices(ctx context.Context, query string, searchersOpt ...map[domain.Retailer]searcher) ([]domain.Offer, error) {
	var searchers map[domain.Retailer]searcher
	if len(searchersOpt) > 0 && searchersOpt[0] != nil {
		searchers = searchersOpt[0]
	} else {
		searchers = map[domain.Retailer]searcher{
			domain.Amazon:  NewAmazonSearcher("https://api.zinc.io/v1/search/amazon"),
			domain.Walmart: NewWalmartSearcher("https://api.zinc.io/v1/search/walmart"),
		}
	}

	var allOffers []domain.Offer
	for retailer, s := range searchers {
		offers, err := s.Search(ctx, query)
		if err != nil {
			// TODO: Wrap/handle network/auth errors as needed
			continue
		}
		// Set retailer field (defensive, in case helpers don't)
		for i := range offers {
			offers[i].Retailer = retailer
		}
		if len(offers) > 3 {
			offers = offers[:3]
		}
		allOffers = append(allOffers, offers...)
	}

	// Enforce invariants
	if len(allOffers) == 0 {
		return nil, domain.ErrNoResults
	}
	if len(allOffers) > 6 {
		allOffers = allOffers[:6]
	}
	// Sort by ascending price
	sort.Slice(allOffers, func(i, j int) bool {
		return allOffers[i].Price < allOffers[j].Price
	})
	// Ensure price >= 0
	for _, offer := range allOffers {
		if offer.Price < 0 {
			return nil, domain.ErrNetwork // or another error
		}
	}

	return allOffers, nil
}
