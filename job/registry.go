package job

var JobRegistry = map[string]func(){
	"notify-hot-leads": NotifyHotLeads,
}
