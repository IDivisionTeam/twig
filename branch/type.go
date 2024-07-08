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

func (t Type) ToString() string {
    switch t {
    case BUILD:
        return "build"
    case CHORE:
        return "chore"
    case CI:
        return "ci"
    case DOCS:
        return "docs"
    case FEAT:
        return "feat"
    case FIX:
        return "fix"
    case PERF:
        return "perf"
    case REFACTOR:
        return "refactor"
    case REVERT:
        return "revert"
    case STYLE:
        return "style"
    case TEST:
        return "test"
    default:
        return ""
    }
}
