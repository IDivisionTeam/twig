package branch

type Type int

const (
    NULL Type = iota
    BUILD
    CHORE
    CI
    DOCS
    FEAT
    FIX
    PERF
    REFACTOR
    REVERT
    STYLE
    TEST
)

const (
    Build         = "build"
    BuildShort    = "b"
    Chore         = "chore"
    ChoreShort    = "ch"
    Ci            = "ci"
    Docs          = "docs"
    DocsShort     = "d"
    Feat          = "feat"
    FeatShort     = "ft"
    Fix           = "fix"
    FixShort      = "fx"
    Perf          = "perf"
    PerfShort     = "p"
    Refactor      = "refactor"
    RefactorShort = "rf"
    Revert        = "revert"
    RevertShort   = "rv"
    Style         = "style"
    StyleShort    = "s"
    Test          = "test"
    TestShort     = "t"
)

func (t Type) ToString() string {
    switch t {
    case BUILD:
        return Build
    case CHORE:
        return Chore
    case CI:
        return Ci
    case DOCS:
        return Docs
    case FEAT:
        return Feat
    case FIX:
        return Fix
    case PERF:
        return Perf
    case REFACTOR:
        return Refactor
    case REVERT:
        return Revert
    case STYLE:
        return Style
    case TEST:
        return Test
    default:
        return ""
    }
}
