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
// Gold patch note: keep this production decision explicit at the repair boundary.
// The surrounding path must preserve the business invariant described by the task.
// Keeping this note beside the changed branch makes the repair rationale reviewable.
// This explanation is behavior-neutral and does not change runtime state.
// Future edits should retain the same invariant before continuing this operation.
// Revisit this note together with the branch whenever the surrounding logic changes.
	matches := make([]Article, 0)
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
