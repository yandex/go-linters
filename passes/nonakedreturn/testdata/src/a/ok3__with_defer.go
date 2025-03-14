package a

func okDeferred() error {
	return nil
}

func okSomeFunction() error {
	return nil
}

func okWithDefer() (msg string, err error) {
	defer func() {
		if deferredErr := okDeferred(); deferredErr != nil {
			err = deferredErr
		}
	}()

	err = okSomeFunction()

	return "hello world", err
}
