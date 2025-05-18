package branch

import (
    "testing"
    "twig/log"
    "twig/network"
)

// Update if not matching config/twig.toml file.
var phrases []string

func init() {
    log.CreateNoOpTestRecorders()

    phrases = []string{"front", "mobile", "android", "ios", "be", "web", "spike", "eval"}
}

func TestReplacePhrasesOptimisticCase(t *testing.T) {
    in := "[Eval] (Mobile) Test Ticket"

    want := "Test Ticket"
    subject := New(NULL, phrases).replacePhrases(in)

    if subject != want {
        t.Errorf(`replacePhrases(in, rawPhrases) = %q, want match for %q`, subject, want)
    }
}

func TestReplacePhrasesNoChanges(t *testing.T) {
    in := "[Unknow] (Unknown) Test_Ticket"

    want := "[Unknow] (Unknown) Test_Ticket"
    subject := New(NULL, phrases).replacePhrases(in)

    if subject != want {
        t.Errorf(`replacePhrases(in, rawPhrases) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexOptimisticCase(t *testing.T) {
    in := "My test STRING"

    want := "my-test-string"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexWithParentheses(t *testing.T) {
    in := "My (test) STRING"

    want := "my-test-string"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexWithSpaceInFront(t *testing.T) {
    in := " My test STRING"

    want := "my-test-string"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexWithSuffix(t *testing.T) {
    in := "My test (STRING)"

    want := "my-test-string"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexWithPrefix(t *testing.T) {
    in := "(My) test STRING"

    want := "my-test-string"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestStripRegexUnknownPhrases(t *testing.T) {
    in := "[Unknow] (Temp) Test Ticket"

    want := "unknow-temp-test-ticket"
    subject := New(NULL, nil).stripRegex(in)

    if subject != want {
        t.Errorf(`stripRegex(in) = %q, want match for %q`, subject, want)
    }
}

func TestCamelToKebabOptimisticCase(t *testing.T) {
    in := "TestTicket"

    want := "test-ticket"
    subject := New(NULL, nil).camelToKebab(in)

    if subject != want {
        t.Errorf(`camelToKebab(in) = %q, want match for %q`, subject, want)
    }
}

func TestCamelToKebabStartWithLowercase(t *testing.T) {
    in := "lowercaseTicketCamel"

    want := "lowercase-ticket-camel"
    subject := New(NULL, nil).camelToKebab(in)

    if subject != want {
        t.Errorf(`camelToKebab(in) = %q, want match for %q`, subject, want)
    }
}

func TestBuildNameOptimisticCase(t *testing.T) {
    branchType := FIX
    summary := "[Android] \"MY\" (super)_branchSummary"
    issue := network.JiraIssue{
        Id:  "",
        Key: "TST-101",
        Fields: network.IssueFields{
            Type: &network.IssueType{
                Id:   "",
                Name: "",
            },
            Summary:  &summary,
            Status:   nil,
            Assignee: nil,
        },
    }

    want := "fix/TST-101_my-super-branch-summary"
    subject := New(branchType, phrases).BuildName(issue)

    if subject != want {
        t.Errorf(`BuildName(type, issue, phrases) = %q, want match for %q`, subject, want)
    }
}

func TestBuildNameAcronym(t *testing.T) {
    branchType := FIX
    summary := "[Android] \"MY\" (super)_branchSummary HTTPClient"
    issue := network.JiraIssue{
        Id:  "",
        Key: "TST-101",
        Fields: network.IssueFields{
            Type: &network.IssueType{
                Id:   "",
                Name: "",
            },
            Summary:  &summary,
            Status:   nil,
            Assignee: nil,
        },
    }

    want := "fix/TST-101_my-super-branch-summary-http-client"
    subject := New(branchType, phrases).BuildName(issue)

    if subject != want {
        t.Errorf(`BuildName(type, issue, phrases) = %q, want match for %q`, subject, want)
    }
}

func TestBuildNameNumbersInBetween(t *testing.T) {
    branchType := FIX
    summary := "[Android] \"MY\" (super)_branchSummary J2K"
    issue := network.JiraIssue{
        Id:  "",
        Key: "TST-101",
        Fields: network.IssueFields{
            Type: &network.IssueType{
                Id:   "",
                Name: "",
            },
            Summary:  &summary,
            Status:   nil,
            Assignee: nil,
        },
    }

    want := "fix/TST-101_my-super-branch-summary-j2k"
    subject := New(branchType, phrases).BuildName(issue)

    if subject != want {
        t.Errorf(`BuildName(type, issue, phrases) = %q, want match for %q`, subject, want)
    }
}
