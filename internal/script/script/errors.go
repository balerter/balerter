package script

type CronAlreadyDefinedError struct {
	prevSpec string
}

func newCronAlreadyDefinedError(prevSpec string) CronAlreadyDefinedError {
	return CronAlreadyDefinedError{
		prevSpec: prevSpec,
	}
}

func (err CronAlreadyDefinedError) Error() string {
	return "cron already defined: " + err.prevSpec
}
