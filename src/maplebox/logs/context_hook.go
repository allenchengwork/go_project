package logs

import (
	"gopkg.in/Sirupsen/logrus.v0"
	"maplebox/util"
	"path"
)

type ContextHook struct {
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	file, line := util.GetCallerIgnoringLogMulti(2)

	entry.Data["file"] = path.Base(file)
	entry.Data["line"] = line

	return nil
}
