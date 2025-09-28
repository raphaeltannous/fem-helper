package api

type lessons map[string]LessonData

type LessonData struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"Title"`
	Description string      `json:"description"`
	Index       int         `json:"index"`
	Timestamp   string      `json:"timestamp"`
	Annotations Annotations `json:"annotations"`
}
