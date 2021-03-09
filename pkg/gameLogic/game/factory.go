package game

type WorkerProvider interface {
	GetConstructor(workerType string) (WorkerConstructor, error)
}

type WorkerFactory struct{}

func (wf WorkerFactory) GetConstructor(workerType string) (WorkerConstructor, error) {
	switch workerType {
	case "TopOfTheTop":
		return newTopOfTheTopWorker, nil
	default:
		return nil, UnknownGameTypeErr
	}
}
