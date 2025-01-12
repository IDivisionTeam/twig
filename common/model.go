package common

type InputType int
type Flag uint8

const (
    Issue InputType = iota
    BranchType
    Remote
    Assignee
)

const (
    HelpFlag Flag = 1 << iota
    CleanFlag
)

const EmptyFlag Flag = 0

type Input struct {
    Flags     Flag
    Arguments map[InputType]string
}

func (i *Input) AddFlag(flag Flag) {
    i.Flags |= flag
}

func (i *Input) RemoveFlag(flag Flag) {
    i.Flags &= ^flag
}

func (i *Input) HasFlag(flag Flag) bool {
    return i.Flags&flag != EmptyFlag
}
