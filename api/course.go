package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/raphaeltannous/fem-helper/cache"
)

const apiUrl = "https://api.frontendmasters.com/v2/kabuki/courses/"

const (
	courseSlugNotFoundErr = "course slug is not correct."
)

type CourseData struct {
	Slug          string `json:"slug"`
	Title         string `json:"Title"`
	DatePublished string `json:"datePublished"`
	Description   string `json:"description"`
	LessonsHash   []string

	Sections Sections

	Lessons lessons `json:"lessonData"`
}

func NewCourse(slug string) (CourseData, error) {
	course := CourseData{
		Slug: slug,
	}

	err := course.fetchAndPopulateJSON()
	if err != nil {
		return CourseData{}, nil
	}

	course.populateLessonsHash()

	return course, nil
}

// Populates course.LessonsHash by lesson hash by index of each lesson.
func (course *CourseData) populateLessonsHash() {
	course.LessonsHash = make([]string, course.lessonsCount())

	for lessonHash, lesson := range course.Lessons {
		course.LessonsHash[lesson.Index] = lessonHash
	}
}

// Convert the fetched data into our courseInfo struct and populate its fields.
func (course *CourseData) fetchAndPopulateJSON() error {
	requestBody, err := course.fetch()
	if err != nil {
		return err
	}

	err = json.Unmarshal(requestBody, &course)
	if err != nil {
		return err
	}

	// Unmarshalling Sections
	var rawSections rawSectionsJSON
	if err := json.Unmarshal(requestBody, &rawSections); err != nil {
		return err
	}

	course.Sections = rawSections.toLessonElements()

	return nil
}

// Returns the body of the api request, and an error if one occurs.
func (course *CourseData) fetch() ([]byte, error) {
	if data, err := course.loadFromCache(); err == nil {
		return data, nil
	}

	body, err := course.fetchFromAPI()
	if err != nil {
		return nil, err
	}

	if _, err := course.addToCache(body); err != nil {
		return nil, err
	}
	return body, nil
}

// Fetch the course json data from the apiURL.
// If course slug is not valid, courseSlugNotFoundErr will be returned as an error.
func (course *CourseData) fetchFromAPI() ([]byte, error) {
	fetchUrl := apiUrl + course.Slug

	resp, err := http.Get(fetchUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(courseSlugNotFoundErr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Returns the json from cache if available.
func (course *CourseData) loadFromCache() ([]byte, error) {
	cacheDir, err := cache.NewCache()
	if err != nil {
		return nil, err
	}

	filename := course.Slug + ".json"

	return cacheDir.Read(filename)
}

// Add fetched json data to cache. Returns the number of bytes written,
// and an error if one occurs.
func (course *CourseData) addToCache(data []byte) (int, error) {
	cacheDir, err := cache.NewCache()
	if err != nil {
		return 0, err
	}

	filename := course.Slug + ".json"

	return cacheDir.Save(filename, data)
}

// Returns the lesson count of the course.
// Needs to be used after course.fetchJson(), otherwise
// the return value is -1.
func (course *CourseData) lessonsCount() int {
	if len(course.Sections) == 0 {
		return -1
	}

	courseSectionsLen := len(course.Sections)
	lastSection := course.Sections[courseSectionsLen-1]
	return lastSection.LessonsIndex[len(lastSection.LessonsIndex)-1] + 1
}

// Returns the course Duration.
func (course *CourseData) courseDuration() (time.Duration, error) {
	duration := 0 * time.Second

	for _, section := range course.Sections {

		currentDur, err := time.ParseDuration(strings.ReplaceAll(section.Duration, " ", ""))
		if err != nil {
			return time.Duration(0), fmt.Errorf("error while parsing section Duration: %w", err)
		}

		duration += currentDur
	}

	return duration, nil
}

func (course CourseData) String() string {
	var result strings.Builder
	result.WriteString("CourseInfo:\n")

	result.WriteString(fmt.Sprintf("\tTitle: %v\n", course.Title))
	result.WriteString(fmt.Sprintf("\tDescription: %v\n", course.Description))
	result.WriteString(fmt.Sprintf("\tNumber of Lessons: %d\n", course.lessonsCount()))

	duration, err := course.courseDuration()
	if err != nil {
		result.WriteString("\tDuration: Failed to parse Duration.")
	} else {
		result.WriteString(fmt.Sprintf("\tDuration: %v\n", duration))
	}

	return result.String()
}
