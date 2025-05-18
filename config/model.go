package config

type Settings struct {
	Project ProjectSettings `mapstructure:"project"`
	Branch  BranchSettings  `mapstructure:"branch"`
	Mapping MappingSettings `mapstructure:"mapping"`
}

type ProjectSettings struct {
	Host  string `mapstructure:"host"`
	Auth  string `mapstructure:"auth"`
	Email string `mapstructure:"email"`
	Token string `mapstructure:"token"`
}

type BranchSettings struct {
	Default string   `mapstructure:"default"`
	Origin  string   `mapstructure:"origin"`
	Exclude []string `mapstructure:"exclude"`
}

type MappingSettings struct {
	Build    []string `mapstructure:"build"`
	Chore    []string `mapstructure:"chore"`
	Ci       []string `mapstructure:"ci"`
	Docs     []string `mapstructure:"docs"`
	Feat     []string `mapstructure:"feat"`
	Fix      []string `mapstructure:"fix"`
	Pref     []string `mapstructure:"pref"`
	Refactor []string `mapstructure:"refactor"`
	Revert   []string `mapstructure:"revert"`
	Style    []string `mapstructure:"style"`
	Test     []string `mapstructure:"test"`
}
