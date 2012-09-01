package components

func NewEngine(kind string, limit byte) *Component {
	return &Component{make(chan Flow), make(map[int]FlowWriter), EngineFuncGenerator(kind, limit)}
}

func NewGenerator(kind string, amount byte) *Component {
	return &Component{make(chan Flow), make(map[int]FlowWriter), GeneratorFuncGenerator(kind, amount)}
}

func EngineFuncGenerator(kind string, limit byte) (func(Flow) Flow) {
	return func(in Flow) Flow {
		if in.kind != kind || in.amount > limit {
			//BLOW SHI(T|P) UP!
		}
		//Propell ship i direction
		return Flow{kind, 0}
	}
}

func GeneratorFuncGenerator(kind string, amount byte) (func(Flow) Flow) {
	return func(in Flow) Flow {
		return Flow{kind, amount}
	}
}
