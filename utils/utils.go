package utils

import (
	"regexp"
)

const (
	// Alphanumeric charset
	alphaNumCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// RFC 3986 Section 2.3 URI Unreserved Characters
	uriUnreservedChars = `^([A-Za-z0-9_.~-])+$`

	// https://gist.github.com/dperini/729294
	urlRegex = `^` +
		// protocol identifier (optional)
		// short syntax // still required
		`(?:(?:(?:https?|ftp):)?\/\/)` +
		// user:pass BasicAuth (optional)
		`(?:\S+(?::\S*)?@)?` +
		`(?:` +
		// IP address dotted notation octets
		// excludes loopback network 0.0.0.0
		// excludes reserved space >= 224.0.0.0
		// excludes network & broadcast addresses
		// (first & last IP address of each class)
		`(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])` +
		`(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}` +
		`(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))` +
		`|` +
		// host & domain names, may end with dot
		// can be replaced by a shortest alternative
		// (?![-_])(?:[-\w\u00a1-\uffff]{0,63}[^-_]\.)+
		`(?:` +
		`(?:` +
		`[a-z0-9\\u00a1-\\uffff]` +
		`[a-z0-9\\u00a1-\\uffff_-]{0,62}` +
		`)?` +
		`[a-z0-9\\u00a1-\\uffff]\.` +
		`)+` +
		// TLD identifier name, may end with dot
		`(?:[a-z\\u00a1-\\uffff]{2,}\.?)` +
		`)` +
		// port number (optional)
		`(?::\d{2,5})?` +
		// resource path (optional)
		`(?:[/?#]\S*)?` +
		`$`
)

const selfLink = "shorturl.ishaandham.com/.*"

// String is valid url
func IsValidURL(str string) bool {
	valid_1, err_1 := regexp.MatchString(urlRegex, str)
	valid_2, err_2 := regexp.MatchString(selfLink, str)
	return (err_1 == nil && err_2 == nil) && (!valid_2 && valid_1)
}

// String has chars - alphabets, numbers, underscores
func IsValidAlias(str string) bool {
	valid, err := regexp.MatchString("^[A-Za-z0-9_]+$", str)
	return valid && (err == nil) && (len(str) < 50)
}
