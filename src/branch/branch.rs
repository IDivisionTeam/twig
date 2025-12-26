use crate::branch::btype::BranchType;
use regex::Regex;
use unicode_normalization::UnicodeNormalization;

const BRANCH_TYPE_SEPARATOR: &str = "/";
const ISSUE_TYPE_SEPARATOR: &str = "_";
const WORD_SEPARATOR: &str = "-";

const ARTICLES: [&str; 3] = ["the", "a", "an"];

struct Branch {
    branch_type: BranchType,
    exclude_phrases: Vec<Regex>,
    #[allow(dead_code)]
    issue_regex: Regex, // FIXME: should be separated from Branch. check branch.go L:115 ExtractIssueNameFromBranch
    strip_regex: Regex,
    first_pass_kebab_regex: Regex,
    second_pass_kebab_regex: Regex,
}

impl Branch {
    pub fn new(branch_type: String, exclude_phrases: Vec<String>) -> Self {
        let bt = BranchType::from_string(branch_type).unwrap_or(BranchType::Unspecified);
        let phrases = exclude_phrases
            .into_iter()
            .map(|w| {
                let word = regex::escape(w.as_str());
                let pattern = format!(r"(?i)(\\[{}\\]|\\({}\\))", word, word);
                Regex::new(&pattern).unwrap()
            })
            .collect();

        Self {
            branch_type: bt,
            exclude_phrases: phrases,
            issue_regex: Regex::new(r"[A-Z]+-\d+_").unwrap(),
            strip_regex: Regex::new(r"[^a-zA-Z0-9]+").unwrap(),
            first_pass_kebab_regex: Regex::new(r"([A-Z]+)([A-Z][a-z])").unwrap(),
            second_pass_kebab_regex: Regex::new(r"([a-z])([A-Z])").unwrap(),
        }
    }

    fn append_branch_type(&self, branch_type: &BranchType, buffer: &mut String) {
        if branch_type.is_specified() {
            buffer.push_str(branch_type.to_string().as_str());
            buffer.push_str(BRANCH_TYPE_SEPARATOR);
        }
    }

    fn append_issue_key(&self, key: String, buffer: &mut String) {
        buffer.push_str(key.as_str());
        buffer.push_str(ISSUE_TYPE_SEPARATOR);
    }

    fn append_issue_summary(&self, summary: Option<String>, buffer: &mut String) {
        if let Some(text) = summary {
            let mut result = self.normalize_summary(text);
            result = self.tokenize(result);
            result = self.replace_phrases(result);
            result = self.camel_to_kebab(result);
            result = self.strip_phrases(result);

            buffer.push_str(&result);
        }
    }

    fn replace_phrases(&self, tokenized: String) -> String {
        let mut subject = tokenized;

        for phrase in self.exclude_phrases.iter() {
            let result = phrase.replace_all(subject.as_str(), "");
            subject = result.into_owned();
        }

        subject
    }

    fn camel_to_kebab(&self, camel_str: String) -> String {
        let kebab = self
            .first_pass_kebab_regex
            .replace_all(camel_str.as_str(), |caps: &regex::Captures| {
                format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
            })
            .into_owned();

        let kebab = self
            .second_pass_kebab_regex
            .replace_all(&kebab, |caps: &regex::Captures| {
                format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
            })
            .into_owned();

        kebab.to_lowercase()
    }

    fn strip_phrases(&self, str: String) -> String {
        self.strip_regex
            .replace_all(&str, WORD_SEPARATOR)
            .trim_start_matches(WORD_SEPARATOR)
            .trim_end_matches(WORD_SEPARATOR)
            .to_string()
    }

    pub fn build_name(&self, key: String, summary: Option<String>) -> String {
        let mut buffer: String = String::new();

        self.append_branch_type(&self.branch_type, &mut buffer);
        self.append_issue_key(key, &mut buffer);
        self.append_issue_summary(summary, &mut buffer);

        buffer
    }

    fn normalize_summary(&self, summary: String) -> String {
        summary.nfkd().filter(|c| c.is_ascii()).collect()
    }

    fn tokenize(&self, normalized: String) -> String {
        normalized
            .split_whitespace()
            .filter(|word| {
                !ARTICLES
                    .iter()
                    .any(|article| word.eq_ignore_ascii_case(article))
            })
            .fold(String::new(), |acc, next| {
                if acc.is_empty() {
                    return next.to_string();
                }
                acc + " " + next
            })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn normalize_summary_removes_non_ascii_characters() {
        let branch = Branch::new("".to_string(), vec!["test".to_string()]);
        let summary = branch.build_name("XXXX-0000".to_string(), Some("ab ®åe©øµ∆e".to_string()));
        assert_eq!("XXXX-0000_ab-aee", summary);
    }

    #[test]
    fn tokenize_splits_normalized_string() {
        let branch = Branch::new("".to_string(), vec!["test".to_string()]);
        let summary = branch.build_name(
            "XXXX-0000".to_string(),
            Some("The string is a typical thing".to_string()),
        );
        assert_eq!("XXXX-0000_string-is-typical-thing", summary);
    }

    #[test]
    fn tokenize_removes_articles_from_normalized_string() {
        let branch = Branch::new("".to_string(), vec!["test".to_string()]);
        let summary = branch.build_name(
            "XXXX-0000".to_string(),
            Some("Has no the THE or thE or AN you name it".to_string()),
        );
        assert_eq!("XXXX-0000_has-no-or-or-you-name-it", summary);
    }

    // TODO(BR-45): port remaining tests from branch.go
}
