package main

import (
	"encoding/json"
	"fmt"
	"github.com/PNarode/vaccine_notifier/helper"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	slack_url = "https://hooks.slack.com/services/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"
	username = "U01SNDMA8RM"
	channel = "coreteam"
	url = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict"
)

type Session struct {
	AvailableCapacity int `json:"available_capacity"`
	MinAgeLimit int `json:"min_age_limit"`
	Vaccine string `json:"vaccine"`
}

type Center struct {
	Address string `json:"address"`
	BlockName string `json:"block_name"`
	Name string `json:"name"`
	Pincode int `json:"pincode"`
	Sessions []Session `json:"sessions"`
}

type Body struct {
	Centers []Center `json:"centers"`
}

func notify(c Center, s Session, client helper.SlackClient)  {
	msg := fmt.Sprintf("******Vaccine Slot*****\nCenter Name: %s \nAddress: %s \nPin Code: %d\nVaccine Name: %s\nAvailable:%d\nMinAgeLimit:%d",
		c.Name, c.Address, c.Pincode, s.Vaccine, s.AvailableCapacity, s.MinAgeLimit)
	sr := helper.SimpleSlackRequest{
		Text: msg,
		IconEmoji: ":syringe:",
	}
	err := client.SendSlackNotification(sr)
	if err != nil {
		fmt.Println("Failed to send slack notification", err)
	}
	helper.SendEmail(msg)
	return
}

func processBody(body Body, client helper.SlackClient) {
	for _, c := range body.Centers {
		for _, s := range c.Sessions {
			if s.MinAgeLimit > 18 && s.MinAgeLimit < 45 && s.AvailableCapacity > 0 {
				go notify(c, s, client)
				break
			}
		}
	}
}

func check_slot(wg *sync.WaitGroup, ID int, sclient helper.SlackClient) {

	client := &http.Client{Timeout: time.Minute * 5}
	var respBody Body
	var resp *http.Response
	defer func() {
		wg.Done()
		resp.Body.Close()
	}()
	for {
		today := time.Now().Format("02-01-2006")
		disUrl := fmt.Sprintf("%s?district_id=%d&date=%s", url, ID, today)
		req, err := http.NewRequest(http.MethodGet, disUrl, nil)
		if err != nil {
			fmt.Println("Failed to make HTTP call ", err)
			time.Sleep(time.Minute * 5)
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Failed to make HTTP call ", err)
			time.Sleep(time.Minute * 5)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("No Response Body ", err)
			time.Sleep(time.Minute * 5)
			continue
		}
		err = json.Unmarshal(body, &respBody)
		if err != nil {
			fmt.Println("Error UnMarshalling Body ", err)
			time.Sleep(time.Minute * 5)
			continue
		}
		processBody(respBody, sclient)
		time.Sleep(time.Minute * 5)
	}
}

func main() {
	var wg sync.WaitGroup
	district := []int{363, 389, 391}
	client := helper.SlackClient{
		WebHookUrl: slack_url,
		UserName:   username,
		Channel:    channel,
	}
	//s := Session{Vaccine: ""}
	//c := Center{Name: "Test Message"}
	fmt.Println("Start Check for Vaccine Slots")
	//notify(c, s, client)
	for _, i := range district {
		wg.Add(1)
		go check_slot(&wg, i, client)
	}
	wg.Wait()
	fmt.Println("Finished")
}
