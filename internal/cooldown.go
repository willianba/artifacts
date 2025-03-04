package internal

import "time"

func Cooldown(seconds int) {
	Logger.Info("Waiting after last action", "cooldown", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
}
