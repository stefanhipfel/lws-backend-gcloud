package models

type Image struct {
	ID     int64   `json:"id`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Ratio  float64 `json:"ratio"`
	Url    string  `json:"url"`
}
