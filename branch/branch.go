package branch

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
    "twig/issue"
    "twig/log"
    "twig/network"
)

const (
    branchTypeSeparator = "/"
    issueTypeSeparator  = "_"
    wordSeparator       = "-"
)

type Branch struct {
    Type           Type
    ExcludePhrases []string

    stripRegx           *regexp.Regexp
    firstPassKebabRegx  *regexp.Regexp
    secondPassKebabRegx *regexp.Regexp
    issueRegx           *regexp.Regexp
    excludePhrasesRegx  []*regexp.Regexp
}

func New(branchType Type, excludePhrases []string) *Branch {
    b := new(Branch)

    b.Type = branchType
    b.ExcludePhrases = excludePhrases

    b.stripRegx = regexp.MustCompile("[^a-zA-Z0-9]+")
    b.firstPassKebabRegx = regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
    b.secondPassKebabRegx = regexp.MustCompile("([a-z])([A-Z])")
    b.issueRegx = regexp.MustCompile(`[A-Z]+-\d+_`) // looking for XXXX-0000_
    b.excludePhrasesRegx = prepareExcludeRegx(excludePhrases)

    return b
}

func prepareExcludeRegx(excludes []string) []*regexp.Regexp {
    regexps := make([]*regexp.Regexp, len(excludes))

    for i, v := range excludes {
        re := regexp.MustCompile("(?i)(\\[" + v + "\\]|\\(" + v + "\\))")
        regexps[i] = re
    }

    return regexps
}

func (b *Branch) BuildName(jiraIssue network.JiraIssue) string {
    log.Info().Println("Prepare branch")

    branchType := b.Type.ToString()
    log.Debug().Println(fmt.Sprintf("Issue %s(%s), type %q", jiraIssue.Key, jiraIssue.Fields.Type.Id, branchType))

    var buffer strings.Builder

    summary := b.replacePhrases(*jiraIssue.Fields.Summary)
    summary = b.camelToKebab(summary)
    summary = b.stripRegex(summary)

    if b.Type != NULL {
        buffer.WriteString(branchType)
        buffer.WriteString(branchTypeSeparator)
    }

    buffer.WriteString(jiraIssue.Key)
    buffer.WriteString(issueTypeSeparator)
    buffer.WriteString(summary)

    return buffer.String()
}

func (b *Branch) stripRegex(in string) string {
    phrase := strings.ToLower(in)
    log.Debug().Println(fmt.Sprintf("Before strip %q", phrase))

    phrase = b.stripRegx.ReplaceAllString(phrase, wordSeparator)
    phrase = strings.TrimPrefix(phrase, wordSeparator)
    phrase = strings.TrimSuffix(phrase, wordSeparator)

    log.Debug().Println(fmt.Sprintf("After strip %s", phrase))
    return phrase
}

func (b *Branch) replacePhrases(in string) string {
    log.Debug().Println(fmt.Sprintf("Before replace %q", in))

    phrase := in
    for _, re := range b.excludePhrasesRegx {
        phrase = re.ReplaceAllString(phrase, "")
    }

    phrase = strings.TrimSpace(phrase)
    log.Debug().Println(fmt.Sprintf("After replace %q", phrase))
    return phrase
}

func (b *Branch) camelToKebab(in string) string {
    log.Debug().Println(fmt.Sprintf("Before kebab %q", in))

    kebab := b.firstPassKebabRegx.ReplaceAllString(in, "${1}"+wordSeparator+"${2}")
    kebab = b.secondPassKebabRegx.ReplaceAllString(kebab, "${1}"+wordSeparator+"${2}")

    log.Debug().Println(fmt.Sprintf("After kebab %q", kebab))
    return strings.ToLower(kebab)
}

func (b *Branch) ExtractIssueNameFromBranch(branchName string) (string, error) {
    log.Debug().Println(fmt.Sprintf("Before extract %q", branchName))

    match := b.issueRegx.FindString(branchName)
    match = strings.TrimSpace(match)
    match = strings.TrimSuffix(match, issueTypeSeparator)

    if match == "" {
        return "", errors.New("no issue match")
    }

    log.Debug().Println(fmt.Sprintf("After extract %q", match))
    return match, nil
}

func InputToBranchType(input string) (Type, error) {
    switch input {
    case Build, BuildShort:
        return BUILD, nil
    case Chore, ChoreShort:
        return CHORE, nil
    case Ci:
        return CI, nil
    case Docs, DocsShort:
        return DOCS, nil
    case Feat, FeatShort:
        return FEAT, nil
    case Fix, FixShort:
        return FIX, nil
    case Perf, PerfShort:
        return PERF, nil
    case Refactor, RefactorShort:
        return REFACTOR, nil
    case Revert, RevertShort:
        return REVERT, nil
    case Style, StyleShort:
        return STYLE, nil
    case Test, TestShort:
        return TEST, nil
    default:
        return NULL, fmt.Errorf("unsupported branch type %q", input)
    }
}

func ConvertIssueTypesToMap(issueTypes []network.IssueType) (map[string]Type, error) {
    issueMap := make(map[string]Type)

    local, err := issue.ParseIssueMapping()
    if err != nil {
        return nil, err
    }

    for _, i := range issueTypes {
        id, ok := local[i.Id]
        if !ok {
            log.Debug().Println(fmt.Sprintf("Unsupported issue type %s (%s)", i.Name, i.Id))
            continue
        }

        name, err := InputToBranchType(id)
        if err != nil {
            log.Warn().Println(err)
            continue
        }

        issueMap[i.Id] = name
    }

    return issueMap, nil
}
