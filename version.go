package suprsend

import (
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/suprsend/suprsend-go"

// Version returns the SDK version reported in the User-Agent.
// A released build uses the module version selected by the consumer, which is the git tag without a leading v.
// A local build that is not a released module version returns "dev".
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	mod, ok := moduleVersion(info)
	if !ok || mod.Version == "" || mod.Version == "(devel)" {
		return "dev"
	}
	return strings.TrimPrefix(mod.Version, "v")
}

func moduleVersion(info *debug.BuildInfo) (debug.Module, bool) {
	if info.Main.Path == modulePath {
		return replaced(info.Main), true
	}
	for _, dep := range info.Deps {
		if dep != nil && dep.Path == modulePath {
			return replaced(*dep), true
		}
	}
	return debug.Module{}, false
}

func replaced(mod debug.Module) debug.Module {
	if mod.Replace != nil {
		return *mod.Replace
	}
	return mod
}
