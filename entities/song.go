package entities

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
