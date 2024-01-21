package client

type Project struct {
	Key string `json:"key"`
}

type Paginated[T any] struct {
	Values []T `json:"values"`
}

type Repository struct {
	Slug      string  `json:"slug"`
	IsPrivate bool    `json:"is_private"`
	Scm       string  `json:"scm"`
	Project   Project `json:"project"`
	Website   *string `json:"website"`
}

type User struct {
	Uuid string `json:"uuid"`
}

type Group struct {
	Slug string `json:"slug"`
}

type BranchRestriction struct {
	ID              int     `json:"id"`
	Kind            string  `json:"kind"`
	BranchMatchKind string  `json:"branch_match_kind"`
	Pattern         string  `json:"pattern"`
	Users           []User  `json:"users"`
	Groups          []Group `json:"groups"`
}

type GroupPermission struct {
	Permission string `json:"permission"`
}
