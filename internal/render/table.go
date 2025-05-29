package render

import (
	"fmt"
	"io"
	"text/tabwriter"

	"savvyshopper/domain"
)

// Table writes the offers to w in a tabular format.
func Table(w io.Writer, offers []domain.Offer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Title\tPrice\tRetailer\tURL")
	for _, offer := range offers {
		title := offer.Title
		if len(title) > 60 {
			title = title[:60]
		}
		fmt.Fprintf(tw, "%s\t$%.2f\t%s\t%s\n", title, offer.Price, offer.Retailer, offer.URL)
	}
	return tw.Flush()
}
