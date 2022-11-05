package paypal

func findLinkURL(links []*link, rel string) string {
	for _, ln := range links {
		if ln.Rel == rel {
			return ln.URL
		}
	}
	return ""
}

func findApproveURL(links []*link) string {
	return findLinkURL(links, "approve")
}
