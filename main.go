package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DB : Any database
type DB struct {
	Campaigns []Campaign
	LastIndex int
}

// GetCampaigns : Returns whole contents of the database
func (d *DB) GetCampaigns() []Campaign {
	return d.Campaigns
}

// GetCampaign : Returns a campaign matching the id argument
func (d *DB) GetCampaign(id int) (found bool, campaign Campaign) {
	for _, campaign := range d.Campaigns {
		if campaign.ID == id {
			return true, campaign
		}
	}
	return false, Campaign{}
}

// SaveCampaign : Adds a new campaign to the database
func (d *DB) SaveCampaign(campaign Campaign) Campaign {
	d.LastIndex++
	campaign.ID = d.LastIndex
	d.Campaigns = append(d.Campaigns, campaign)
	return campaign
}

// DeleteCampaign : Remove a campaign by ID
func (d *DB) DeleteCampaign(id int) (campaign Campaign, err error) {
	for index, campaign := range d.Campaigns {
		if campaign.ID == id {
			d.Campaigns = append(d.Campaigns[:index], d.Campaigns[index+1:]...)
			return campaign, nil
		}
	}
	return Campaign{}, fmt.Errorf("Could not find campaign %d", id)
}

// Campaign : An AwesomeAds campaign
type Campaign struct {
	ID           int           `json:"id,omitempty"`
	Name         string        `json:"name"`
	Company      string        `json:"company"`
	IO           string        `json:"io"`
	House        bool          `json:"house"`
	SubCampaigns []SubCampaign `json:"subcampaigns"`
}

// SubCampaign : An AwesomeAds sub-campaign
type SubCampaign struct {
	ID         int    `json:"id,omitempty"`
	CampaignID int    `json:"campaign_id"`
	Name       string `json:"name"`
	SubIO      string `json:"sub_io"`
	Countries  string `json:"countries"`
	Devices    string `json:"devices"`
}

// GetCampaigns : Get a list of all AwesomeAds campaigns
func GetCampaigns(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(db.GetCampaigns())
}

// GetCampaign : Get an AwesomeAds campaign by ID
func GetCampaign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ok, campaign := db.GetCampaign(campaignID); ok {
		json.NewEncoder(w).Encode(campaign)
		return
	}
	http.Error(w, fmt.Sprintf("Campaign %d not found", campaignID), http.StatusNotFound)
}

// SaveCampaign : Save a new AwesomeAds campaign
func SaveCampaign(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var campaign Campaign
	err := decoder.Decode(&campaign)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	savedCampaign := db.SaveCampaign(campaign)
	json.NewEncoder(w).Encode(savedCampaign)
}

// DeleteCampaign : Delete an AwesomeAds campaign
func DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	campaign, err := db.DeleteCampaign(campaignID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(campaign)
}

func initDB() *DB {
	content, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Fatal(err)
	}

	var campaigns []Campaign
	err = json.Unmarshal(content, &campaigns)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &DB{
		Campaigns: campaigns,
		LastIndex: 100,
	}
}

var db = initDB()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/campaigns", GetCampaigns).Methods("GET")
	router.HandleFunc("/campaign/{id}", GetCampaign).Methods("GET")
	router.HandleFunc("/campaign", SaveCampaign).Methods("POST")
	router.HandleFunc("/campaign/{id}", DeleteCampaign).Methods("DELETE")

	/*
		In the Hello World example from Session 1, we initalised the server like this
		log.Fatal(http.ListenAndServe(":8080", nil))
		ListenAndServe starts an HTTP server with a given address and handler.
		The handler is usually nil, which means to use DefaultServeMux.
		Handle and HandleFunc add handlers to DefaultServeMux
		e.g.: http.HandleFunc("/", handler)

		A multiplexer is usually something that does routing
	*/

	log.Fatal(http.ListenAndServe(":8080", router))
}
