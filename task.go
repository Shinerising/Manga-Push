package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func startTask() {
	fmt.Println("Start Task!")
	_, err := os.Stat("./config/")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./config/", os.ModePerm)
	}

	ticker := time.NewTicker(4 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				taskHandler()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func taskHandler() {
	config := Config{}
	_, err := os.Stat("./config/config.json")
	if os.IsNotExist(err) {
		config.LastUpdate = 10210
	} else {
		getJson("./config/config.json", &config)
	}

	id := config.LastUpdate
	sum := 0
	for sum < 10 {
		fmt.Println("Checking " + strconv.Itoa(id+sum))
		success, bookid := bookManagement(strconv.Itoa(id + sum))
		if success {
			pushBook(bookid, strconv.Itoa(id+sum))
			sum += 1
		} else {
			break
		}
	}
	config.LastUpdate = id + sum
	jsonString, _ := json.Marshal(config)
	err = ioutil.WriteFile("./config/config.json", jsonString, 0644)
	if err != nil {
		log.Panic(err)
		return
	}
}
