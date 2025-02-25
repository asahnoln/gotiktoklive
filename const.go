package gotiktoklive

import (
	"errors"
	"regexp"
)

const (
	// Base URL
	tiktokBaseUrl = "https://www.tiktok.com/"
	tiktokAPIUrl  = "https://webcast.tiktok.com/webcast/"
	// TODO: Use fallback signer as well
	// signProviderHost: 'https://tiktok.eulerstream.com/',
	// signProviderFallbackHosts: ['https://tiktok-sign.zerody.one/'],
	tiktokSigner = "https://tiktok.eulerstream.com/"
	// Seems not to work?
	// tiktokSigner = "https://tiktok-sign.zerody.one/"

	// Endpoints
	urlLive      = "live/"
	urlFeed      = "feed/"
	urlRankList  = "ranklist/online_audience/"
	urlPriceList = "diamond/"
	urlUser      = "@%s/"
	// Think this changed to room/enter/
	urlRoomInfo = "room/info/"
	urlRoomData = "im/fetch/"
	urlGiftInfo = "gift/list/"
	urlSignReq  = "webcast/sign_url/"

	// Default Request Headers
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	referer   = "https://www.tiktok.com/"
	origin    = "https://www.tiktok.com"
	clientId  = "ttlive-golang"
)

var (
	// Default GET Request parameters
	defaultGETParams = map[string]string{
		"aid":                 "1988",
		"app_language":        "en-US",
		"app_name":            "tiktok_web",
		"browser_language":    "en",
		"browser_name":        "Mozilla",
		"browser_online":      "true",
		"browser_platform":    "Win32",
		"browser_version":     "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"cookie_enabled":      "true",
		"cursor":              "",
		"internal_ext":        "",
		"device_platform":     "web",
		"focus_state":         "true",
		"from_page":           "user",
		"history_len":         "0",
		"is_fullscreen":       "false",
		"is_page_visible":     "true",
		"did_rule":            "3",
		"fetch_rule":          "1",
		"last_rtt":            "0",
		"live_id":             "12",
		"resp_content_type":   "protobuf",
		"screen_height":       "1152",
		"screen_width":        "2048",
		"tz_name":             "Europe/Berlin",
		"referer":             "https://www.tiktok.com/",
		"root_referer":        "https://www.tiktok.com/",
		"host":                "https://webcast.tiktok.com",
		"webcast_sdk_version": "1.3.0",
		"update_version_code": "1.3.0",
	}
	reJsonData = []*regexp.Regexp{
		regexp.MustCompile(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__"[^>]+>(.*?)</script>`),
	}
	reVerify = regexp.MustCompile(`tiktok-verify-page`)
)

var (
	ErrUserOffline       = errors.New("user might be offline, Room ID not found")
	ErrIPBlocked         = errors.New("your IP or country might be blocked by TikTok")
	ErrJSONNotFound      = errors.New("unable to find required json data in HTML page")
	ErrLiveHasEnded      = errors.New("livestream has ended")
	ErrMsgNotImplemented = errors.New("message protobuf type has not been implemented, please report")
	ErrNoMoreFeedItems   = errors.New("no more feed items available")
	ErrUserNotFound      = errors.New("user not found")
	ErrCaptcha           = errors.New("captcha detected, unable to proceed")
	ErrURLNotFound       = errors.New("unable to download stream, URL not found")
	ErrFFMPEGNotFound    = errors.New("please install ffmpeg before downloading")
	ErrRateLimitExceeded = errors.New("you have exceeded the rate limit, please wait a few min")
	ErrUserInfoNotFound  = errors.New("user info not found")
)
