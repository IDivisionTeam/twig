use std::fmt;
use std::fmt::Display;
use std::str::FromStr;

/// Categorizes branches by the primary type of change they introduce.
///
/// This enum is intended for use with conventional-commit–style workflows
/// and branch naming, where each branch represents a single dominant intent
/// (e.g. feature development, bug fixes, refactoring, or maintenance).
#[derive(Debug, PartialEq, Eq)]
pub enum BranchType {
    /// Changes that affect the build system or external dependencies.
    Build,
    /// Routine maintenance tasks that do not modify source or test files.
    Chore,
    /// Changes to the CI configuration files and scripts.
    Ci,
    /// Documentation only changes.
    Docs,
    /// A new feature.
    Feature,
    /// A bug fix.
    Fix,
    /// A code change that improves performance.
    Performance,
    /// A code change that neither fixes a bug nor adds a feature.
    Refactor,
    /// Reverts a previous commit(s).
    Revert,
    /// Changes that do not affect the meaning of the code (white-space, formatting, missing semicolons, etc.).
    Style,
    /// Temporary or experimental changes not intended for long-term use.
    Temporary,
    /// Adding missing tests or correcting existing tests.
    Test,
    /// Changes that do not clearly fit any other category or undefined.
    Unspecified,
}

impl BranchType {
    /// Returns `true` if the branch type is `Unspecified`.
    #[allow(dead_code)]
    pub fn is_unspecified(&self) -> bool {
        *self == BranchType::Unspecified
    }

    /// Returns `true` if the branch type is anything other than `Unspecified`.
    pub fn is_specified(&self) -> bool {
        *self != BranchType::Unspecified
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
