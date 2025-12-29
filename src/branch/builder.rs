#![allow(dead_code)]

use crate::branch::btype::BranchType;
use regex::{Error, Regex};
use std::str::FromStr;
use unicode_normalization::UnicodeNormalization;

const BRANCH_TYPE_SEPARATOR: &str = "/";
const ISSUE_TYPE_SEPARATOR: &str = "_";
const WORD_SEPARATOR: &str = "-";

const ARTICLES: [&str; 3] = ["the", "a", "an"];

// FIXME: this object is responsible for many things. Should be decomposed into several objects instead.
//  Current implementation almost direct copy of the Branch class in Golang. Most parts were unchanged to preserve the behavior.
struct Branch {
    pub branch_type: BranchType,
    exclude_phrases: Vec<Regex>,
    #[allow(dead_code)]
    issue_regex: Regex, // FIXME: should be separated from Branch. check branch.go L:115 ExtractIssueNameFromBranch
    strip_regex: Regex,
    first_pass_kebab_regex: Regex,
    second_pass_kebab_regex: Regex,
}

impl Branch {
    pub fn new(branch_type: &str, exclude_phrases: Vec<&str>) -> Self {
        let branch_type = BranchType::from_str(branch_type).unwrap_or(BranchType::Unspecified);
        let exclude_phrases = build_exclude_phrases_regex_list(exclude_phrases);

        Self {
            branch_type,
            exclude_phrases,
            issue_regex: Regex::new(r"[A-Z]+-\d+_").unwrap(),
            strip_regex: build_strip_regex().unwrap(),
            first_pass_kebab_regex: build_first_pass_kebab_regex().unwrap(),
            second_pass_kebab_regex: build_second_pass_kebab_regex().unwrap(),
        }
    }

    #[allow(dead_code)]
    pub fn set_branch_type(&mut self, branch_type: BranchType) {
        self.branch_type = branch_type;
    }

    pub fn build_name(&self, key: &str, summary: &Option<String>) -> String {
        let mut buffer: String = String::new();

        append_branch_type(&self.branch_type, &mut buffer);
        append_issue_key(key, &mut buffer);
        append_issue_summary(
            &self.exclude_phrases,
            &self.first_pass_kebab_regex,
            &self.second_pass_kebab_regex,
            &self.strip_regex,
            summary,
            &mut buffer,
        );

        buffer
    }
}

fn build_exclude_phrases_regex_list(exclude_phrases: Vec<&str>) -> Vec<Regex> {
    exclude_phrases
        .into_iter()
        .map(|w| {
            let word = regex::escape(w);
            let pattern = format!(r"(?i)(\[{}\]|\({}\))", word, word);
            Regex::new(&pattern).unwrap()
        })
        .collect()
}

fn build_strip_regex() -> Result<Regex, Error> {
    Regex::new(r"[^a-zA-Z0-9]+")
}

fn build_first_pass_kebab_regex() -> Result<Regex, Error> {
    Regex::new(r"([A-Z]+)([A-Z][a-z])")
}

fn build_second_pass_kebab_regex() -> Result<Regex, Error> {
    Regex::new(r"([a-z])([A-Z])")
}

fn append_issue_summary(
    exclude_phrases: &[Regex],
    first_pass_kebab_regex: &Regex,
    second_pass_kebab_regex: &Regex,
    strip_regex: &Regex,
    summary: &Option<String>,
    buffer: &mut String,
) {
    if let Some(text) = summary {
        let mut result = normalize(text);
        result = tokenize(&result);
        result = replace_phrases(exclude_phrases, &result);
        result = camel_to_kebab(first_pass_kebab_regex, second_pass_kebab_regex, &result);
        result = strip(strip_regex, &result);

        buffer.push_str(&result);
    }
}

fn append_branch_type(branch_type: &BranchType, buffer: &mut String) {
    if branch_type.is_specified() {
        buffer.push_str(branch_type.to_string().as_str());
        buffer.push_str(BRANCH_TYPE_SEPARATOR);
    }
}

fn append_issue_key(key: &str, buffer: &mut String) {
    buffer.push_str(key);
    buffer.push_str(ISSUE_TYPE_SEPARATOR);
}

fn replace_phrases(exclude_phrases: &[Regex], tokenized: &str) -> String {
    let mut subject = tokenized.to_owned();

    for phrase in exclude_phrases.iter() {
        subject = phrase.replace_all(&subject, "").into_owned();
    }

    subject.trim().to_owned()
}

fn camel_to_kebab(
    first_pass_kebab_regex: &Regex,
    second_pass_kebab_regex: &Regex,
    camel_str: &str,
) -> String {
    let kebab = first_pass_kebab_regex.replace_all(camel_str, |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    let kebab = second_pass_kebab_regex.replace_all(&kebab, |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    kebab.to_lowercase()
}

fn strip(strip_regex: &Regex, str: &str) -> String {
    strip_regex
        .replace_all(str, WORD_SEPARATOR)
        .trim_start_matches(WORD_SEPARATOR)
        .trim_end_matches(WORD_SEPARATOR)
        .to_string()
}

fn normalize(summary: &str) -> String {
    summary.nfkd().filter(|c| c.is_ascii()).collect()
}

fn tokenize(normalized: &str) -> String {
    normalized
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

        append_branch_type(&BranchType::Unspecified, &mut buffer);

        let actual = buffer;
        assert_eq!(expected, actual);
    }

    #[test]
    fn append_branch_type_adds_branch_type_when_specified() {
        let expected = format!("feat{}", BRANCH_TYPE_SEPARATOR);
        let mut buffer = String::new();

        append_branch_type(&BranchType::Feature, &mut buffer);

        let actual = buffer;
        assert_eq!(expected, actual);
    }

    #[test]
    fn append_issue_key_adds_issue_key() {
        let expected = format!("XXXX-0000{}", ISSUE_TYPE_SEPARATOR);
        let mut buffer = String::new();

        append_issue_key("XXXX-0000", &mut buffer);

        let actual = buffer;
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
    fn camel_to_kebab_changes_text_as_expected() {
        let expected = String::from("test-ticket");
        let first_pass = build_first_pass_kebab_regex().unwrap();
        let second_pass = build_second_pass_kebab_regex().unwrap();

        let actual = camel_to_kebab(&first_pass, &second_pass, "TestTicket");
        assert_eq!(expected, actual);
    }

    #[test]
    fn camel_to_kebab_takes_place_even_if_starts_with_lowercase() {
        let expected = String::from("lowercase-ticket-camel");
        let first_pass = build_first_pass_kebab_regex().unwrap();
        let second_pass = build_second_pass_kebab_regex().unwrap();

        let actual = camel_to_kebab(&first_pass, &second_pass, "lowercaseTicketCamel");
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
    fn tokenize_splits_normalized_string() {
        let expected = String::from("string is typical thing");

        let actual = tokenize("The string is a typical thing");
        assert_eq!(expected, actual);
    }

    #[test]
    fn tokenize_removes_articles_from_normalized_string() {
        let expected = String::from("Has no or or you name it");

        let actual = tokenize("Has no the THE or thE or AN you name it");
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_constructs_correct_branch_name() {
        let expected = String::from("fix/TST-101_my-super-branch-summary");
        let branch_type = BranchType::Fix.to_string();
        let subject = Branch::new(&branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            &Some("[Android] \"MY\" (super)_branchSummary".to_string()),
        );
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_handles_acronym_case_correctly() {
        let expected = String::from("ci/TST-101_my-super-branch-summary-http-client");
        let branch_type = BranchType::Ci.to_string();
        let subject = Branch::new(&branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            &Some("[Android] \"MY\" (super)_branchSummary HTTPClient".to_string()),
        );
        assert_eq!(expected, actual);
    }

    #[test]
    fn branch_build_name_handles_numeric_acronym_case_correctly() {
        let expected = String::from("build/TST-101_my-super-branch-summary-j2k");
        let branch_type = BranchType::Build.to_string();
        let subject = Branch::new(&branch_type, EXCLUDE_PHRASES.to_vec());

        let actual = subject.build_name(
            "TST-101",
            &Some("[Android] \"MY\" (super)_branchSummary J2K".to_string()),
        );
        assert_eq!(expected, actual);
    }
}
