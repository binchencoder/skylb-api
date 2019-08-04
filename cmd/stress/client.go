package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"

	"github.com/binchencoder/letsgo"
	"github.com/binchencoder/letsgo/token"
	pb "github.com/binchencoder/skylb-api/cmd/stress/proto"
)

const urlPattern = "http://%s/stress/v1/say-hello/dory/gender/FEMALE"
const disabledUrlPattern = "http://%s/stress/v1/say-hello-disabled/dory/gender/FEMALE"
const signingKey = "BX4DR-FY8YF-SECRET-G7RHK-KBXFV-"

var (
	flagServer = flag.String("server", "", "The server host:port")
	onceOnly   = flag.Bool("send-once", false, "Send the request once only")
	clientId   = flag.String("client-id", "", "Client ID")
	xSource    = flag.String("x-source", "web", "X-Source to use")
)

func main() {
	letsgo.Init()

	url := fmt.Sprintf(urlPattern, *flagServer)
	disabledUrl := fmt.Sprintf(disabledUrlPattern, *flagServer)

	// If specified, send request once and return.
	if *onceOnly {
		sendRequest(url)
		sendRequest(disabledUrl)
		return
	}

	for i := 0; i < 10; i++ {
		go func() {
			for range time.Tick(50 * time.Millisecond) {
				sendRequest(url)
			}
		}()
	}

	select {}
}

func createToken(clientId string) (string, error) {
	info := token.TokenClientInfo{
		ClientId: "clientId",
		CorpCode: "TestCid",
	}

	return token.CreateSignedJwtToken(signingKey, &info, 2000)
}

func sendRequest(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return
	}
	req.Header.Set("x-source", *xSource)
	req.Header.Set("Content-Type", "application/json")

	if *clientId != "" {
		token, err2 := createToken(*clientId)
		if err2 != nil {
			panic(fmt.Sprintf("%v", err2))
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set("x-Uid", "marlin")
		req.Header.Set("X-Cid", "disney")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		record := pb.SayHelloResp{}
		if err = jsonpb.UnmarshalString(string(bodyBytes), &record); err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%d, %s from %s.\n", resp.StatusCode, record.Greeting, record.Peer)
	} else {
		fmt.Println(resp.StatusCode)
	}
}
