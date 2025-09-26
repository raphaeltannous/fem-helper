package api

import (
	"errors"
	"testing"
)

func TestCourseData_FetchFromAPI(t *testing.T) {
	fetchTests := []struct {
		courseSlug string
		want       string
	}{
		{"basics-go", courseSlugNotFoundErr},
		{"go-basics", ""},
	}

	for _, c := range fetchTests {
		testName := c.courseSlug
		t.Run(testName, func(t *testing.T) {
			course := CourseData{
				Slug: c.courseSlug,
			}
			_, answerErr := course.fetchFromAPI()
			if answerErr == nil {
				answerErr = errors.New("")
			}

			if answerErr.Error() != c.want {
				t.Errorf("got %s, want %s", answerErr, c.want)
			}
		})
	}
}
