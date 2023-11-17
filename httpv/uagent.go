package httpv

import (
	"fmt"
	"math/rand"
	"strings"
)

var uaGens = []func() string{
	GenFirefoxUA,
	GenChromeUA,
	GenEdgeUA,
}

var uaGensMobile = []func() string{
	GenMobilePixel7UA,
}

func RandomUserAgent() string {
	return uaGens[rand.Intn(len(uaGens))]()
}

func RandomMobileUserAgent() string {
	return uaGensMobile[rand.Intn(len(uaGensMobile))]()
}

var ffVersions = []float32{
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://en.wikipedia.org/wiki/Firefox_version_history

	// 2022
	102.0,
	103.0,
	104.0,
	105.0,
	106.0,
	107.0,
	108.0,

	// 2023
	109.0,
	110.0,
	111.0,
	112.0,
	113.0,
	114.0,
	115.0,
	116.0,
	117.0,
	118.0,
	119.0,
	120.0,
}

var chromeVersions = []string{
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://chromereleases.googleblog.com/search/label/Stable%20updates

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop.html
	"102.0.5005.115",

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop_21.html
	"103.0.5060.53",

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop_27.html
	"103.0.5060.66",

	// https://chromereleases.googleblog.com/2022/07/stable-channel-update-for-desktop.html
	"103.0.5060.114",

	// https://chromereleases.googleblog.com/2022/07/stable-channel-update-for-desktop_19.html
	"103.0.5060.134",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop.html
	"104.0.5112.79",
	"104.0.5112.80",
	"104.0.5112.81",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop_16.html
	"104.0.5112.101",
	"104.0.5112.102",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop_30.html
	"105.0.5195.52",
	"105.0.5195.53",
	"105.0.5195.54",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop.html
	"105.0.5195.102",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_14.html
	"105.0.5195.125",
	"105.0.5195.126",
	"105.0.5195.127",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_27.html
	"106.0.5249.61",
	"106.0.5249.62",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_30.html
	"106.0.5249.91",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop.html
	"106.0.5249.103",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_11.html
	"106.0.5249.119",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_25.html
	"107.0.5304.62",
	"107.0.5304.63",
	"107.0.5304.68",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_27.html
	"107.0.5304.87",
	"107.0.5304.88",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop.html
	"107.0.5304.106",
	"107.0.5304.107",
	"107.0.5304.110",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop_24.html
	"107.0.5304.121",
	"107.0.5304.122",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop_29.html
	"108.0.5359.71",
	"108.0.5359.72",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop.html
	"108.0.5359.94",
	"108.0.5359.95",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop_7.html
	"108.0.5359.98",
	"108.0.5359.99",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop_13.html
	"108.0.5359.124",
	"108.0.5359.125",

	// https://chromereleases.googleblog.com/2023/01/stable-channel-update-for-desktop.html
	"109.0.5414.74",
	"109.0.5414.75",
	"109.0.5414.87",

	// https://chromereleases.googleblog.com/2023/01/stable-channel-update-for-desktop_24.html
	"109.0.5414.119",
	"109.0.5414.120",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-update-for-desktop.html
	"110.0.5481.77",
	"110.0.5481.78",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update.html
	"110.0.5481.96",
	"110.0.5481.97",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_14.html
	"110.0.5481.100",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_16.html
	"110.0.5481.104",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_22.html
	"110.0.5481.177",
	"110.0.5481.178",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_97.html
	"109.0.5414.129",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop.html
	"111.0.5563.64",
	"111.0.5563.65",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop_21.html
	"111.0.5563.110",
	"111.0.5563.111",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop_27.html
	"111.0.5563.146",
	"111.0.5563.147",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop.html
	"112.0.5615.49",
	"112.0.5615.50",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_12.html
	"112.0.5615.86",
	"112.0.5615.87",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_14.html
	"112.0.5615.121",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_18.html
	"112.0.5615.137",
	"112.0.5615.138",
	"112.0.5615.165",

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop.html
	"113.0.5672.63",
	"113.0.5672.64",

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop_8.html
	"113.0.5672.92",
	"113.0.5672.93",
}

var edgeVersions = []string{
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://learn.microsoft.com/en-us/deployedge/microsoft-edge-release-schedule

	// 2022
	"103.0.0.0,103.0.1264.37",
	"104.0.0.0,104.0.1293.47",
	"105.0.0.0,105.0.1343.25",
	"106.0.0.0,106.0.1370.34",
	"107.0.0.0,107.0.1418.24",
	"108.0.0.0,108.0.1462.42",

	// 2023
	"109.0.0.0,109.0.1518.49",
	"110.0.0.0,110.0.1587.41",
	"111.0.0.0,111.0.1661.41",
	"112.0.0.0,112.0.1722.34",
	"113.0.0.0,113.0.1774.3",
}

var pixel7AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"13",
}

var osStrings = []string{

	// Windows
	"Windows NT 10.0; Win64; x64",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",

	// Linux
	"X11; Linux x86_64",
}

// Generates Firefox Browser User-Agent (Desktop)
//
// -> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0"
func GenFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

// Generates Chrome Browser User-Agent (Desktop)
//
// -> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
func GenChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Generates Microsoft Edge User-Agent (Desktop)
//
// -> "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.39"
func GenEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	edgeVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s", os, chromeVersion, edgeVersion)
}

// Generates Pixel 7 Browser User-Agent (Mobile)
//
// -> Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36
func GenMobilePixel7UA() string {
	android := pixel7AndroidVersions[rand.Intn(len(pixel7AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, chrome)
}

// // RandomUserAgent generates a random DESKTOP browser user-agent on every requests
// func RandomUserAgent(c *colly.Collector) {
// 	c.OnRequest(func(r *colly.Request) {
// 		r.Headers.Set("User-Agent", uaGens[rand.Intn(len(uaGens))]())
// 	})
// }

// // RandomMobileUserAgent generates a random MOBILE browser user-agent on every requests
// func RandomMobileUserAgent(c *colly.Collector) {
// 	c.OnRequest(func(r *colly.Request) {
// 		r.Headers.Set("User-Agent", uaGensMobile[rand.Intn(len(uaGensMobile))]())
// 	})
// }
