# spt-scaffold

A TUI wizard that scaffolds ready-to-build **SPT 4.0 mod** projects — both **server** and **client** mods.

![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

---

## Usage

Download the latest binary from [Releases](https://github.com/viniHNS/spt-scaffold-go/releases) and run it inside the folder where you want your mod created.

```sh
./spt-scaffold
```

The wizard will ask for:

| Field | Description |
|---|---|
| Mod Type | `server` or `client` |
| Template | Project template (depends on mod type) |
| SPT Install Path | Absolute path to SPT installation (client mods only) |
| Mod Name | No spaces, used as the C# namespace |
| Author | Your username |
| Version | Mod version (semver) |
| SPT Version | Fetched live from NuGet (scrollable picker) |
| Description | Short description, max 120 chars (optional) |
| Repository URL | `https://` URL (optional) |
| License | MIT, Apache-2.0, GPL-3.0, and more |

### Available Templates

**Server:**
- **Empty** — Minimal `IOnLoad` entry point
- **Edit Database** — `DatabaseService` example
- **Read Json Config** — Custom `config.json` via `ModHelper`

**Client:**
- **Empty** — Minimal BepInEx plugin with `Awake`
- **BepInEx Configuration** — Plugin with `Config.Bind` entries (F12 configurable)
- **Harmony Patch** — Plugin with example `ModulePatch` (Prefix/Postfix)

### Generated Output

**Server mod** (`server/empty`):
```
MyMod/
├── MyMod.sln
├── MyMod.csproj
├── Mod.cs
├── README.md
└── .gitignore
```

**Client mod** (`client/empty`):
```
MyMod/
├── MyMod.sln
├── MyMod.csproj
├── Plugin.cs
├── README.md
└── .gitignore
```

> Some templates may generate additional files (e.g., `readJsonConfig` also creates `config.json`).

---

## Build from source

**Windows**
```sh
git clone https://github.com/viniHNS/spt-scaffold-go
cd spt-scaffold-go
go mod tidy
go build -ldflags="-s -w" -o spt-scaffold.exe
.\spt-scaffold.exe
```

**Linux**
```sh
git clone https://github.com/viniHNS/spt-scaffold-go
cd spt-scaffold-go
go mod tidy
go build -ldflags="-s -w" -o spt-scaffold
./spt-scaffold
```

> **Note:** The prebuilt Windows binary does not display correctly under Wine.
> Build natively for Linux instead.

---

---

## How to Contribute a New Template

`spt-scaffold-go` is designed to be extensible. Since it uses Go's `embed` package, adding a new template (Server or Client) requires **zero changes** to the Go source code!

### Quick Step-by-Step

> Use existing folders in `internal/templates/` as a base for your new template.

1.  **Create your folder**: Navigate to `internal/templates/server/` or `internal/templates/client/` and create a descriptive folder name (e.g., `myTraderTemplate`).
2.  **Add Metadata**: Create a `template.json` file to define how it appears in the TUI:
    ```json
    {
      "Label": "Custom Trader",
      "Desc": "A starter kit for adding a new trader with custom items."
    }
    ```
3.  **Create Template Files**: Add your `.cs`, `.csproj`, or other files using the `.tmpl` extension.

### Template Syntax & Variables

The engine uses standard Go `text/template`. You can use these variables anywhere in the **file content** or even in the **filename**:

| Variable | Description |
|---|---|
| `{{.ModName}}` | Mod identifier (e.g., `MyAwesomeMod`) |
| `{{.Author}}` | Your username |
| `{{.Version}}` | Mod version (e.g., `1.0.0`) |
| `{{.SptVersion}}` | Pinned SPT version (e.g., `4.0.13`) |
| `{{.SptVersionRange}}` | Version range (e.g., `~4.0.0`) |
| `{{.Desc}}` | Short description |
| `{{.RepoURL}}` | Repository URL |
| `{{.License}}` | Selected license (SPDX identifier) |
| `{{.ProjectGuid}}` | Auto-generated UUID for Visual Studio |
| `{{.SptInstallPath}}` | **Client mods only** — Path to SPT install |

> **Dynamic Filenames:** To name a file based on user input, use the variable in the filename!  
> **(Incorrect)** `Mod.csproj.tmpl`  
> **(Correct)** `{{.ModName}}.csproj.tmpl`

### Test & Submit

1.  **Rebuild**: Run `go build -ldflags="-s -w" -o spt-scaffold.exe` to bundle the new files.
2.  **Run**: Launch the tool. Your new template will appear automatically in the list.
3.  **PR**: Everything working? Submit a Pull Request! We love community contributions.

---

---

## Resources

- [SPT Server C# Overview](https://deepwiki.com/sp-tarkov/server-csharp/1-overview)
- [Server Mod Examples](https://github.com/sp-tarkov/server-mod-examples)
- [SPT Client Mod Examples](https://github.com/Jehree/SPTClientModExamples)
- [SPT Wiki](https://wiki.sp-tarkov.com/modding/Modding_Resources)
