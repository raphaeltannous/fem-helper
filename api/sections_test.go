package api

import "testing"

func TestSectionData_SlugifiedSectionTitle(t *testing.T) {
	sectionsTests := []struct {
		sectionSlug string
		want        string
	}{
		{"Portfolio Footer & Projects", "portfolio-footer-and-projects"},
		{"Adding a Light/Dark Theme Switcher", "adding-a-light_dark-theme-switcher"},
		{"Adding a Light|Dark Theme Switcher", "adding-a-light_dark-theme-switcher"},
		{"Adding a Light\\Dark Theme Switcher", "adding-a-light_dark-theme-switcher"},
		{"Wrapping up", "wrapping-up"},
		{"Routing Q&A", "routing-qanda"},
		{"Protecting Client-Side Routes", "protecting-client-side-routes"},
		{" Scaffolding an API Project ", "scaffolding-an-api-project"},
		{"Search, Filter, & Sort", "search-filter-and-sort"},
	}

	for _, c := range sectionsTests {
		testName := c.sectionSlug
		t.Run(testName, func(t *testing.T) {
			answer := SectionData{
				Title: c.sectionSlug,
			}.SlugifiedSectionTitle()

			if answer != c.want {
				t.Errorf("got %s, want %s", answer, c.want)
			}
		})
	}
}
