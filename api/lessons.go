package api

type lessons map[string]lessonData

type lessonData struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"Title"`
	Description string      `json:"description"`
	Index       int         `json:"index"`
	Timestamp   string      `json:"timestamp"`
	Annotations annotations `json:"annotations"`
}
