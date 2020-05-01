package script

type CronAlreadyDefinedError struct {
	PrevSpec string
}

func newCronAlreadyDefinedError(prevSpec string) CronAlreadyDefinedError {
	return CronAlreadyDefinedError{
		PrevSpec: prevSpec,
	}
}

func (err CronAlreadyDefinedError) Error() string {
	return "cron already defined: " + err.PrevSpec
}
