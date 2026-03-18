package scaffold

// CsprojTemplate is the .csproj XML template.
const CsprojTemplate = `<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net9.0</TargetFramework>
    <OutputType>Library</OutputType>
    <ImplicitUsings>enable</ImplicitUsings>
    <AppendTargetFrameworkToOutputPath>false</AppendTargetFrameworkToOutputPath>
    <RootNamespace>{{.ModName}}</RootNamespace>
    <Version>{{.Version}}</Version>
    <Authors>{{.Author}}</Authors>
    <Company>{{.Author}}</Company>
    <Product>{{.ModName}}</Product>
    <Description>{{.Desc}}</Description>
    <Copyright>Copyright © {{.Author}}</Copyright>
    <RepositoryUrl>{{.RepoURL}}</RepositoryUrl>
    <PackageLicenseExpression>{{.License}}</PackageLicenseExpression>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="SPTarkov.Common" Version="{{.SptVersion}}" />
    <PackageReference Include="SPTarkov.DI" Version="{{.SptVersion}}" />
    <PackageReference Include="SPTarkov.Server.Core" Version="{{.SptVersion}}" />
  </ItemGroup>

  <ItemGroup>
    <_ContentIncludedByDefault Remove="dist\package.json" />
  </ItemGroup>

  <Target Name="PackageModForDistribution" AfterTargets="Build">
    <PropertyGroup>
      <ModName>$(AssemblyName)</ModName>
      <ZipFileName>$(ModName).zip</ZipFileName>
      <TempPackageDir>$(ProjectDir)obj\Package</TempPackageDir>
      <ModPackageDir>$(TempPackageDir)\SPT\user\mods\$(ModName)</ModPackageDir>
    </PropertyGroup>
    <Message Text="--- Starting Mod packaging for distribution ---" Importance="high" />
    <RemoveDir Directories="$(TempPackageDir)" />
    <MakeDir Directories="$(ModPackageDir)" />
    <Copy SourceFiles="$(TargetPath)" DestinationFolder="$(ModPackageDir)" />
    <ZipDirectory SourceDirectory="$(TempPackageDir)" DestinationFile="$(ProjectDir)$(ZipFileName)" Overwrite="true" />
    <RemoveDir Directories="$(TempPackageDir)" />
    <Message Text="--- Mod successfully packaged at: $(ProjectDir)$(ZipFileName) ---" Importance="high" />
  </Target>
</Project>
`

// ModCSTemplate is the main Mod.cs entry-point template.
const ModCSTemplate = `using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.DI;
using SPTarkov.Server.Core.Models.Enums;
using SPTarkov.Server.Core.Models.Spt.Mod;
using SPTarkov.Server.Core.Models.Utils;

namespace {{.ModName}};

public record ModMetadata : AbstractModMetadata
{
    /// <summary>
    /// Any string can be used for a ModGuid, but it should ideally be unique and not easily duplicated.
    /// A 'bad' ID would be: "mymod", "mod1", "questmod"
    /// It is recommended (but not mandatory) to use the reverse domain name notation,
    /// see: https://docs.oracle.com/javase/tutorial/java/package/namingpkgs.html
    /// </summary>
    public override string ModGuid { get; init; } = "com.{{.Author}}.{{.ModName}}";
    public override string Name { get; init; } = "{{.ModName}}";
    public override string Author { get; init; } = "{{.Author}}";
    public override List<string>? Contributors { get; init; }
    public override SemanticVersioning.Version Version { get; init; } = new("{{.Version}}");
    public override SemanticVersioning.Range SptVersion { get; init; } = new("{{.SptVersionRange}}");
    public override List<string>? Incompatibilities { get; init; }
    public override Dictionary<string, SemanticVersioning.Range>? ModDependencies { get; init; }
    public override string? Url { get; init; } = "{{.RepoURL}}";
    public override bool? IsBundleMod { get; init; }
    public override string License { get; init; } = "{{.License}}";
}

[Injectable(TypePriority = OnLoadOrder.PostDBModLoader + 1)]
public class {{.ModName}}Mod(ISptLogger<{{.ModName}}Mod> logger) : IOnLoad
{
    public Task OnLoad()
    {
        logger.Info("{{.ModName}} loaded successfully!");
        // TODO: Add your mod initialization logic here
        return Task.CompletedTask;
    }
}
`

// ReadmeTemplate is the README.md template.
const ReadmeTemplate = `# {{.ModName}}

> {{.Desc}}

**Author:** {{.Author}}
**Version:** {{.Version}}
**SPT Version:** {{.SptVersion}}
**License:** {{.License}}

---

## What This Mod Does

Describe what your mod does here. Be specific about game systems it modifies.

---

## Requirements

- [SPT](https://www.sp-tarkov.com/) **{{.SptVersion}}** or compatible
- .NET 9 SDK (for building from source)

---

## Building

` + "```" + `sh
git clone {{.RepoURL}}
cd {{.ModName}}
dotnet build -c Release
` + "```" + `

The build target automatically packages the mod into a distributable ` + "`{{.ModName}}.zip`" + `.

---

## Installation

1. Build the project (see above) **or** download the latest release zip.
2. Extract the zip so that ` + "`{{.ModName}}.dll`" + ` ends up in:
   ` + "```" + `
   <SPT root>/user/mods/{{.ModName}}/
   ` + "```" + `
3. Launch SPT server as usual.

---

## Configuration

No configuration file is required by default. Extend ` + "`Mod.cs`" + ` to add your own settings.

---

## Project Structure

` + "```" + `
{{.ModName}}/
├── {{.ModName}}.csproj   ← project file + packaging target
├── Mod.cs                ← entry point (ModMetadata + IOnLoad)
├── README.md
└── .gitignore
` + "```" + `

---

## Learning Resources

| Resource | URL |
|---|---|
| SPT Server (C#) — Overview | https://deepwiki.com/sp-tarkov/server-csharp/1-overview |
| Server Mod Examples | https://github.com/sp-tarkov/server-mod-examples |
| SPT Wiki Modding Resources | https://wiki.sp-tarkov.com/modding/Modding_Resources |
| SPT Client Mod Examples | https://github.com/Jehree/SPTClientModExamples |

---

`

// GitignoreTemplate is the .gitignore template.
const GitignoreTemplate = `## .NET / C# standard gitignore

# Build results
[Dd]ebug/
[Dd]ebugPublic/
[Rr]elease/
[Rr]eleases/
x64/
x86/
[Ww][Ii][Nn]32/
[Aa][Rr][Mm]/
[Aa][Rr][Mm]64/
bld/
[Bb]in/
[Oo]bj/
[Ll]og/
[Ll]ogs/

# MSTest test Results
[Tt]est[Rr]esult*/
[Bb]uild[Ll]og.*

# NuGet
*.nupkg
*.snupkg
.nuget/
packages/
!**/packages/build/
*.nuspec
project.lock.json
project.fragment.lock.json
artifacts/

# User-specific files
*.rsuser
*.suo
*.user
*.userosscache
*.sln.docstates

# Visual Studio
.vs/
*.vsp
*.vsps
*.vspsx

# JetBrains Rider
.idea/
*.sln.iml

# mono auto created files
mono_crash.*

# Windows image file caches
Thumbs.db
ehthumbs.db

# Folder config file
Desktop.ini

# Mac OS
.DS_Store

# SPT mod packaging output
*.zip
obj/Package/
`
