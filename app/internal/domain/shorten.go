package domain

type Shorten struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	LongURL   string   `json:"long_url"`
	ShortURL  string   `json:"short_url"`
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

type Shortens []Shorten
