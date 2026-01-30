package chromium

import "github.com/vo1dFl0w/marketplace-parser-service/internal/repository"

type Browser struct {
	browserRepo  repository.BrowserRepository
}

// Creates a new browser struct.
func NewBrowser(browserRepo repository.BrowserRepository) *Browser {
	return &Browser{browserRepo: browserRepo}
}

func (b *Browser) Chromium() repository.BrowserRepository {
	return b.browserRepo
}
