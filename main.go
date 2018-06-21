package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	gomail "gopkg.in/gomail.v2"
)

func main() {

	args := os.Args[1:]

	if len(args) != 5 {
		panic("Please provide exactly 5 argument in the order [email] [email_password] [email2notify] [cpu_threshold] [check_every_X_minute]")
	}

	emailUsername := args[0]
	emailPassword := args[1]
	email2Notify := args[2]
	cpuThreshold := args[3]
	checkEveryXMinute := args[4]

	d := gomail.NewDialer("smtp.gmail.com", 587, emailUsername, emailPassword)

	xMinute, err := strconv.ParseInt(checkEveryXMinute, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("You provided X minutes %d, please try integer.", xMinute))
	}

	c := time.Tick(time.Duration(xMinute) * time.Minute)
	for range c {
		v, _ := cpu.Percent(0, false)

		if len(v) != 1 {
			panic("Something is wrong with program, cpu percentage should be a one element array.")
		}

		// almost every return value is a struct
		fmt.Printf("UsedPercent:%f\n", v)

		thresholdFloat, err := strconv.ParseFloat(cpuThreshold, 32)
		if err != nil {
			panic(fmt.Sprintf("You provided threshold %s, which is not a valid float number.", cpuThreshold))
		}

		if v[0] > thresholdFloat {
			m := gomail.NewMessage()
			m.SetHeader("From", emailUsername)
			m.SetHeader("To", email2Notify)
			m.SetHeader("Subject", fmt.Sprintf("[Warning] CPU Usage at %f!", v))
			m.SetBody("text/html", "Hello <b>Walter</b> and <i>CPU usage at your PC is HIGH</i>!")

			// Send the email to Bob, Cora and Dan.
			if err := d.DialAndSend(m); err != nil {
				panic(err)
			}
		}
	}

}
