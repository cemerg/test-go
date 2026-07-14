package main

import (
	"html/template"
	"log"
	"net/http"
)

// PageData represents the context passed to the HTML templates
type PageData struct {
	Title       string
	ActiveTab   string
	IconSizes   []string
}

// Global layout containing the Sidebar and Header (configured with Gemini-standard icon definitions)
const layout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | Backoffice</title>
    
    <link rel="icon" type="image/png" sizes="16x16" href="/static/favicon-16x16.png">
    <link rel="icon" type="image/png" sizes="32x32" href="/static/favicon-32x32.png">
    <link rel="apple-touch-icon" sizes="180x180" href="/static/apple-touch-icon.png">
    <link rel="manifest" href="/static/site.webmanifest">
    
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 font-sans flex h-screen overflow-hidden">

    <aside class="w-64 bg-slate-900 text-white flex flex-col justify-between">
        <div class="p-5">
            <h1 class="text-xl font-bold tracking-wider flex items-center gap-2">
                <span class="w-4 h-4 rounded-full bg-indigo-500 animate-pulse"></span>
                GEMINI ADMIN
            </h1>
            <nav class="mt-8 space-y-2">
                <a href="/" class="flex items-center px-4 py-2.5 rounded transition-all {{if eq .ActiveTab "dashboard"}}bg-indigo-600 text-white{{else}}text-gray-400 hover:bg-slate-800{{end}}">
                    Dashboard
                </a>
                <a href="/users" class="flex items-center px-4 py-2.5 rounded transition-all {{if eq .ActiveTab "users"}}bg-indigo-600 text-white{{else}}text-gray-400 hover:bg-slate-800{{end}}">
                    Users
                </a>
                <a href="/settings" class="flex items-center px-4 py-2.5 rounded transition-all {{if eq .ActiveTab "settings"}}bg-indigo-600 text-white{{else}}text-gray-400 hover:bg-slate-800{{end}}">
                    Settings
                </a>
            </nav>
        </div>
        <div class="p-5 border-t border-slate-800 text-xs text-gray-500">
            v1.0.0 • Go 1.22+
        </div>
    </aside>

    <div class="flex-1 flex flex-col overflow-y-auto">
        <header class="bg-white border-b px-8 py-4 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-gray-800">{{.Title}}</h2>
            <div class="flex items-center gap-3">
                <span class="text-sm text-gray-600">Admin Account</span>
                <div class="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-white font-bold text-sm">A</div>
            </div>
        </header>

        <main class="p-8">
            {{template "content" .}}
        </main>
    </div>

</body>
</html>
`

const dashboardTemplate = `
{{define "content"}}
<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
    <div class="bg-white p-6 rounded-xl shadow-sm border">
        <p class="text-sm text-gray-500 font-medium">Total Users</p>
        <h3 class="text-3xl font-bold text-gray-900 mt-2">1,248</h3>
        <span class="text-green-500 text-xs font-semibold">↑ 12% this week</span>
    </div>
    <div class="bg-white p-6 rounded-xl shadow-sm border">
        <p class="text-sm text-gray-500 font-medium">Monthly Active</p>
        <h3 class="text-3xl font-bold text-gray-900 mt-2">842</h3>
        <span class="text-green-500 text-xs font-semibold">↑ 4.3% this week</span>
    </div>
    <div class="bg-white p-6 rounded-xl shadow-sm border">
        <p class="text-sm text-gray-500 font-medium">Server Status</p>
        <h3 class="text-3xl font-bold text-green-600 mt-2">Online</h3>
        <span class="text-gray-400 text-xs">Uptime: 99.98%</span>
    </div>
</div>

<div class="bg-white p-6 rounded-xl shadow-sm border">
    <h4 class="text-lg font-bold mb-4 text-gray-800">System Information</h4>
    <p class="text-gray-600">This Go micro-backoffice uses custom Web Icon templates scaling up to 512x512 pixels to mimic native-app fidelity.</p>
</div>
{{end}}
`

const usersTemplate = `
{{define "content"}}
<div class="bg-white rounded-xl shadow-sm border overflow-hidden">
    <table class="w-full text-left border-collapse">
        <thead>
            <tr class="bg-gray-50 border-b">
                <th class="p-4 font-semibold text-gray-600 text-sm">Name</th>
                <th class="p-4 font-semibold text-gray-600 text-sm">Email</th>
                <th class="p-4 font-semibold text-gray-600 text-sm">Role</th>
                <th class="p-4 font-semibold text-gray-600 text-sm">Status</th>
            </tr>
        </thead>
        <tbody>
            <tr class="border-b">
                <td class="p-4 text-sm text-gray-800 font-medium">Sarah Connor</td>
                <td class="p-4 text-sm text-gray-600">sconnor@cyberdyne.com</td>
                <td class="p-4 text-sm text-gray-600">Administrator</td>
                <td class="p-4 text-sm"><span class="bg-green-100 text-green-800 text-xs px-2.5 py-1 rounded-full font-semibold">Active</span></td>
            </tr>
            <tr class="border-b">
                <td class="p-4 text-sm text-gray-800 font-medium">John Doe</td>
                <td class="p-4 text-sm text-gray-600">johndoe@example.com</td>
                <td class="p-4 text-sm text-gray-600">Editor</td>
                <td class="p-4 text-sm"><span class="bg-green-100 text-green-800 text-xs px-2.5 py-1 rounded-full font-semibold">Active</span></td>
            </tr>
            <tr>
                <td class="p-4 text-sm text-gray-800 font-medium">Miles Dyson</td>
                <td class="p-4 text-sm text-gray-600">mdyson@cyberdyne.com</td>
                <td class="p-4 text-sm text-gray-600">Viewer</td>
                <td class="p-4 text-sm"><span class="bg-red-100 text-red-800 text-xs px-2.5 py-1 rounded-full font-semibold">Inactive</span></td>
            </tr>
        </tbody>
    </table>
</div>
{{end}}
`

const settingsTemplate = `
{{define "content"}}
<div class="bg-white p-6 rounded-xl shadow-sm border max-w-2xl">
    <h4 class="text-lg font-bold mb-6 text-gray-800">Application Icon Settings</h4>
    <div class="space-y-4">
        <div>
            <label class="block text-sm font-medium text-gray-700">App Name</label>
            <input type="text" class="mt-1 block w-full rounded-md border border-gray-300 p-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm" value="Gemini Backoffice">
        </div>
        <div>
            <label class="block text-sm font-medium text-gray-700">PWA Manifest (Required Icon Sizes)</label>
            <ul class="mt-2 space-y-1 text-sm text-gray-500 list-disc list-inside">
                {{range .IconSizes}}
                <li>{{.}}</li>
                {{end}}
            </ul>
        </div>
        <button class="mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-semibold rounded-md shadow-sm">Save Changes</button>
    </div>
</div>
{{end}}
`

func main() {
	// Root base sizes
	sizes := []string{"16x16 (Browser Tab)", "32x32 (Browser Tab HD)", "180x180 (Apple Touch Icon)", "192x192 (PWA Standard)", "512x512 (Splash Screen)"}

	// Router setup
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl, _ := template.New("layout").Parse(layout + dashboardTemplate)
		tmpl.Execute(w, PageData{Title: "Dashboard Overview", ActiveTab: "dashboard", IconSizes: sizes})
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("layout").Parse(layout + usersTemplate)
		tmpl.Execute(w, PageData{Title: "Users Management", ActiveTab: "users", IconSizes: sizes})
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("layout").Parse(layout + settingsTemplate)
		tmpl.Execute(w, PageData{Title: "System Settings", ActiveTab: "settings", IconSizes: sizes})
	})

	log.Println("Server starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}