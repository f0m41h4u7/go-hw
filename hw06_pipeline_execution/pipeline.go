package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		empty := make(Bi)
		close(empty)
		return empty
	}

	res := make(Bi)

	pipeline := make(In)
	pipeline = stages[0](in)
	for i := 1; i < len(stages); i++ {
		pipeline = stages[i](pipeline)
	}

	go func() {
		defer close(res)
		for {
			select {
			case val, ok := <-pipeline:
				if !ok {
					return
				}
				res <- val
			case <-done:
				return
			}
		}
	}()

	return res
}
