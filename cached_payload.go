package fiber

// NewCachedPayload is a creator factory for a CachedPayload structure
func NewCachedPayload(data []byte) *CachedPayload {
	return &CachedPayload{data: data}
}

// CachedPayload caches []byte contents
type CachedPayload struct {
	data []byte
}

// Payload returns the cached []byte contents
func (b *CachedPayload) Payload() interface{} {
	return b.data
}
