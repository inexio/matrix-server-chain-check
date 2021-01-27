package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/inexio/go-monitoringplugin"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"os"
	"strconv"
	"time"
)

// sender credentials
var sendingHomeServer = flag.String("sending-homeserver", "", "Sending Matrix homeserver")
var sendingUsername = flag.String("sending-username", "", "Sending Matrix username")
var sendingPassword = flag.String("sending-password", "", "Sending Matrix password")

// receiver credentials
var receivingHomeServer = flag.String("receiving-homeserver", "", "Receiving Matrix homeserver")
var receivingUsername = flag.String("receiving-username", "", "Receiving Matrix username")
var receivingPassword = flag.String("receiving-password", "", "Receiving Matrix password")

var timeoutSeconds = flag.Int("timeout", 10, "Timeout before throwing critical error if there is no matching message received")
var roomID = flag.String("room-id", "", "ID of the Matrix room")

type ErrorAndCode struct {
	ExitCode int
	Error    error
}

func main() {
	var errSlice []ErrorAndCode
	flag.Parse()

	if *sendingHomeServer == "" || *sendingPassword == "" || *sendingUsername == "" || *receivingHomeServer == "" || *receivingUsername == "" || *receivingPassword == "" || *roomID == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	//fmt.Println("Logging into", *sendingHomeServer, "as", *sendingUsername)
	client, err := mautrix.NewClient(*sendingHomeServer, "", "")
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.New("Homeserver not known " + *sendingHomeServer)})
		OutputMonitoring(errSlice, "checked")
	}
	_, err = client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: *sendingUsername},
		Password:         *sendingPassword,
		StoreCredentials: true,
	})
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.New("Could not login as " + *sendingUsername + " to " + *sendingHomeServer)})
		OutputMonitoring(errSlice, "checked")
	}

	sendingText := "chainTestText" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	_, err = client.SendText(id.RoomID(*roomID), sendingText)
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.New("Could not send message in" + *roomID)})
		OutputMonitoring(errSlice, "checked")
	}

	client2, err := mautrix.NewClient(*receivingHomeServer, "", "")
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.New("Homeserver not known" + *receivingHomeServer)})
		OutputMonitoring(errSlice, "checked")
	}
	_, err = client2.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: *receivingUsername},
		Password:         *receivingPassword,
		StoreCredentials: true,
	})
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.New("Could not login to " + *sendingHomeServer)})
		OutputMonitoring(errSlice, "checked")
	}

	signal := make(chan bool, 1)
	errChan := make(chan bool, 1)
	monErrChan := make(chan ErrorAndCode)

	syncer := client2.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		//fmt.Printf("%[5]s <%[1]s> %[4]s (%[2]s/%[3]s)\n", evt.Sender, evt.Type.String(), evt.ID, evt.Content.AsMessage().Body, evt.Timestamp)
		err := client2.MarkRead(id.RoomID(*roomID), evt.ID)
		if err != nil {
			monErrChan <- ErrorAndCode{3, errors.New("Could not mark message as read")}
		}
		if evt.Content.AsMessage().Body == sendingText {
			monErrChan <- ErrorAndCode{0, errors.New("The chain check was successfull")}
			signal <- true
		}
	})

	go func() {
		err = client2.Sync()
		if err != nil {
			monErrChan <- ErrorAndCode{3, errors.New("sync stopped with error")}
		}
	}()

	go func() {
		time.Sleep(time.Duration(*timeoutSeconds) * time.Second)
		errChan <- true
	}()

	select {
	case <-signal:
		client2.StopSync()
	case <-errChan:
		client2.StopSync()
		errSlice = append(errSlice, ErrorAndCode{2, errors.New("Message was not received")})
	}
	for i := range monErrChan {
		errSlice = append(errSlice, i)
	}
	OutputMonitoring(errSlice, "checked")
}

func OutputMonitoring(errSlice []ErrorAndCode, defaultMessage string) {
	response := monitoringplugin.NewResponse(defaultMessage)
	for i := 0; i < len(errSlice); i++ {
		response.UpdateStatus(errSlice[i].ExitCode, errSlice[i].Error.Error())
	}
	response.OutputAndExit()
}
