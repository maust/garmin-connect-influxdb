package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type DateObject struct {
	Date string `json:"date"`
}

type DateSummary struct {
	UserPK                     string     `json:"userProfilePK"`
	CalendarDate               DateObject `json:"calendarDate"`
	UUID                       string     `json:"uuid"`
	Duration                   int        `json:"durationInMilliseconds"`
	ActiveKilocalories         int        `json:"activeKilocalories"`
	BmrKilocalories            int        `json:"bmrKilocalories"`
	WellnessKilocalories       int        `json:"wellnessKilocalories"`
	RemainingKilocalories      int        `json:"remainingKilocalories"`
	WellnessTotalKilocalories  int        `json:"wellnessTotalKilocalories"`
	WellnessActiveKilocalories int        `json:"wellnessActiveKilocalories"`
	TotalSteps                 int        `json:"totalSteps"`
	DailyStepGoal              int        `json:"dailyStepGoal"`
	TotalDistanceMeters        int        `json:"totalDistanceMeters"`
	WellnessDistanceMeters     int        `json:"wellnessDistanceMeters"`
	HighlyActiveSeconds        int        `json:"highlyActiveSeconds"`
	ModerateIntensityMinutes   int        `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes   int        `json:"vigorousIntensityMinutes"`
	FloorsAscendedInMeters     int        `json:"floorsAscendedInMeters"`
	UserIntensityMinutesGoal   int        `json:"userIntensityMinutesGoal"`
	UserFloorsAscendedGoal     int        `json:"userFloorsAscendedGoal"`
	MinHeartRate               int        `json:"minHeartRate"`
	MaxHeartRate               int        `json:"maxHeartRate"`
	RestingHeartRate           int        `json:"restingHeartRate"`
	CurrentDayRestingHeartRate int        `json:"currentDayRestingHeartRate"`
}

func processFile(c client.Client, file string) {
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened file")
	}

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var values []DateSummary

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &values)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "garmin",
	})

	var layout = "Jan 2, 2006 03:04:05 PM"

	for i := 0; i < len(values); i++ {
		t, err := time.Parse(layout, values[i].CalendarDate.Date)
		if err != nil {
			fmt.Println("Error while parsing date")
		}
		pt := writePoint(t, values[i])
		bp.AddPoint(pt)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	err = c.Write(bp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func writePoint(timestamp time.Time, values DateSummary) *client.Point {
	tags := map[string]string{}

	fields := map[string]interface{}{
		"ActiveKilocalories":         values.ActiveKilocalories,
		"BmrKilocalories":            values.BmrKilocalories,
		"WellnessKilocalories":       values.WellnessKilocalories,
		"RemainingKilocalories":      values.RemainingKilocalories,
		"WellnessTotalKilocalories":  values.WellnessTotalKilocalories,
		"WellnessActiveKilocalories": values.WellnessActiveKilocalories,
		"TotalSteps":                 values.TotalSteps,
		"DailyStepGoal":              values.DailyStepGoal,
		"TotalDistanceMeters":        values.TotalDistanceMeters,
		"WellnessDistanceMeters":     values.WellnessDistanceMeters,
		"HighlyActiveSeconds":        values.HighlyActiveSeconds,
		"ModerateIntensityMinutes":   values.ModerateIntensityMinutes,
		"VigorousIntensityMinutes":   values.VigorousIntensityMinutes,
		"FloorsAscendedInMeters":     values.FloorsAscendedInMeters,
		"UserIntensityMinutesGoal":   values.UserIntensityMinutesGoal,
		"UserFloorsAscendedGoal":     values.UserFloorsAscendedGoal,
		"MinHeartRate":               values.MinHeartRate,
		"MaxHeartRate":               values.MaxHeartRate,
		"RestingHeartRate":           values.RestingHeartRate,
		"CurrentDayRestingHeartRate": values.CurrentDayRestingHeartRate,
	}

	pt, err := client.NewPoint(
		"health",
		tags,
		fields,
		timestamp,
	)
	if err != nil {
		println("Error:", err.Error())
	}
	return pt
}

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	var folder = "./export/DI_CONNECT/Di-Connect-User/"

	file, err := os.Open(folder)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		if strings.HasPrefix(name, "UDSFile") {
			fmt.Println(name)
			processFile(c, folder+name)
		}
	}
}
