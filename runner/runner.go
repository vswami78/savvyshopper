package runner

import (
	"context"
	"fmt"
	"io"
	"os"

	"savvyshopper/domain"
	"savvyshopper/internal/config"
	"savvyshopper/internal/price"
	"savvyshopper/internal/render"
)

// Run executes the CLI logic.
// If searchersOpt is provided, it uses those searchers instead of the default ones.
func Run(ctx context.Context, args []string, w io.Writer, searchersOpt ...map[domain.Retailer]price.Searcher) error {
	var query string
	if len(args) > 0 {
		query = args[0]
	} else {
		fmt.Fprint(w, "Enter product: ")
		if _, err := fmt.Fscanln(os.Stdin, &query); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
	}

	_, err := config.APIKey()
	if err != nil {
		fmt.Fprintf(w, "\033[31mError: %v\033[0m\n", err)
		return err
	}

	var offers []domain.Offer
	if len(searchersOpt) > 0 && searchersOpt[0] != nil {
		offers, err = price.SearchPrices(ctx, query, searchersOpt[0])
	} else {
		offers, err = price.SearchPrices(ctx, query)
	}
	if err != nil {
		switch err {
		case domain.ErrNoResults:
			fmt.Fprintf(w, "\033[33mNo results found.\033[0m\n")
		case domain.ErrNetwork:
			fmt.Fprintf(w, "\033[31mNetwork error: %v\033[0m\n", err)
		default:
			fmt.Fprintf(w, "\033[31mError: %v\033[0m\n", err)
		}
		return err
	}

	return render.Table(w, offers)
}
