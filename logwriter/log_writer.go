package logwriter

import log "github.com/sirupsen/logrus"

type LogWriter struct {
	Message string
}

func (l LogWriter) Write(p []byte) (n int, err error) {
	if log.GetLevel() == log.TraceLevel || log.GetLevel() == log.DebugLevel {
		log.Debugf("%s: %s", l.Message, string(p))
		return len(p), nil
	}
	return 0, nil
}
