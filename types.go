package main

type configServer struct {
	DisplayName  string `yaml:"displayName" hc:"Server name displayed to users"`
	DisplayDesc  string `yaml:"displayDesc" hc:"Server description displayed to users, supports markdown"`
	DisplayIcon  string `yaml:"displayIcon" hc:"Server icon displayed to users (URL)"`
	DownloadLink string `yaml:"downloadLink" hc:"Client modpack download link, or leave empty if none"`

	ServerID      string `yaml:"serverId" hc:"Server ID, found below Server Details label in the web panel"`
	ResourceSlots int    `yaml:"resourceSlots" hc:"Number of resource slots this server takes up"`
	BootTime      int    `yaml:"bootTime" hc:"Boot time in seconds (estimated) - set to 0 if you got a mod that sends a 'Server is online' message or if you want to disable this feature"`
}
