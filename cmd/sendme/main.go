package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/miquella/ask"
	"gopkg.in/gomail.v2"
)

func main() {
	var pfrom = flag.String("ffrom", "", "from")
	var pto = flag.String("to", "", "to")
	var puser = flag.String("user", "", "user")
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
	m.SetHeader("To", *pto)
	m.SetHeader("Subject", *psbj)
	m.SetBody("text/html", string(in))

	d := gomail.NewPlainDialer(*psmtp, 587, *puser, pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
