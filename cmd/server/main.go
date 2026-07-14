package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PageData struct {
	Title     string
	ActiveTab string
	IconSizes []string
	Users     []User
}

type User struct {
	ID   int
	Name string
}

func getDatabaseDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgresql://neondb_owner:npg_HucV8d0RzvKG@ep-divine-heart-asty9beu-pooler.c-4.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
}

func loadUsers(ctx context.Context, pool *pgxpool.Pool) ([]User, error) {
	rows, err := pool.Query(ctx, `SELECT id, "Name" FROM "User" ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

const layout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | Gemini Backoffice</title>
    
    <!-- Tailwind CSS with custom scale & pulse-glow styles -->
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        @keyframes cosmic-pulse {
            0%, 100% { transform: scale(1); filter: drop-shadow(0 0 15px rgba(99, 102, 241, 0.6)); }
            50% { transform: scale(1.08); filter: drop-shadow(0 0 25px rgba(168, 85, 247, 0.9)); }
        }
        .gemini-glow { animation: cosmic-pulse 4s infinite ease-in-out; }
    </style>
</head>
<!-- Global scaling: text-lg sets a larger structural baseline -->
<body class="bg-slate-50 text-slate-800 text-lg font-sans flex h-screen overflow-hidden antialiased">

    <!-- Sidebar Navigation -->
    <aside class="w-72 bg-slate-950 text-white flex flex-col justify-between border-r border-slate-900 shadow-2xl">
        <div class="p-6">
            <!-- Gemini Spark Header Effect -->
            <h1 class="text-2xl font-bold tracking-tight flex items-center gap-3 bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
                <span class="w-5 h-5 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 gemini-glow"></span>
                Gemini Admin
            </h1>
            
            <nav class="mt-12 space-y-3">
                <a href="/" class="flex items-center px-5 py-3.5 rounded-xl transition-all duration-300 font-medium {{if eq .ActiveTab "dashboard"}}bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-600/20{{else}}text-slate-400 hover:bg-slate-900 hover:text-slate-200{{end}}">
                    Dashboard
                </a>
                <a href="/users" class="flex items-center px-5 py-3.5 rounded-xl transition-all duration-300 font-medium {{if eq .ActiveTab "users"}}bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-600/20{{else}}text-slate-400 hover:bg-slate-900 hover:text-slate-200{{end}}">
                    User Management
                </a>
                <a href="/settings" class="flex items-center px-5 py-3.5 rounded-xl transition-all duration-300 font-medium {{if eq .ActiveTab "settings"}}bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-600/20{{else}}text-slate-400 hover:bg-slate-900 hover:text-slate-200{{end}}">
                    System Settings
                </a>
            </nav>
        </div>
        <div class="p-6 border-t border-slate-900 text-sm text-slate-500 tracking-wide">
            Platform Engine v1.2.0 • Go
        </div>
    </aside>

    <!-- Main Content Panel -->
    <div class="flex-1 flex flex-col overflow-y-auto">
        <!-- Expanded Header Header with Subtle Blur -->
        <header class="bg-white/80 backdrop-blur-md border-b border-slate-200 px-10 py-6 flex items-center justify-between sticky top-0 z-50">
            <h2 class="text-2xl font-bold text-slate-900 tracking-tight">{{.Title}}</h2>
            <div class="flex items-center gap-4">
                <span class="text-base text-slate-500 font-medium">Enterprise Core</span>
                <div class="w-11 h-11 rounded-xl bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center text-white font-extrabold text-base shadow-md">
                    GM
                </div>
            </div>
        </header>

        <!-- Dynamic View Content -->
        <main class="p-10 max-w-[1600px] w-full mx-auto space-y-8">
            {{template "content" .}}
        </main>
    </div>

</body>
</html>
`

const dashboardTemplate = `
{{define "content"}}
<div class="grid grid-cols-1 md:grid-cols-3 gap-8">
    <div class="bg-white p-8 rounded-2xl shadow-xl shadow-slate-100 border border-slate-200/60 transition-transform hover:-translate-y-1 duration-300">
        <p class="text-base text-slate-400 font-semibold tracking-wide uppercase">Core Node Users</p>
        <h3 class="text-4xl font-extrabold text-slate-900 mt-3 tracking-tight">{{len .Users}}</h3>
        <span class="text-emerald-500 text-sm font-bold bg-emerald-50 px-2.5 py-1 rounded-md inline-block mt-3">Live from PostgreSQL</span>
    </div>
    <div class="bg-white p-8 rounded-2xl shadow-xl shadow-slate-100 border border-slate-200/60 transition-transform hover:-translate-y-1 duration-300">
        <p class="text-base text-slate-400 font-semibold tracking-wide uppercase">Active Contexts</p>
        <h3 class="text-4xl font-extrabold text-slate-900 mt-3 tracking-tight">94,201</h3>
        <span class="text-emerald-500 text-sm font-bold bg-emerald-50 px-2.5 py-1 rounded-md inline-block mt-3">↑ 8.6% run-rate</span>
    </div>
    <div class="bg-white p-8 rounded-2xl shadow-xl shadow-indigo-500/5 border border-slate-200/60 transition-transform hover:-translate-y-1 duration-300">
        <p class="text-base text-slate-400 font-semibold tracking-wide uppercase">AI Cluster Health</p>
        <h3 class="text-4xl font-extrabold text-indigo-600 mt-3 tracking-tight">Nominal</h3>
        <span class="text-slate-400 text-sm font-medium inline-block mt-3">Operational • 99.99%</span>
    </div>
</div>

<div class="bg-white rounded-2xl shadow-xl shadow-slate-100 border border-slate-200/60 overflow-hidden">
    <div class="p-6 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
        <h3 class="text-xl font-bold text-slate-800">Recent Users</h3>
        <span class="text-sm font-medium text-slate-500">Loaded from the database</span>
    </div>
    <table class="w-full text-left border-collapse">
        <thead>
            <tr class="bg-slate-50/50 border-b border-slate-200">
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">ID</th>
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">Name</th>
            </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
            {{range .Users}}
            <tr class="hover:bg-slate-50/80 transition-colors">
                <td class="p-5 text-base text-slate-900 font-bold">{{.ID}}</td>
                <td class="p-5 text-base text-slate-600 font-medium">{{.Name}}</td>
            </tr>
            {{else}}
            <tr>
                <td colspan="2" class="p-5 text-base text-slate-600">No users found.</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</div>
{{end}}
`

const usersTemplate = `
{{define "content"}}
<div class="bg-white rounded-2xl shadow-xl shadow-slate-100 border border-slate-200/60 overflow-hidden">
    <div class="p-6 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
        <h3 class="text-xl font-bold text-slate-800">Identity Directory</h3>
        <button class="px-5 py-2.5 bg-slate-900 hover:bg-slate-800 text-white font-semibold rounded-xl text-sm transition-all">Add Directory Entity</button>
    </div>
    <table class="w-full text-left border-collapse">
        <thead>
            <tr class="bg-slate-50/50 border-b border-slate-200">
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">Identity</th>
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">Secure Anchor Routing Address</th>
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">System Role</th>
                <th class="p-5 font-bold text-slate-500 text-sm uppercase tracking-wider">Status Cluster</th>
            </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
            <tr class="hover:bg-slate-50/80 transition-colors">
                <td class="p-5 text-base text-slate-900 font-bold">Sarah Connor</td>
                <td class="p-5 text-base text-slate-600 font-mono">sconnor@cyberdyne.io</td>
                <td class="p-5 text-base text-slate-600 font-medium">Cluster Director</td>
                <td class="p-5"><span class="bg-emerald-100 text-emerald-800 text-xs px-3 py-1.5 rounded-lg font-bold">Active</span></td>
            </tr>
            <tr class="hover:bg-slate-50/80 transition-colors">
                <td class="p-5 text-base text-slate-900 font-bold">John Doe</td>
                <td class="p-5 text-base text-slate-600 font-mono">johndoe@gemini.net</td>
                <td class="p-5 text-base text-slate-600 font-medium">Model Specialist</td>
                <td class="p-5"><span class="bg-emerald-100 text-emerald-800 text-xs px-3 py-1.5 rounded-lg font-bold">Active</span></td>
            </tr>
            <tr class="hover:bg-slate-50/80 transition-colors">
                <td class="p-5 text-base text-slate-900 font-bold">Miles Dyson</td>
                <td class="p-5 text-base text-slate-600 font-mono">mdyson@cyberdyne.io</td>
                <td class="p-5 text-base text-slate-600 font-medium">System Engineer</td>
                <td class="p-5"><span class="bg-rose-100 text-rose-800 text-xs px-3 py-1.5 rounded-lg font-bold">Revoked</span></td>
            </tr>
        </tbody>
    </table>
</div>
{{end}}
`

const settingsTemplate = `
{{define "content"}}
<div class="bg-white p-8 rounded-2xl shadow-xl shadow-slate-100 border border-slate-200/60 max-w-3xl">
    <h4 class="text-xl font-bold mb-6 text-slate-900">Application Configuration</h4>
    
    <div class="space-y-6">
        <div>
            <label class="block text-base font-bold text-slate-700 tracking-wide">Platform Title Identity</label>
            <input type="text" class="mt-2 block w-full rounded-xl border border-slate-300 p-3.5 shadow-sm text-lg focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all" value="Gemini Backoffice Platform">
        </div>
        
        <div>
            <label class="block text-base font-bold text-slate-700 tracking-wide mb-2">Required Asset Icons Manifest Mapping</label>
            <div class="bg-slate-50 rounded-xl p-5 border border-slate-200">
                <ul class="space-y-2 text-base text-slate-600 font-medium list-disc list-inside">
                    {{range .IconSizes}}
                    <li>{{.}}</li>
                    {{end}}
                </ul>
            </div>
        </div>
        
        <div class="pt-4">
            <button class="px-6 py-3.5 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white text-base font-bold rounded-xl shadow-lg shadow-indigo-600/20 transition-all transform active:scale-[0.98]">
                Deploy Parameters
            </button>
        </div>
    </div>
</div>
{{end}}
`

func main() {
	sizes := []string{
		"16x16 — Standard Browser Tab Node",
		"32x32 — High-Density Monitor Icon Matrix",
		"180x180 — Apple Mobile Hardware Web Clip Container",
		"192x192 — Chromium Progressive Web Application Vector",
		"512x512 — High Fidelity PWA Splash Screen Vector",
	}

	dsn := getDatabaseDSN()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create PostgreSQL pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		users, err := loadUsers(r.Context(), pool)
		if err != nil {
			log.Printf("failed to load users for dashboard: %v", err)
			users = nil
		}

		tmpl, _ := template.New("layout").Parse(layout + dashboardTemplate)
		tmpl.Execute(w, PageData{Title: "Dashboard Space", ActiveTab: "dashboard", IconSizes: sizes, Users: users})
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("layout").Parse(layout + usersTemplate)
		tmpl.Execute(w, PageData{Title: "Identity Directory", ActiveTab: "users", IconSizes: sizes})
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("layout").Parse(layout + settingsTemplate)
		tmpl.Execute(w, PageData{Title: "Control Configurations", ActiveTab: "settings", IconSizes: sizes})
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthCtx, healthCancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer healthCancel()

		var now time.Time
		if err := pool.QueryRow(healthCtx, "select now()").Scan(&now); err != nil {
			http.Error(w, fmt.Sprintf("database query failed: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "PostgreSQL OK at %s", now.UTC().Format(time.RFC3339Nano))
	})

	log.Println("Server running at http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
