package nsync

type Waiter interface {
	Add(delta int)

	Done()

	Wait()
}
