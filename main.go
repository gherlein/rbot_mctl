package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math"
	"strconv"
	"strings"
)

var (
	support_xy      bool   = true
	support_vector  bool   = false
	support_buttons bool   = true
	debugraw        bool   = false
	debugvector     bool   = false
	debugjoy        bool   = false
	debugbutton     bool   = false
	rbroker         string = "tcp://<rbot-ip>:1883"
	rtext           string = "pi-blaster-mqtt/text"
	jbroker         string = "tcp://localhost:1883"
	joysticks       string = "xb/1/joysticks"
	triggers        string = "xb/1/triggers"
	buttons         string = "xb/1/buttons"
	xb              string = "xb/#"
	qos             int    = 0
	rbot            MQTT.Client
	joy             MQTT.Client
	m1pin           int = 4
)

func init() {
	flag.StringVar(&rbroker, "rbotbroker", "tcp://rbot-ip:1883",
		"rbot broker connection string")
	flag.StringVar(&jbroker, "joybroker", "tcp://localhost:1883",
		"rbot broker connection string")
	flag.Parse()
}

func onXB(client MQTT.Client, message MQTT.Message) {
	var jmatch, bmatch, tmatch int

	m := string(message.Payload())
	s := strings.Split(m, "|")

	fmt.Printf("JOY--> {%s}\n", m)

	jmatch = strings.Compare(joysticks, message.Topic())
	bmatch = strings.Compare(buttons, message.Topic())
	tmatch = strings.Compare(triggers, message.Topic())

	if jmatch == 0 {
		var msg string
		if s[0] == "L" {
			if s[1] == "Y" {
				y, _ := strconv.Atoi(s[3])
				var v float64 = float64(y) / float64(math.MaxInt16)
				o := float64(0.05) * float64(v)
				n := 0.15 + o

				// totally faking it - only one pin, hardcoded
				msg = fmt.Sprintf("%d=%f", m1pin, n)
				fmt.Printf("%s\n", msg)
			}
			token := rbot.Publish(rtext, 0, false, msg)
			token.Wait()
			fmt.Printf("%s\n", msg)

		}
	}

	if bmatch == 0 {
	}
	if tmatch == 0 {
	}
}

func onRBOT(client MQTT.Client, message MQTT.Message) {

	m := string(message.Payload())
	fmt.Printf("RBOT--> {%s}\n", m)
}

func main() {

	optsJ := MQTT.NewClientOptions()
	optsJ.AddBroker(jbroker)
	optsJ.OnConnect = func(c MQTT.Client) {
		fmt.Printf("Connected to joystick broker %s\n", jbroker)
		if token := c.Subscribe(xb, 0, onXB); token.Wait() &&
			token.Error() != nil {
			panic(token.Error())
		}
	}
	joy = MQTT.NewClient(optsJ)
	if tokenJ := joy.Connect(); tokenJ.Wait() && tokenJ.Error() != nil {
		panic(tokenJ.Error())
	}

	optsR := MQTT.NewClientOptions()
	optsR.AddBroker(rbroker)
	optsR.OnConnect = func(c MQTT.Client) {
		fmt.Printf("Connected to bot broker %s\n", rbroker)
		if token := c.Subscribe(rtext, 0, onRBOT); token.Wait() &&
			token.Error() != nil {
			panic(token.Error())
		}
	}
	rbot = MQTT.NewClient(optsR)
	if tokenR := rbot.Connect(); tokenR.Wait() && tokenR.Error() != nil {
		fmt.Println("Error connecting to rbot motor event queue")
		panic(tokenR.Error())
	}
	for {

	}
}
