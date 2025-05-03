package enums

// HTTPContentType enum
type HTTPContentType string

func (h HTTPContentType) String() string {
	return string(h)
}

const (
	ContentTypeText HTTPContentType = "text/plain"
	ContentTypeJSON HTTPContentType = "application/json"
)
