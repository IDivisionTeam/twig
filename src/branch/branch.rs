use crate::branch::btype::BranchType;
use regex::Regex;
use unicode_normalization::UnicodeNormalization;

const BRANCH_TYPE_SEPARATOR: &str = "/";
const ISSUE_TYPE_SEPARATOR: &str = "_";
const WORD_SEPARATOR: &str = "-";

const ARTICLES: [&str; 3] = ["the", "a", "an"];

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
    pub fn new(branch_type: String, exclude_phrases: Vec<String>) -> Self {
        let bt = BranchType::from_string(branch_type).unwrap_or(BranchType::Unspecified);
        let phrases = build_exclude_phrases_regex_list(exclude_phrases);

        Self {
            branch_type: bt,
            exclude_phrases: phrases,
            issue_regex: Regex::new(r"[A-Z]+-\d+_").unwrap(),
            strip_regex: Regex::new(r"[^a-zA-Z0-9]+").unwrap(),
            first_pass_kebab_regex: Regex::new(r"([A-Z]+)([A-Z][a-z])").unwrap(),
            second_pass_kebab_regex: Regex::new(r"([a-z])([A-Z])").unwrap(),
        }
    }

    #[allow(dead_code)]
    pub fn set_branch_type(&mut self, branch_type: BranchType) {
        self.branch_type = branch_type;
    }

    pub fn build_name(&self, key: &String, summary: &Option<String>) -> String {
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

fn build_exclude_phrases_regex_list(exclude_phrases: Vec<String>) -> Vec<Regex> {
    exclude_phrases
        .into_iter()
        .map(|w| {
            let word = regex::escape(w.as_str());
            let pattern = format!(r"(?i)(\[{}\]|\({}\))", word, word);
            Regex::new(&pattern).unwrap()
        })
        .collect()
}

fn append_issue_summary(
    exclude_phrases: &Vec<Regex>,
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
        result = strip_phrases(strip_regex, &result);

        buffer.push_str(&result);
    }
}

fn append_branch_type(branch_type: &BranchType, buffer: &mut String) {
    if branch_type.is_specified() {
        buffer.push_str(branch_type.to_string().as_str());
        buffer.push_str(BRANCH_TYPE_SEPARATOR);
    }
}

fn append_issue_key(key: &String, buffer: &mut String) {
    buffer.push_str(key.as_str());
    buffer.push_str(ISSUE_TYPE_SEPARATOR);
}

fn replace_phrases(exclude_phrases: &Vec<Regex>, tokenized: &String) -> String {
    let mut subject = tokenized.to_owned();

    for phrase in exclude_phrases.iter() {
        subject = phrase.replace_all(&subject, "").into_owned();
    }

    subject.trim().to_owned()
}

fn camel_to_kebab(
    first_pass_kebab_regex: &Regex,
    second_pass_kebab_regex: &Regex,
    camel_str: &String,
) -> String {
    let kebab = first_pass_kebab_regex.replace_all(camel_str.as_str(), |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    let kebab = second_pass_kebab_regex.replace_all(&kebab, |caps: &regex::Captures| {
        format!("{}{}{}", &caps[1], WORD_SEPARATOR, &caps[2])
    });

    kebab.to_lowercase()
}

fn strip_phrases(strip_regex: &Regex, str: &String) -> String {
    strip_regex
        .replace_all(str, WORD_SEPARATOR)
        .trim_start_matches(WORD_SEPARATOR)
        .trim_end_matches(WORD_SEPARATOR)
        .to_string()
}

fn normalize(summary: &String) -> String {
    summary.nfkd().filter(|c| c.is_ascii()).collect()
}

fn tokenize(normalized: &String) -> String {
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

#[cfg(test)]
mod tests {
    use super::*;
    use lazy_static::lazy_static;

    lazy_static! {
        // Update if not matching config/twig.toml file.
        static ref EXCLUDE_PHRASES: Vec<String> = {
             [
                "front",
                "mobile",
                "android",
                "ios",
                "be",
                "web",
                "spike",
                "eval"
             ].map(String::from).to_vec()
        };
    }

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

        append_issue_key(&String::from("XXXX-0000"), &mut buffer);

        let actual = buffer;
        assert_eq!(expected, actual);
    }

    #[test]
    fn replace_phrases_removes_matched_phrases() {
        let expected = String::from("Test Ticket");
        let phrases = build_exclude_phrases_regex_list(EXCLUDE_PHRASES.to_vec());

        let actual = replace_phrases(&phrases, &"[Eval] (Mobile) Test Ticket".to_string());
        assert_eq!(expected, actual);
    }

    // TODO(BR-45): add tests.
    // camel_to_kebab
    // strip_phrases

    #[test]
    fn normalize_removes_non_ascii_characters() {
        let expected = String::from("ab aee");

        let actual = normalize(&"ab ®åe©øµ∆e".to_string());
        assert_eq!(expected, actual);
    }

    #[test]
    fn tokenize_splits_normalized_string() {
        let expected = String::from("string is typical thing");

        let actual = tokenize(&"The string is a typical thing".to_string());
        assert_eq!(expected, actual);
    }

    #[test]
    fn tokenize_removes_articles_from_normalized_string() {
        let expected = String::from("Has no or or you name it");

        let actual = tokenize(&"Has no the THE or thE or AN you name it".to_string());
        assert_eq!(expected, actual);
    }

    // TODO(BR-45): port remaining tests from branch.go
}
