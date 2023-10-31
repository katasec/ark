package handlers

import "embed"

var (
	//go:embed all:assets
	assets        embed.FS
	storagePrefix = "https://stkatasecassets.blob.core.windows.net/public"
)
