package api

import (
	"encoding/json"
	"strings"
	"sync"
	"unicode"
)

type sections []SectionData

func newSection(title, duration string, lessonsIndexes []int) SectionData {
	return SectionData{
		Title:        title,
		Duration:     duration,
		LessonsIndex: lessonsIndexes,
	}
}

type SectionData struct {
	Title        string
	Duration     string
	LessonsIndex []int
}

// Initialize a map that is used by createSectionDir,
// to "sanitize" the sectionSlug before creating it.
var charsMap map[rune]string = make(map[rune]string)

func init() {
	channel := make(chan uint16)
	var wg sync.WaitGroup

	unicodeSlice := [][]unicode.Range16{
		unicode.S.R16,
		unicode.P.R16,
	}

	wg.Add(len(unicodeSlice))
	for _, un := range unicodeSlice {
		go func(un []unicode.Range16) {
			defer wg.Done()
			asciiUnicode16(un, channel)
		}(un)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for char := range channel {
		charMap := ""
		switch rune(char) {
		case '&':
			charMap = "and"
		case '-':
			charMap = "-"
		case '(', ')', '[', ']', '/', '\\', '|':
			charMap = "_"
		}

		charsMap[rune(char)] = charMap
	}

	charsMap[' '] = "-"
}

func asciiUnicode16(ranges []unicode.Range16, channel chan uint16) {
	for _, r := range ranges {
		if r.Lo > unicode.MaxASCII {
			break
		}

		for char := r.Lo; char <= r.Hi && char <= unicode.MaxASCII; char += r.Stride {
			channel <- char
		}
	}
}

func (section SectionData) SlugifiedSectionTitle() string {
	sectionTitle := strings.Trim(section.Title, " ")

	var slugifiedTitle strings.Builder

	for _, char := range sectionTitle {
		if mappedValue, ok := charsMap[char]; ok {
			slugifiedTitle.WriteString(mappedValue)
			continue
		}

		slugifiedTitle.WriteRune(char)
	}

	return strings.ToLower(slugifiedTitle.String())
}

type rawSectionsJSON struct {
	RawJSON []json.RawMessage `json:"lessonElements"`
}

func (r *rawSectionsJSON) toLessonElements() sections {
	var secs sections
	var currentSection *SectionData

	for _, jsonElement := range r.RawJSON {
		var jsonObject map[string]any

		if err := json.Unmarshal(jsonElement, &jsonObject); err == nil {
			section := newSection(jsonObject["title"].(string), jsonObject["duration"].(string), []int{})

			secs = append(secs, section)
			currentSection = &secs[len(secs)-1]
		} else {
			var num int

			if err := json.Unmarshal(jsonElement, &num); err == nil && currentSection != nil {
				currentSection.LessonsIndex = append(currentSection.LessonsIndex, num)
			}
		}
	}

	return secs
}
