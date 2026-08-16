package articles

type Service struct {
	articles []Article
}

func NewService(articles []Article) *Service {
	copyOfArticles := append([]Article(nil), articles...)
	return &Service{articles: copyOfArticles}
}

func (s *Service) List() []Article {
	return append([]Article{}, s.articles...)
}

func (s *Service) Search(query string) []Article {
	needle := normalize(query)
	if needle == "" {
		return s.List()
	}
	// Start with a non-nil empty slice so that a no-match search serializes
	// to [] (not null). The list endpoint is expected to return an array
	// even when nothing matches, which keeps the page in a consistent state.
	matches := []Article{}
	for _, article := range s.articles {
		searchable := normalize(article.Title + " " + article.Summary + " " + article.Category + " " + article.Audience)
		if contains(searchable, needle) {
			matches = append(matches, article)
		}
	}
	return matches
}

func (s *Service) Find(id string) (Article, bool) {
	for _, article := range s.articles {
		if article.ID == id {
			return article, true
		}
	}
	return Article{}, false
}

func contains(value, needle string) bool {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
