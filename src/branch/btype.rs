#[derive(Debug)]
#[allow(dead_code)]
enum BranchType {
    Build,
    Chore,
    Ci,
    Docs,
    Feature,
    Fix,
    Performance,
    Refactor,
    Revert,
    Style,
    Temporary,
    Test,
    Unspecified,
}

impl BranchType {
    #[allow(dead_code)]
    fn to_string(&self) -> String {
        match self {
            BranchType::Build => String::from("build"),
            BranchType::Chore => String::from("chore"),
            BranchType::Ci => String::from("ci"),
            BranchType::Docs => String::from("docs"),
            BranchType::Feature => String::from("feat"),
            BranchType::Fix => String::from("fix"),
            BranchType::Performance => String::from("perf"),
            BranchType::Refactor => String::from("refactor"),
            BranchType::Revert => String::from("revert"),
            BranchType::Style => String::from("style"),
            BranchType::Temporary => String::from("temp"),
            BranchType::Test => String::from("test"),
            BranchType::Unspecified => String::from(""),
        }
    }

    #[allow(dead_code)]
    fn from_string(input: String) -> Option<BranchType> {
        match input.as_str() {
            "build" => Some(BranchType::Build),
            "chore" => Some(BranchType::Chore),
            "ci" => Some(BranchType::Ci),
            "docs" => Some(BranchType::Docs),
            "feat" => Some(BranchType::Feature),
            "fix" => Some(BranchType::Fix),
            "perf" => Some(BranchType::Performance),
            "refactor" => Some(BranchType::Refactor),
            "revert" => Some(BranchType::Revert),
            "style" => Some(BranchType::Style),
            "temp" => Some(BranchType::Temporary),
            "test" => Some(BranchType::Test),
            _ => None,
        }
    }
}
