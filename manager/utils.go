package manager

import "runtime"

func GetHostsFile() string {
	switch runtime.GOOS {
	case "darwin":
		return "/etc/hosts"
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		return "/etc/hosts"
	default:
		return ""
	}
}

func GetLineSeparator() string {
	switch runtime.GOOS {
	case "darwin":
		return "\n"
	case "windows":
		return "\r\n"
	case "linux":
		return "\n"
	default:
		return "\n"
	}
}
