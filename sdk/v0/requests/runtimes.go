package requests

type RuntimeSpec struct {
	Go         string
	Dotnet     string
	Java       string
	TypeScript string
	JavaScript string
}

var Runtimes = &RuntimeSpec{
	Go:         "go",
	Dotnet:     "dotnet",
	Java:       "java",
	TypeScript: "nodejs",
	JavaScript: "nodejs",
}
