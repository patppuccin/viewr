package constants

// Compile-time configurations ///////////////////
var AppVersion string

// App branding configurations ///////////////////
const AppFullName = "Viewr"
const AppAbbrName = "viewr"
const AppDescription = "Web-based file browser"
const AppBanner = `
█ █ █ ██▀ █   █ █▀▄ │ WEB-BASED, READ-ONLY
▀▄▀ █ █▄▄ ▀▄▀▄▀ █▀▄ │ FILE BROWSER ───── •
`

// Web Meta Configurations ///////////////////////

const SEOPageTitleSuffix = "File Browser"
const SEOPageDescription = "VIEWR is a web-based simple, local-first, read-only file browser."
const SEORobotsDirective = "index, follow"
const AssetLogoURL = "/assets/icons/logo.png"
const AssetFaviconURL = "/assets/icons/favicon.ico"
const AssetOGImageURL = "assets/images/og.png"

// CLI Configurations ////////////////////////////

var LogLevels = []string{"debug", "info", "warn", "error"}
