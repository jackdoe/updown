package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/miquella/ask"
	"gopkg.in/gomail.v2"
)

func main() {
	var pfrom = flag.String("ffrom", "", "from")
	var pto = flag.String("to", "", "to comma separated list")
	var puser = flag.String("user", "", "user")
	var pct = flag.String("content-type", "text/html", "content type")
	var psmtp = flag.String("smtp", "smtp.gmail.com", "smtp")
	var psbj = flag.String("subject", fmt.Sprintf("testing %d", time.Now().UnixNano()), "subject")
	flag.Parse()

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	pass, err := ask.HiddenAsk("Password: ")
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	if *pfrom == "" {
		m.SetHeader("From", *puser)
	} else {
		m.SetHeader("From", *pfrom)
	}
	m.SetHeader("To", strings.Split(*pto, ",")...)
	m.SetHeader("Subject", *psbj)
	m.SetBody(*pct, string(in))

	d := gomail.NewPlainDialer(*psmtp, 587, *puser, pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
