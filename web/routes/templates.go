package routes

type templateData struct {
	name string
	path string
}

func templates() []templateData {
	return []templateData{
		{"error.html", "web/views/error.html" },
		{"createwallet.html", "web/views/createwallet.html" },
		{"balance.html", "web/views/balance.html" },
		{"send.html", "web/views/send.html" },
		{"receive.html", "web/views/receive.html" },
		{"history.html", "web/views/history.html" },
	}
}