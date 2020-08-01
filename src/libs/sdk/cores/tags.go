package cores

type TagsResult []string

type TagSource string

type ListTagsInput struct {
	Source TagSource `param:"source"`
}
