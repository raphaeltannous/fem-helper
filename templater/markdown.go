package templater

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/raphaeltannous/fem-helper/api"
	"github.com/raphaeltannous/fem-helper/outputdir"
)

//go:embed templates
var templatesFolder embed.FS

type MarkdownTemplater struct {
	course struct {
		api.CourseData
		Tags []string
	}
	outputDirectory outputdir.OutputDirectory

	courseTemplate *template.Template
	lessonTemplate *template.Template
}

var markdownTemplateFunctions = template.FuncMap{
	"formattags":        formatTagsToMarkdown,
	"formatcoursedata":  formatCourseDataToMarkdown,
	"formatannotations": formatAnnotationsToMarkdown,
}

func NewMarkdownTemplater(course api.CourseData, outputDirector outputdir.OutputDirectory, tags []string, courseTemplate, lessonTemplate string) MarkdownTemplater {
	markdownTemp := MarkdownTemplater{
		course: struct {
			api.CourseData
			Tags []string
		}{
			course, tags,
		},
		outputDirectory: outputDirector,
	}

	markdownTemp.courseTemplate = template.Must(
		template.New(
			filepath.Base(courseTemplate),
		).Funcs(
			markdownTemplateFunctions,
		).ParseFS(
			templatesFolder,
			courseTemplate,
		),
	)

	markdownTemp.lessonTemplate = template.Must(
		template.New(
			filepath.Base(lessonTemplate),
		).Funcs(
			markdownTemplateFunctions,
		).ParseFS(
			templatesFolder,
			lessonTemplate,
		),
	)

	return markdownTemp
}

func (markdown MarkdownTemplater) GenerateCourseMarkdown() error {
	if err := markdown.GenerateCourseFromTemplate(); err != nil {
		return err
	}

	for x, section := range markdown.course.Sections {
		sectionDir, err := markdown.outputDirectory.Create(
			fmt.Sprintf("%d-%s", x, section.SlugifiedSectionTitle()),
		)
		if err != nil {
			return err
		}

		for _, lessonIndex := range section.LessonsIndex {
			lessonHash := markdown.course.LessonsHash[lessonIndex]
			lesson := markdown.course.Lessons[lessonHash]

			if err := markdown.GenerateLessonFromTemplate(sectionDir, lesson); err != nil {
				return err
			}
		}
	}

	return nil
}

func (markdown MarkdownTemplater) GenerateCourseFromTemplate() error {
	outputFile := filepath.Join(
		markdown.outputDirectory.String(),
		fmt.Sprintf("%s.md", markdown.course.Slug),
	)

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return markdown.courseTemplate.Execute(file, markdown.course)
}

func (markdown MarkdownTemplater) GenerateLessonFromTemplate(outputDirectory outputdir.OutputDirectory, lesson api.LessonData) error {
	outputFile := filepath.Join(
		outputDirectory.String(),
		fmt.Sprintf("%02d-%s.md", lesson.Index, lesson.Slug),
	)

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return markdown.lessonTemplate.Execute(file, struct {
		api.LessonData
		Tags       []string
		CourseSlug string
	}{lesson, markdown.course.Tags, markdown.course.Slug})
}

func formatCourseDataToMarkdown(course api.CourseData) string {
	var result strings.Builder

	for x, section := range course.Sections {
		result.WriteString(
			fmt.Sprintf("%d. %s\n", x, section.Title),
		)

		sectionSlug := section.SlugifiedSectionTitle()
		for _, lessonIndex := range section.LessonsIndex {
			lessonHash := course.LessonsHash[lessonIndex]
			lesson := course.Lessons[lessonHash]

			result.WriteString(
				fmt.Sprintf("  - [[%d-%s/%s.md|%d. %s]]\n", x, sectionSlug, lesson.Slug, lessonIndex, lesson.Title),
			)
		}
	}

	return result.String()
}

func formatTagsToMarkdown(tags []string) string {
	var result strings.Builder

	result.WriteByte('\n')
	for _, tag := range tags {
		result.WriteString(
			fmt.Sprintf("  - %s\n", tag),
		)
	}

	return result.String()
}

func formatAnnotationsToMarkdown(annos api.Annotations) string {
	var result strings.Builder

	for _, anno := range annos {
		result.WriteString(
			fmt.Sprintf("\n> [!NOTE]+ %s\n", strings.Join(anno.GetReadableRange(), " -> ")),
		)
		result.WriteString(fmt.Sprintf("> %s\n", anno.Message))
	}

	return result.String()
}
