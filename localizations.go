package main

var localizations = map[string]string{
	"cmd.packs.title":       "packs",
	"cmd.packs.description": "See and boot available packs",

	"selector.title":             "switch 'n' craft!",
	"selector.titleSwitch":       "Not what you're looking for? Just switch!",
	"selector.desc":              "What do you wanna play today? Select a pack and get playing.",
	"selector.count":             "%d pack/s available",
	"selector.selectPlaceholder": "Select a pack",

	"pack.onlineCount":       "%d player/s online",
	"pack.onlineCount.empty": "No players online",
	"pack.switching":         "Please wait, launching pack!",
	"pack.switch":            "Launch pack",
	"pack.switch.stop":       "Stop server",
	"pack.switch.occupied":   "Server is occupied",
	"pack.switchingFailed":   "Failed to launch pack! Not enough resources are available. I cannot automatically stop servers if someone's logged in, so try again later.",
	"pack.shutDown":          "%s was successfully shut down.",
	"pack.download":          "Get this pack",

	"process.kill":       "Shutting down %s",
	"process.kill.sub":   "This should only take a few seconds. Please wait...",
	"process.boot":       "Booting up %s",
	"process.boot.sub":   "Boot up the client -- the pack should be ready within a few minutes.",
	"process.booted":     "%s is ready!",
	"process.booted.sub": "Join the server and start playing!",

	"generic.back":   "Back",
	"generic.failed": "Sorry, something went wrong!",

	"color.generic": "#cba6f7",
	"color.error":   "#f38ba8",
	"color.success": "#a6e3a1",
}

func loc(key string) string {
	if value, ok := localizations[key]; ok {
		return value
	}
	return "<" + key + ">"
}
