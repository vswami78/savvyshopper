package price

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"savvyshopper/domain"
)

// SearchPrices queries both Amazon and Walmart concurrently, merges, sorts, and enforces invariants.
// If searchers is nil, uses the default real searchers.
func SearchPrices(ctx context.Context, query string, searchersOpt ...map[domain.Retailer]Searcher) ([]domain.Offer, error) {
	var searchers map[domain.Retailer]Searcher
	if len(searchersOpt) > 0 && searchersOpt[0] != nil {
		searchers = searchersOpt[0]
	} else {
		searchers = map[domain.Retailer]Searcher{
			domain.Amazon:  NewAmazonSearcher("https://api.zinc.io/v1/search/amazon"),
			domain.Walmart: NewWalmartSearcher("https://api.zinc.io/v1/search/walmart"),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	type result struct {
		retailer domain.Retailer
		offers   []domain.Offer
		err      error
	}
	ch := make(chan result, len(searchers))
	var wg sync.WaitGroup

	for retailer, s := range searchers {
		wg.Add(1)
		go func(retailer domain.Retailer, s Searcher) {
			defer wg.Done()
			offers, err := s.Search(ctx, query)
			ch <- result{retailer: retailer, offers: offers, err: err}
		}(retailer, s)
	}

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	var allOffers []domain.Offer
	var firstErr error
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: search timed out", domain.ErrNetwork)
		case res, ok := <-ch:
			if !ok {
				goto DONE
			}
			if res.err != nil {
				if firstErr == nil && (errors.Is(res.err, domain.ErrNetwork) || errors.Is(res.err, domain.ErrAuth)) {
					firstErr = res.err
				}
				continue
			}
			// Set retailer field (defensive, in case helpers don't)
			for i := range res.offers {
				res.offers[i].Retailer = res.retailer
			}
			if len(res.offers) > 3 {
				res.offers = res.offers[:3]
			}
			allOffers = append(allOffers, res.offers...)
		}
	}
DONE:
	if len(allOffers) == 0 {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, domain.ErrNoResults
	}
	if len(allOffers) > 6 {
		allOffers = allOffers[:6]
	}
	sort.Slice(allOffers, func(i, j int) bool {
		return allOffers[i].Price < allOffers[j].Price
	})
	for _, offer := range allOffers {
		if offer.Price < 0 {
			return nil, fmt.Errorf("%w: negative price found", domain.ErrNetwork)
		}
	}
	return allOffers, nil
}
