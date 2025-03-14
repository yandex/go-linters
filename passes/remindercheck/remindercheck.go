package remindercheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"

	"golang.yandex/linters/internal/lintutils"
	"golang.yandex/linters/internal/nogen"
)

func Analyzer() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "remindercheck",
		Doc:  "Checks remind comments are formatted properly",
		Run:  run,
		Requires: []*analysis.Analyzer{
			nogen.Analyzer,
		},
	}

	a.Flags.String("keywords", defaultKeywords, "Comment patterns to check")
	a.Flags.String("format", defaultFormat, "Regular expression for get content groups")

	return a
}

const (
	defaultKeywords = "TODO,FIXME,BUG"
	defaultFormat   = `^([a-zA-Z\-]+\d+)?:?(\s+.*)?$`

	hintTemplate = `'// %s: %s: comment'`
	taskIDHint   = "TASKID-1"
)

func run(pass *analysis.Pass) (interface{}, error) {
	files := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files).List()

	format := pass.Analyzer.Flags.Lookup("format").Value.String()
	re, err := regexp.Compile(format)
	if err != nil {
		return nil, err
	}

	keywords := strings.Split(pass.Analyzer.Flags.Lookup("keywords").Value.String(), ",")
	for i, keyword := range keywords {
		keywords[i] = strings.ToLower(strings.TrimSpace(keyword))
	}

	for _, file := range files {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				err := checkComment(c.Text, keywords, re)
				if err != nil {
					pass.Reportf(c.Pos(), "%s", err.Error())
				}
			}
		}
	}

	return nil, nil
}

func checkComment(text string, keywords []string, re *regexp.Regexp) error {
	const (
		doubleSlashes = "//"
	)

	if text[0:2] != doubleSlashes {
		return nil
	}

	if len(text) < 3 {
		return nil
	}

	text = text[3:]
	from, to := findKeyword(text, keywords)
	if to == 0 || from != 0 {
		return nil
	}

	keyword := text[from:to]
	if !isAllUpperCase(keyword) {
		hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskIDHint)

		return fmt.Errorf("keyword '%s' must be upper case. Required template: %s", keyword, hint)
	}

	if idx := strings.Index(text, doubleSlashes); idx != -1 {
		text = text[:idx]
	}

	shift := len(keyword) + 2
	if len(text) <= shift {
		return makeReportWithRightParts(keyword)
	}

	match := re.FindStringSubmatch(text[shift:])
	const countStringsMatchParts = 3
	if len(match) < countStringsMatchParts {
		return makeReportWithRightParts(keyword)
	}

	taskID, summary := match[1], match[2]
	if taskID == "" {
		hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskIDHint)

		return fmt.Errorf("%s must include task id. Required template: %s", keyword, hint)
	}

	taskArr := strings.Split(taskID, "-")
	const taskArrCount = 2
	if len(taskArr) < taskArrCount {
		hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskIDHint)

		return fmt.Errorf("%s must use valid task id: %s. Required template: %s", keyword, taskID, hint)
	}

	id, err := strconv.Atoi(taskArr[1])
	if err != nil || id < 1 {
		hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskIDHint)

		return fmt.Errorf("%s must use task id number greater zero: %s. Required template: %s", keyword, taskID, hint)
	}

	if idx := strings.Index(summary, doubleSlashes); idx != -1 {
		summary = summary[:id]
	}

	if strings.TrimSpace(summary) == "" {
		hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskID)

		return fmt.Errorf("%s must describe what needs to remind. Required template: %s", keyword, hint)
	}

	return nil
}

func makeReportWithRightParts(keyword string) error {
	hint := fmt.Sprintf(hintTemplate, strings.ToUpper(keyword), taskIDHint)

	return fmt.Errorf("%s must be contains right parts. Required template: %s", keyword, hint)
}

func findKeyword(str string, words []string) (from, to int) {
	str = strings.ToLower(str)

	for _, w := range words {
		if i := strings.Index(str, w); i >= 0 {
			return i, i + len(w)
		}
	}

	return 0, 0
}

func isAllUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}

	return true
}
