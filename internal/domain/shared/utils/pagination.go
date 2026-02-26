package utils

// Paginator handles pagination logic for API calls
type Paginator struct {
	Limit  int
	Offset int
}

// NewPaginator creates a new paginator with the specified limit
func NewPaginator(limit int) *Paginator {
	return &Paginator{
		Limit:  limit,
		Offset: 0,
	}
}

// Next advances to the next page
func (p *Paginator) Next() {
	p.Offset += p.Limit
}

// HasMore checks if there are more items based on the current count
func (p *Paginator) HasMore(currentCount int) bool {
	return currentCount >= p.Limit
}

// Reset resets the paginator to the first page
func (p *Paginator) Reset() {
	p.Offset = 0
}
