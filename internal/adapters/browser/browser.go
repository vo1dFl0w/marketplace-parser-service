package browser

import "github.com/vo1dFl0w/marketplace-parser-service/internal/repository"

type Browser interface {
	Chromium() repository.BrowserRepository
}