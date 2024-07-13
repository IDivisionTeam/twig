package issue

type Type struct {
    Id   string
    Name string
}

// TODO: replace with your <issue-types>
const (
    Build    = "1000"
    Chore    = "2000"
    Ci       = "3000"
    Docs     = "4000"
    Feat     = "5000"
    Fix      = "6000"
    Perf     = "7000"
    Refactor = "8000"
    Revert   = "9000"
    Style    = "1010"
    Test     = "1100"
)

// TODO: replace with your ignored <issue-types>
var Ignored = buildMap()

func buildMap() map[string]bool {
    builder := make(map[string]bool)

    builder.Set("4444", true) // Subtask
    // builder.Set("3333", true) // Bug
    // etc.

    return builder
}
