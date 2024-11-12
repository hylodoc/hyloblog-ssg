package sitefile

type CustomPage interface {
	Template() string
	Data() map[string]string
}
