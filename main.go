package main

import (
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
	rbroker         string = "tcp://192.168.2.18:1883"
	rtext           string = "pi-blaster-mqtt/text"
	jbroker         string = "tcp://localhost:1883"
	xy              string = "xb/1/joy-xy"
	vector          string = "xb/1/joy-vector"
	buttons         string = "xb/1/buttons"
	qos             int    = 0
	xmult           int16  = 1
	ymult           int16  = -1
	rbot            MQTT.Client
	joy             MQTT.Client
	m1pin           int = 4
)

func init() {
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {

	//	fmt.Printf("->[%s]\n", message.Payload())

	m := string(message.Payload())
	s := strings.Split(m, "|")

	var msg string

	if s[0] == "L" {
		if s[1] == "Y" {
			y, _ := strconv.Atoi(s[3])
			var v float64 = float64(y) / float64(math.MaxInt16)
			o := float64(0.05) * float64(v)
			n := 0.15 + o
			//			fmt.Printf("%d %v %f %f\n", y, v, o, n)
			msg = fmt.Sprintf("%d=%f", m1pin, n)
			fmt.Printf("%s\n", msg)
		}
		token := rbot.Publish(rtext, 0, false, msg)
		token.Wait()
	}

}

func main() {

	optsJ := MQTT.NewClientOptions()
	optsJ.AddBroker(jbroker)

	optsJ.OnConnect = func(c MQTT.Client) {
		fmt.Printf("Connected to %s\n", jbroker)
		if token := c.Subscribe(xy, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	joy = MQTT.NewClient(optsJ)
	if tokenJ := joy.Connect(); tokenJ.Wait() && tokenJ.Error() != nil {
		panic(tokenJ.Error())
	}

	optsR := MQTT.NewClientOptions()
	optsR.AddBroker(rbroker)
	rbot = MQTT.NewClient(optsR)
	if tokenR := rbot.Connect(); tokenR.Wait() && tokenR.Error() != nil {
		fmt.Println("Error connecting to rbot motor event queue")
		panic(tokenR.Error())
	}
	for {

	}
}
