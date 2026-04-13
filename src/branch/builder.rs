#![allow(dead_code)]

use regex::{Error, Regex};
use unicode_normalization::UnicodeNormalization;

const BRANCH_TYPE_SEPARATOR: &str = "/";
pub const ISSUE_TYPE_SEPARATOR: &str = "_";
const WORD_SEPARATOR: &str = "-";

const ARTICLES: [&str; 3] = ["the", "a", "an"];

/// Categorizes branches by the primary type of change they introduce.
///
/// This enum is intended for use with conventional-commit–style workflows
/// and branch naming, where each branch represents a single dominant intent
/// (e.g. feature development, bug fixes, refactoring, or maintenance).
pub type BranchType = Option<String>;

// FIXME: this object is responsible for many things. Should be decomposed into several objects instead.
//  Current implementation almost direct copy of the Branch class in Golang. Most parts were unchanged to preserve the behavior.
/// Represents a branch with its type and rules for generating normalized names.
pub struct Branch {
    pub r#type: BranchType,
    exclude_phrases: Vec<Regex>,
    #[allow(dead_code)]
    issue_regex: Regex, // FIXME: should be separated from Branch. check branch.go L:115 ExtractIssueNameFromBranch
    strip_regex: Regex,
    pascal_case_regex: Regex,
    camel_case_regex: Regex,
}

impl Branch {
    /// Constructs a new instance with the given branch type and excluded phrases.
    pub fn new(branch_type: BranchType, exclude_phrases: Vec<&str>) -> Self {
        let exclude_phrases = build_exclude_phrases_regex_list(exclude_phrases);

        Self {
            r#type: branch_type,
            exclude_phrases,
            issue_regex: Regex::new(r"[A-Z]+-\d+_").unwrap(),
            strip_regex: build_strip_regex().unwrap(),
            pascal_case_regex: build_pascal_case_regex().unwrap(),
            camel_case_regex: build_camel_case_regex().unwrap(),
        }
    }

    /// Sets the branch type for this instance.
    #[allow(dead_code)]
    pub fn set_branch_type(&mut self, branch_type: BranchType) {
        self.r#type = branch_type;
    }

    /// Builds a normalized branch name based on the branch type, issue key, and summary.
    pub fn build_name(&self, key: &str, summary: Option<&String>) -> String {
        let mut buffer: String = String::new();

        buffer = append_branch_type(self.r#type.as_deref(), &buffer);
        buffer = append_issue_key(key, &buffer);
        buffer = append_issue_summary(
            &self.exclude_phrases,
            &self.pascal_case_regex,
            &self.camel_case_regex,
            &self.strip_regex,
            summary,
            &buffer,
        );

        buffer
    }
}

/// Builds a list of case-insensitive regex patterns for phrases to exclude.
/// Each phrase is matched if enclosed in `[]` or `()`.
///
/// Examples: [`tests::EXCLUDE_PHRASES`]
fn build_exclude_phrases_regex_list(exclude_phrases: Vec<&str>) -> Vec<Regex> {
    exclude_phrases
        .into_iter()
        .map(|w| {
            let word = regex::escape(w);
            let pattern = format!(r"(?i)(\[{word}\]|\({word}\))");
            Regex::new(&pattern).unwrap()
        })
        .collect()
}

/// Regex to strip all non-alphanumeric characters for transforming into kebab-case.
fn build_strip_regex() -> Result<Regex, Error> {
    Regex::new(r"[^a-zA-Z0-9]+")
}

/// Regex to detect `PascalCase` boundaries for transforming into kebab-case.
fn build_pascal_case_regex() -> Result<Regex, Error> {
    Regex::new(r"([A-Z]+)([A-Z][a-z])")
}

/// Regex to detect camelCase boundaries for transforming into kebab-case.
fn build_camel_case_regex() -> Result<Regex, Error> {
    Regex::new(r"([a-z])([A-Z])")
}

/// Returns a new `String` containing `buffer` followed by the processed summary, or just `buffer` if `summary` is `None`.
///
/// The summary is processed through multiple steps:
/// 1. Normalization to ASCII using [normalize].
/// 2. Removal of articles via [`filter_articles`].
/// 3. Removal of excluded phrases using [`replace_phrases`].
/// 4. Conversion from PascalCase/camelCase to kebab-case using [`pascal_camel_to_kebab`].
/// 5. Stripping unwanted characters via [strip].
fn append_issue_summary(
    exclude_phrases: &[Regex],
    pascal_case_regex: &Regex,
    camel_case_regex: &Regex,
    strip_regex: &Regex,
    summary: Option<&String>,
    buffer: &str,
) -> String {
    match summary {
        Some(text) => {
            let mut result = normalize(text);
            result = filter_articles(&result);
            result = replace_phrases(exclude_phrases, &result);
            result = pascal_camel_to_kebab(pascal_case_regex, camel_case_regex, &result);
            result = strip(strip_regex, &result);

            format!("{buffer}{result}")
        }
        None => buffer.to_owned(),
    }
}

/// Returns a new `String` containing the original `buffer` followed by the [`BranchType`]
/// and [`BRANCH_TYPE_SEPARATOR`] if the type is specified; otherwise returns the original `buffer`.
fn append_branch_type(branch_type: Option<&str>, buffer: &str) -> String {
    match branch_type {
        Some(branch_type) => format!("{buffer}{branch_type}{BRANCH_TYPE_SEPARATOR}"),
        None => buffer.to_owned(),
    }
}

/// Returns a new `String` containing issue key and [`ISSUE_TYPE_SEPARATOR`] appended to the given string.
fn append_issue_key(key: &str, buffer: &str) -> String {
    format!("{buffer}{key}{ISSUE_TYPE_SEPARATOR}")
}

/// Returns a `String` with all matching phrases removed and leading/trailing whitespace trimmed.
fn replace_phrases(exclude_phrases: &[Regex], value: &str) -> String {
    let mut subject = value.to_owned();

    for phrase in exclude_phrases {
        subject = phrase.replace_all(&subject, "").into_owned();
    }

    subject.trim().to_owned()
}

/// Converts a camelCase or `PascalCase` string to kebab-case.
///
/// The transformation is performed in two regex passes to correctly split
/// lowercase–uppercase and acronym boundaries (e.g. `HTTPServer` → `http-server`).
///
/// Note: cases when string has numeric acronym (e.g. `J2K`) are not processed intentionally.
fn pascal_camel_to_kebab(
    pascal_case_regex: &Regex,
    camel_case_regex: &Regex,
    value: &str,
) -> String {
    let kebab = pascal_case_regex.replace_all(value, |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    let kebab = camel_case_regex.replace_all(&kebab, |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    kebab.to_lowercase()
}

/// Normalizes a string by replacing unwanted characters with [`WORD_SEPARATOR`].
/// Any leading or trailing separators are removed from the final output.
fn strip(strip_regex: &Regex, value: &str) -> String {
    strip_regex
        .replace_all(value, WORD_SEPARATOR)
        .trim_start_matches(WORD_SEPARATOR)
        .trim_end_matches(WORD_SEPARATOR)
        .to_owned()
}

/// Decomposes Unicode characters using **NFKD** and filters out any non-ASCII output.
///
/// This converts accented characters (e.g. `é`) into their ASCII base (`e`) where possible,
/// drops characters without an ASCII representation.
fn normalize(value: &str) -> String {
    value.nfkd().filter(char::is_ascii).collect()
}

/// Removes standalone English articles from a string.
///
/// Words are compared case-insensitively against [ARTICLES] and removed if they match exactly.
/// Articles followed or preceded by punctuation (e.g. `an,`,`the.`) are not removed.
fn filter_articles(value: &str) -> String {
    value
        .split_whitespace()
        .filter(|word| {
            !ARTICLES
                .iter()
                .any(|article| word.eq_ignore_ascii_case(article))
        })
        .collect::<Vec<_>>()
        .join(" ")
}

#[cfg(test)]
mod tests {
    use super::*;

    // Update if not matching config/twig.toml file.
    const EXCLUDE_PHRASES: [&str; 8] = [
        "front", "mobile", "android", "ios", "be", "web", "spike", "eval",
    ];

    #[test]
    fn append_branch_type_does_nothing_when_unspecified() {
        let expected = String::new();
        let mut buffer = String::new();

        append_branch_type(None, &mut buffer);

        let actual = buffer;
        assert_eq!(expected, actual);
    }

    #[test]
    fn append_branch_type_adds_branch_type_when_specified() {
        let expected = format!("feat{}", BRANCH_TYPE_SEPARATOR);
        let actual = append_branch_type(Some("feat"), "");

        assert_eq!(expected, actual);
    }

    #[test]
    fn append_issue_key_adds_issue_key() {
        let expected = format!("XXXX-0000{}", ISSUE_TYPE_SEPARATOR);
        let actual = append_issue_key("XXXX-0000", "");

        assert_eq!(expected, actual);
    }

    #[test]
    fn replace_phrases_removes_matched_phrases() {
        let expected = String::from("Test Ticket");
        let phrases = build_exclude_phrases_regex_list(EXCLUDE_PHRASES.to_vec());

        let actual = replace_phrases(&phrases, "[Eval] (Mobile) Test Ticket");
        assert_eq!(expected, actual);
    }

    #[test]
    fn replace_phrases_outputs_unchanged_text() {
        let expected = String::from("[Unknow] (Unknown) Test_Ticket");
        let phrases = build_exclude_phrases_regex_list(EXCLUDE_PHRASES.to_vec());

        let actual = replace_phrases(&phrases, &expected);
        assert_eq!(expected, actual);
    }

    #[test]
    fn pascal_camel_to_kebab_changes_text_as_expected() {
        let expected = String::from("test-ticket-pascal");
        let pascal_case_regex = build_pascal_case_regex().unwrap();
        let camel_case_regex = build_camel_case_regex().unwrap();

        let actual =
            pascal_camel_to_kebab(&pascal_case_regex, &camel_case_regex, "TestTicketPascal");
        assert_eq!(expected, actual);
    }

    #[test]
    fn pascal_camel_to_kebab_takes_place_even_if_starts_with_lowercase() {
        let expected = String::from("lowercase-ticket-camel");
        let pascal_case_regex = build_pascal_case_regex().unwrap();
        let camel_case_regex = build_camel_case_regex().unwrap();

        let actual = pascal_camel_to_kebab(
            &pascal_case_regex,
            &camel_case_regex,
            "lowercaseTicketCamel",
        );
        assert_eq!(expected, actual);
    }

    #[test]
    fn strip_removes_undesired_symbols() {
        let expected = String::from("My-test-STRING");
        let regex = build_strip_regex().unwrap();

        let actual = strip(&regex, "My test STRING");
        assert_eq!(expected, actual);
    }

    #[test]
    fn strip_must_remove_parentheses() {
        let expected = String::from("My-test-STRING");
        let regex = build_strip_regex().unwrap();

        let actual = strip(&regex, "My (test) STRING");
        assert_eq!(expected, actual);
    }

    #[test]
    fn strip_must_remove_prefix() {
        let expected = String::from("My-test-STRING");
        let regex = build_strip_regex().unwrap();

        let actual = strip(&regex, "(My) test STRING");
        assert_eq!(expected, actual);
    }

    #[test]
    fn strip_must_remove_suffix() {
        let expected = String::from("My-test-STRING");
        let regex = build_strip_regex().unwrap();

        let actual = strip(&regex, "My test (STRING)");
        assert_eq!(expected, actual);
    }

    #[test]
    fn strip_must_treat_unknown_phrases_as_regular_text() {
        let expected = String::from("Unknow-Temp-Test-Ticket");
        let regex = build_strip_regex().unwrap();

        let actual = strip(&regex, "[Unknow] (Temp) Test Ticket");
        assert_eq!(expected, actual);
    }

    #[test]
    fn normalize_removes_non_ascii_characters() {
        let expected = String::from("ab aee");

        let actual = normalize("ab ®åe©øµ∆e");
        assert_eq!(expected, actual);
    }

    #[test]
    fn filter_articles_removes_articles_from_normalized_string() {
        let expected = String::from("Has no or or you name it");

        let actual = filter_articles("Has no the THE or thE or AN you name it");
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_constructs_correct_branch_name() {
        let expected = String::from("fix/TST-101_my-super-branch-summary");
        let branch_type = Some("fix".to_string());
        let subject = Branch::new(branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            Some(&"[Android] \"MY\" (super)_branchSummary".to_string()),
        );
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_handles_acronym_case_correctly() {
        let expected = String::from("ci/TST-101_my-super-branch-summary-http-client");
        let branch_type = Some("ci".to_string());
        let subject = Branch::new(branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            Some(&"[Android] \"MY\" (super)_branchSummary HTTPClient".to_string()),
        );
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_handles_numeric_acronym_case_correctly() {
        let expected = String::from("build/TST-101_my-super-branch-summary-j2k");
        let branch_type = Some("build".to_string());
        let subject = Branch::new(branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            Some(&"[Android] \"MY\" (super)_branchSummary J2K".to_string()),
        );
        assert_eq!(expected, actual);
    }
}
