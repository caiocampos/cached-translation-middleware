package model

// TranslationRequest represents the incoming translation request payload.
type TranslationRequest struct {
	Q      string `json:"q"      binding:"required"`
	Source string `json:"source" binding:"required"`
	Target string `json:"target" binding:"required"`
}

// TranslationResponse represents the response returned by the middleware and the upstream API.
type TranslationResponse struct {
	TranslatedText string `json:"translatedText"`
}

// CacheSource indicates whether the response was served from cache or the upstream API.
type CacheSource string

const (
	CacheSourceHit  CacheSource = "cache"
	CacheSourceMiss CacheSource = "upstream"
)
