package sitefile

type CustomPage interface {
	Title() string
	Content() string
}
