package process

func StateString(state string) string {
	switch state {
	case "R":
		return "Running"
	case "S":
		return "Sleeping"
	case "D":
		return "DiskSleep"
	case "Z":
		return "Zombie"
	case "T":
		return "Stopped"
	case "t":
		return "TracingStop"
	default:
		return "Unknwon"
	}
}
