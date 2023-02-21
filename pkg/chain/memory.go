package chain

type Memory interface {
	Clear()
	SaveContext(input string, output string)
	LoadMemoryVariables() map[string]string
}
