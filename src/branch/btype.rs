use std::fmt;
use std::fmt::Display;
use std::str::FromStr;

#[derive(Debug, PartialEq, Eq)]
pub enum BranchType {
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
    pub fn is_unspecified(&self) -> bool {
        self == &BranchType::Unspecified
    }

    pub fn is_specified(&self) -> bool {
        self != &BranchType::Unspecified
    }
}

#[derive(Debug, PartialEq, Eq)]
pub struct BranchTypeError;

impl FromStr for BranchType {
    type Err = BranchTypeError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "build" => Ok(BranchType::Build),
            "chore" => Ok(BranchType::Chore),
            "ci" => Ok(BranchType::Ci),
            "docs" => Ok(BranchType::Docs),
            "feat" => Ok(BranchType::Feature),
            "fix" => Ok(BranchType::Fix),
            "perf" => Ok(BranchType::Performance),
            "refactor" => Ok(BranchType::Refactor),
            "revert" => Ok(BranchType::Revert),
            "style" => Ok(BranchType::Style),
            "temp" => Ok(BranchType::Temporary),
            "test" => Ok(BranchType::Test),
            _ => Err(BranchTypeError),
        }
    }
}

impl Display for BranchType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match *self {
            BranchType::Build => write!(f, "build"),
            BranchType::Chore => write!(f, "chore"),
            BranchType::Ci => write!(f, "ci"),
            BranchType::Docs => write!(f, "docs"),
            BranchType::Feature => write!(f, "feat"),
            BranchType::Fix => write!(f, "fix"),
            BranchType::Performance => write!(f, "perf"),
            BranchType::Refactor => write!(f, "refactor"),
            BranchType::Revert => write!(f, "revert"),
            BranchType::Style => write!(f, "style"),
            BranchType::Temporary => write!(f, "temp"),
            BranchType::Test => write!(f, "test"),
            BranchType::Unspecified => write!(f, ""),
        }
    }
}
