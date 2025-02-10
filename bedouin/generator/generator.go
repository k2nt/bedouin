package generator

type Generator interface {
	Run()
	run()
}

type Field struct {
	HandleFunc func()
	IsAsync    bool
}
