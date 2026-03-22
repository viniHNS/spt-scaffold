<div align="center">

# SPT Scaffold

A TUI wizard that scaffolds ready-to-build SPT 4.x mod projects — server (C# / .NET 9) and client (BepInEx).

![Version](https://img.shields.io/github/v/release/viniHNS/spt-scaffold?label=version&color=orange&style=flat)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/platform-Windows%20%C2%B7%20Linux-lightgrey?style=flat)

[Server Mods](#server-mods) · [Client Mods](#client-mods) · [6 Templates](#templates) · [Windows · Linux](#build-from-source)

</div>

---

## Usage

Download the latest binary from [Releases](https://github.com/viniHNS/spt-scaffold-go/releases) and run it inside the folder where you want your mod created:

```sh
.\spt-scaffold.exe    # Windows
./spt-scaffold        # Linux
```

The wizard walks you through:

| Field | Description |
|---|---|
| Mod Type | `server` or `client` |
| Template | Project template (depends on mod type) |
| SPT Install Path | Absolute path to SPT installation (client mods only) |
| Mod Name | No spaces — used as C# namespace and folder name |
| Author | Your username |
| Version | Mod version in `x.y.z` format |
| SPT Version | Fetched live from NuGet — scrollable picker. Falls back to `4.0.0` if NuGet is unreachable. |
| Description | Optional, max 120 characters |
| Repository URL | Optional, must start with `https://` |
| License | SPDX license picker (MIT, Apache-2.0, GPL-3.0, and more) |

---

## Templates

### Server Mods

**C# / .NET 9 / NuGet**

| Template | Description | Files |
|---|---|---|
| `empty` | Minimal `IOnLoad` entry point | 5 |
| `editDatabase` | `DatabaseService` example with item injection | 5 |
| `readJsonConfig` | Custom `config.json` loaded via `ModHelper` | 6 |

> `readJsonConfig` generates an additional `config.json` file alongside the standard 5.

### Client Mods

**BepInEx / netstandard2.1**

| Template | Description | Files |
|---|---|---|
| `empty` | Minimal `BaseUnityPlugin` with `Awake` | 5 |
| `bepinexConfig` | Plugin with `Config.Bind` entries (F12 menu configurable) | 5 |
| `harmonyPatch` | Plugin with `ModulePatch` example (Prefix/Postfix) | 6 |

> `harmonyPatch` generates an additional `Patches/ExamplePatch.cs` file inside a subdirectory.

### Generated Output

**`server/empty`**
```
MyMod/
├── MyMod.sln
├── MyMod.csproj
├── Mod.cs
├── README.md
└── .gitignore
```

**`client/harmonyPatch`**
```
MyMod/
├── MyMod.sln
├── MyMod.csproj
├── Plugin.cs
├── README.md
├── .gitignore
└── Patches/
    └── ExamplePatch.cs
```

---

## Build from Source

**Requirements:** Go 1.22+

```sh
git clone https://github.com/viniHNS/spt-scaffold-go
cd spt-scaffold-go
go mod tidy
```

**Windows**
```sh
go build -ldflags="-s -w" -o spt-scaffold.exe
.\spt-scaffold.exe
```

**Linux**
```sh
go build -ldflags="-s -w" -o spt-scaffold
./spt-scaffold
```

> The prebuilt Windows binary does not display correctly under Wine. Build natively on Linux instead.

---

## Contributing a Template

Templates are plain files — no Go knowledge required. Adding a new template means creating a folder with `.tmpl` files and a `template.json`. The scaffold engine discovers and bundles templates automatically at compile time via Go's `embed` package.

### 1. Create the folder

Navigate to the correct directory and create a descriptive folder name:

```
internal/templates/server/<your-template>/    # for server mods
internal/templates/client/<your-template>/    # for client mods
```

The folder name becomes the template's identifier (e.g. `myTrader`).

### 2. Add `template.json`

Every template folder must contain a `template.json` that defines how it appears in the TUI picker:

```json
{
  "Label": "Custom Trader",
  "Desc":  "A starter kit for adding a new trader with custom items."
}
```

| Field | Type | Description |
|---|---|---|
| `Label` | string | Display name shown in the template picker |
| `Desc` | string | Short description shown below the label in the picker |

### 3. Add `.tmpl` files

Each `.tmpl` file in the folder becomes a generated file. The engine uses Go `text/template` syntax — variables are written as `{{.FieldName}}`.

**The filename itself can also use template variables:**

| Filename on disk | Generated as |
|---|---|
| `{{.ModName}}.csproj.tmpl` | `MyMod.csproj` |
| `{{.ModName}}.sln.tmpl` | `MyMod.sln` |
| `Mod.cs.tmpl` | `Mod.cs` |
| `Patches/ExamplePatch.cs.tmpl` | `Patches/ExamplePatch.cs` |

Every template must include `{{.ModName}}.csproj.tmpl` and `{{.ModName}}.sln.tmpl` — these are required to produce a buildable project.

Subdirectories are supported. The engine walks the template folder recursively, so `Patches/ExamplePatch.cs.tmpl` generates `Patches/ExamplePatch.cs` inside the mod directory.

#### Template Variables

All fields from `ModConfig` are available in every `.tmpl` file and in filenames:

| Variable | Example value | Description |
|---|---|---|
| `{{.ModName}}` | `MyAwesomeMod` | Mod identifier, C# namespace, folder name |
| `{{.Author}}` | `viniHNS` | Author username |
| `{{.Version}}` | `1.0.0` | Mod version |
| `{{.SptVersion}}` | `4.0.13` | Pinned SPT NuGet version |
| `{{.SptVersionRange}}` | `~4.0.0` | SemanticVersioning range for `Mod.cs` |
| `{{.Desc}}` | `A cool mod` | Short description |
| `{{.RepoURL}}` | `https://github.com/…` | Repository URL |
| `{{.License}}` | `MIT` | SPDX license identifier |
| `{{.ProjectGuid}}` | `{A1B2-…}` | Auto-generated UUID for `.sln` files |
| `{{.SptInstallPath}}` | `D:\Games\SPT` | SPT install path **(client mods only)** |

### 4. Test locally

Templates are bundled at compile time by Go's `embed` package. After adding your files, recompile the binary to include them:

**Windows**
```sh
go build -ldflags="-s -w" -o spt-scaffold.exe
.\spt-scaffold.exe
```

**Linux**
```sh
go build -ldflags="-s -w" -o spt-scaffold
./spt-scaffold
```

Your template will appear automatically in the picker. Verify the generated output looks correct before submitting.

### 5. Submit a Pull Request

Open a PR against `main`. Before submitting, confirm:

- [ ] `template.json` present with `Label` and `Desc`
- [ ] At least one `.tmpl` file
- [ ] `{{.ModName}}.csproj.tmpl` and `{{.ModName}}.sln.tmpl` present (required for a buildable project)
- [ ] Generated output builds with `dotnet build`
- [ ] No hardcoded author names or paths in template files

---

## Resources

| Resource | URL |
|---|---|
| SPT Server C# Overview | https://deepwiki.com/sp-tarkov/server-csharp/1-overview |
| Server Mod Examples | https://github.com/sp-tarkov/server-mod-examples |
| Client Mod Examples | https://github.com/Jehree/SPTClientModExamples |
| SPT Wiki — Modding Resources | https://wiki.sp-tarkov.com/modding/Modding_Resources |
| NuGet — SPTarkov.Server.Core | https://www.nuget.org/packages/SPTarkov.Server.Core |
