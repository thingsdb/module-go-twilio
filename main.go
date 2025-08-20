// ThingsDB module for using Twilio.
//
// For example:
//
//		// Create the module (@thingsdb scope)
//		new_module('twilio', 'github.com/thingsdb/module-go-twilio');
//
//		// Configure the module
//	  	set_module_conf("twilio", {
//	        TWILIO_ACCOUNT_SID: "REPLACE WITH ACCOUNT_SID",
//	        TWILIO_AUTH_TOKEN: "REPLACE WITH AUTH_TOKEN",
//	  	});
//
//		// Use the module
//		twilio.message( {
//		    to: "+31 612345678",
//		    from: "+31 687654321",
//		    body: "sample SMS message",
//		}).else(|err| {
//		    // error handling....
//		}));
package main

import (
	"fmt"
	"log"
	"os"

	twiml "github.com/homie-dev/gotwiml/twiml"
	attr "github.com/homie-dev/gotwiml/twiml/attr"
	voice "github.com/homie-dev/gotwiml/twiml/attr/const/voice"
	timod "github.com/thingsdb/go-timod"
	twilio "github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	msgpack "github.com/vmihailenco/msgpack"
)

type confTwilio struct {
	TwilioAccountSid string `msgpack:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken  string `msgpack:"TWILIO_AUTH_TOKEN"`
}

func handleConf(conf *confTwilio) error {
	if conf.TwilioAccountSid == "" {
		return fmt.Errorf("TWILIO_ACCOUNT_SID must not be empty")
	}
	if conf.TwilioAuthToken == "" {
		return fmt.Errorf("TWILIO_AUTH_TOKEN must not be empty")
	}
	os.Setenv("TWILIO_ACCOUNT_SID", conf.TwilioAccountSid)
	os.Setenv("TWILIO_AUTH_TOKEN", conf.TwilioAuthToken)
	return nil
}

type twilioMessage struct {
	Body             string `msgpack:"body"`
	To               string `msgpack:"to"`
	From             string `msgpack:"from"`
	ContentSid       string `msgpack:"sid"`
	ContentVariables string `msgpack:"variables"`
	DisableRiskCheck bool   `msgpack:"disableRiskCheck"`
}

type twilioCall struct {
	Body string `msgpack:"body"`
	To   string `msgpack:"to"`
	From string `msgpack:"from"`
}

type reqTwilio struct {
	Call    *twilioCall    `msgpack:"call"`
	Message *twilioMessage `msgpack:"message"`
}

func handleCall(pkg *timod.Pkg, call *twilioCall) {

	if call.Body == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio call requires a non empty `body`")
		return
	}

	if call.To == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio call requires a non empty `to`")
		return
	}

	if call.From == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio call requires a non empty `from`")
		return
	}

	client := twilio.NewRestClient()
	params := &api.CreateCallParams{}
	msg := twiml.NewVoiceResponse().Say(call.Body, attr.Voice(voice.Alice))

	txt, err := msg.ToXML()
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to create call body: %s", err.Error()))
		return
	}

	params.SetTwiml(txt)
	params.SetTo(call.To)
	params.SetFrom(call.From)

	resp, err := client.Api.CreateCall(params)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			fmt.Sprintf("Failed to create call: %s", err.Error()))
		return
	}

	timod.WriteResponse(pkg.Pid, resp.Sid)
}

func handleMessage(pkg *timod.Pkg, message *twilioMessage) {
	riskCheck := "enable"

	if message.Body == "" && message.ContentSid == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio message requires a non empty `body` or `sid`")
		return
	}

	if message.To == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio message requires a non empty `to`")
		return
	}

	if message.From == "" {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio message requires a non empty `from`")
		return
	}

	if message.DisableRiskCheck {
		riskCheck = "disable"
	}

	client := twilio.NewRestClient()
	params := &api.CreateMessageParams{}

	if message.Body != "" {
		params.SetBody(message.Body)
	}
	if message.ContentSid != "" {
		params.SetContentSid(message.ContentSid)
		params.SetContentVariables(message.ContentVariables)
	}
	params.SetTo(message.To)
	params.SetFrom(message.From)
	params.SetRiskCheck(riskCheck)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			fmt.Sprintf("Failed to create message: %s", err.Error()))
		return
	}

	timod.WriteResponse(pkg.Pid, resp.Sid)
}

func onModuleReq(pkg *timod.Pkg) {
	var req reqTwilio
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Failed to unpack Twilio request")
		return
	}

	if req.Call == nil && req.Message == nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio requires either `call` or `message`")
		return
	}

	if req.Call != nil && req.Message != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Twilio requires either `call` or `message`, not both")
		return
	}

	if req.Call != nil {
		go handleCall(pkg, req.Call)
	} else {
		go handleMessage(pkg, req.Message)
	}
}

func handler(buf *timod.Buffer, quit chan bool) {
	for {
		select {
		case pkg := <-buf.PkgCh:
			switch timod.Proto(pkg.Tp) {
			case timod.ProtoModuleConf:
				var conf confTwilio

				err := msgpack.Unmarshal(pkg.Data, &conf)
				if err != nil {
					log.Println("Missing or invalid Twilio configuration")
					timod.WriteConfErr()
					break
				}

				err = handleConf(&conf)
				if err != nil {
					log.Println(err.Error())
					timod.WriteConfErr()
					break
				}

				timod.WriteConfOk()

			case timod.ProtoModuleReq:
				onModuleReq(pkg)

			default:
				log.Printf("Unexpected package type: %d", pkg.Tp)
			}
		case err := <-buf.ErrCh:
			// In case of an error you probably want to quit the module.
			// ThingsDB will try to restart the module a few times if this
			// happens.
			log.Printf("Error: %s", err)
			quit <- true
		}
	}
}

func main() {
	// Starts the module
	timod.StartModule("twilio", handler)
}
