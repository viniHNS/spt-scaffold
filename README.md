# SPT Scaffold

A TUI wizard that scaffolds ready-to-build **SPT 4.0 server mod** projects.

![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

![terminal](https://media0.giphy.com/media/v1.Y2lkPTc5MGI3NjExZjVqbm9ueTVzeTFob3Y0MmxhdGh4Y3BvZm8ybHZ0bDBpZnR6aDNtMiZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/1ZUBF93ibXGuPnLT6u/giphy.gif)

---

## Usage

Download the latest binary from [Releases](https://github.com/viniHNS/spt-scaffold-go/releases) and run it inside the folder where you want your mod created.

```sh
./spt-scaffold
```

The wizard will ask for:

| Field | Description |
|---|---|
| Mod Name | No spaces, used as the C# namespace |
| Author | Your username |
| Version | Mod version (semver) |
| SPT Compatibility | All 4.x or a specific version |
| Description | Short description (optional) |
| Repository URL | `https://` URL (optional) |
| License | MIT, Apache-2.0, GPL-3.0, and more |

Generated files inside `<cwd>/<ModName>/`:

```
MyMod/
├── MyMod.csproj
├── Mod.cs
├── README.md
└── .gitignore
```

---

## Build from source

```sh
git clone https://github.com/viniHNS/spt-scaffold-go
cd spt-scaffold-go
go mod tidy
go build -ldflags="-s -w" -o spt-scaffold.exe
```

---

## Resources

- [SPT Server C# Overview](https://deepwiki.com/sp-tarkov/server-csharp/1-overview)
- [Server Mod Examples](https://github.com/sp-tarkov/server-mod-examples)
- [SPT Wiki](https://wiki.sp-tarkov.com/modding/Modding_Resources)
