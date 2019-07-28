// gtts.go
package zuke

import (
	"fmt"
	"net/url"
)

func GetUrlFromSpeech(text string, lang string) string {
	return fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), url.QueryEscape(lang))
}
