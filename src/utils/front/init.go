package front

func GenerateAllFront() error {
	if err := GenerateIndex(); err != nil {
		return err
	}
	if err := GeneratePages(); err != nil {
		return err
	}

	return nil
}
