package main

import (
	"flag"
	"os"

	"github.com/jackdoe/updown/util"
	"github.com/satyrius/gonx"
	"zgo.at/isbot"
)

func main() {
	var format = flag.String("format", "$remote_addr - $remote_user [$time_local] \"$request\" $status $size \"$http_referer\" \"$http_user_agent\"", "format")
	flag.Parse()
	parser := gonx.NewParser(*format)
	util.ForeachLine(os.Stdin, func(text string, _ bool) {
		rec, err := parser.ParseString(text)
		if err != nil {
			return
		}

		ip, _ := rec.Field("remote_addr")
		ua, _ := rec.Field("http_user_agent")

		a := isbot.UserAgent(ua)
		b := isbot.IPRange(ip)
		if a+b < 2 {
			os.Stdout.Write([]byte(text))
			os.Stdout.Write([]byte("\n"))
		}
	})
}
