package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/raphaeltannous/fem-helper/api"
	"github.com/raphaeltannous/fem-helper/outputdir"
	"github.com/raphaeltannous/fem-helper/templater"
)

var (
	courseSlug string
	outputDir  string
)

func init() {
	const (
		outputDirHelpString  = "Output directory of the course."
		courseSlugHelpString = "Slug of the course."
	)

	flag.StringVar(&courseSlug, "course-slug", "", courseSlugHelpString+" (required)")
	flag.StringVar(&courseSlug, "c", "", courseSlugHelpString+" (shorthand/required)")

	flag.StringVar(&outputDir, "output-dir", "", outputDirHelpString)
	flag.StringVar(&outputDir, "o", "", outputDirHelpString+" (shorthand)")

	// todo:
	// implement flag
	// --clean will clean cache
}

type customUserTemplates []string

func (cUT *customUserTemplates) String() string {
	return fmt.Sprint(*cUT)
}

func (cUT *customUserTemplates) Set(path string) error {
	value := filepath.Base(path)
	if slices.Contains([]string{"course.tmpl", "lesson.tmpl"}, value) {
		*cUT = append(*cUT, path)
		return nil
	}

	return errors.New("allowed template filenames: course.tmpl and lesson.tmpl")
}

func (cUT *customUserTemplates) checkTemplates() {
	defaultTemplatePath := "templates/obsidian/"

	seen := make(map[string]bool)
	for _, templatePath := range *cUT {
		seen[filepath.Base(templatePath)] = true
	}

	for _, templateName := range []string{"course.tmpl", "lesson.tmpl"} {
		if !seen[templateName] {
			*cUT = append(*cUT, filepath.Join(defaultTemplatePath, templateName))
		}
	}
}

func (cUT *customUserTemplates) getTemplateByName(name string) string {
	for _, templatePath := range *cUT {
		if name == filepath.Base(templatePath) {
			return templatePath
		}
	}

	return ""
}

var customTemplates customUserTemplates

func init() {
	flag.Var(&customTemplates, "custom-template", "Custom templates for course and lesson. (Allowed filenames: course.tmpl and lesson.tmpl)")
}

type tags []string

func (t *tags) String() string {
	return fmt.Sprint(*t)
}

func (t *tags) Set(value string) error {
	if len(*t) > 0 {
		return errors.New("tags flag already set.")
	}

	if len(value) == 0 {
		return errors.New("tags flag cannot be an empty string.")
	}

	for tag := range strings.SplitSeq(value, ",") {
		*t = append(*t, tag)
	}

	return nil
}

var tagsFlag tags

func init() {
	var (
		tagsFlagHelp = "comma-seperated list of tags."
	)

	flag.Var(&tagsFlag, "tags", tagsFlagHelp)
	flag.Var(&tagsFlag, "t", tagsFlagHelp+" (shorthand)")
}

func main() {
	flag.Parse()

	requiredFlags([][2]string{
		{"course-slug", "c"},
		{"output-dir", "o"},
	})
	customTemplates.checkTemplates()

	course, err := api.NewCourse(courseSlug)
	if err != nil {
		log.Fatal(err)
	}

	outputDirectory, err := outputdir.NewOutputDirectory(outputDir)
	if err != nil {
		log.Fatal(err)
	}

	markdown := templater.NewMarkdownTemplater(
		course,
		outputDirectory,
		tagsFlag,
		customTemplates.getTemplateByName("course.tmpl"),
		customTemplates.getTemplateByName("lesson.tmpl"),
	)

	if err := markdown.GenerateCourseMarkdown(); err != nil {
		log.Fatal(err)
	}
}

func requiredFlags(requiredFlags [][2]string) {
	givenFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		givenFlags[f.Name] = true
	})

	for _, duo := range requiredFlags {
		if givenFlags[duo[0]] || givenFlags[duo[1]] {
			continue
		}

		fmt.Fprintf(os.Stderr, "%s flag is required.\n", duo[0])
		os.Exit(2)
	}
}
