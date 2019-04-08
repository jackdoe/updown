package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/miquella/ask"
	"gopkg.in/gomail.v2"
)

func main() {
	var pinhtml = flag.String("html", "", "html file")
	var pintext = flag.String("text", "", "text file")
	var pfrom = flag.String("ffrom", "", "from")
	var pto = flag.String("to", "", "to comma separated list")
	var puser = flag.String("user", "", "user")
	var psmtp = flag.String("smtp", "smtp.gmail.com", "smtp")
	var psbj = flag.String("subject", fmt.Sprintf("testing %d", time.Now().UnixNano()), "subject")
	flag.Parse()

	in, err := ioutil.ReadFile(*pinhtml)
	if err != nil {
		panic(err)
	}

	intext, err := ioutil.ReadFile(*pintext)
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
	m.SetBody("text/plain", string(intext))
	m.SetHeader("Content-Transfer-Encoding", "quoted-printable")
	m.AddAlternative("text/html", string(in))
	d := gomail.NewPlainDialer(*psmtp, 587, *puser, pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
