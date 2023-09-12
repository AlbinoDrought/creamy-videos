package tmpl

type PublicURLGenerator func(relativeURL string) string

type AppState struct {
	ReadOnly      bool
	SortDirection string
	Sortable      bool
	SearchText    string

	PUG PublicURLGenerator
}
