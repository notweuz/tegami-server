package util

import "runtime/debug"

func GetCommit() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return TruncateString(setting.Value, 7)
			}
		}
	}
	return "no git revision detected"
}
