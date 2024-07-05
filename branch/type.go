package branch

type Type struct {
    Build    string
    Chore    string
    Ci       string
    Docs     string
    Feat     string
    Fix      string
    Perf     string
    Refactor string
    Revert   string
    Style    string
    Test     string
}

func NewBranchType() *Type {
    return &Type{
        Build:    "build",
        Chore:    "chore",
        Ci:       "ci",
        Docs:     "docs",
        Feat:     "feat",
        Fix:      "fix",
        Perf:     "perf",
        Refactor: "refactor",
        Revert:   "revert",
        Style:    "style",
        Test:     "test",
    }
}
