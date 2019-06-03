package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"net/http"
	"os"
	"strings"
)

type songProperties []struct {
	TrackID string `json:"track_id"`
	//TrackTitle       string `json:"track_title"`
	//TrackStreetTitle string `json:"track_street_title"`
	//AlbumTitle       string `json:"album_title"`
	//ArtistName       string `json:"artist_name"`
	//ProviderName     string `json:"provider_name"`
	//StreamURL        string `json:"stream_url"`
	Genres           string `json:"genres"`
	SongGeneratedID  string `json:"song_generated_id"`
	SongID string `json:"song_id"`
}

func main() {
	getAssetinfo()
}

func getAssetinfo() {

	//Check if .env file exists
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//load environment
	jsonurl := os.Getenv("JSON_URL")

	// Build the request
	req, err := http.NewRequest("GET", jsonurl, nil)
	if err != nil {
		log.Fatal(err, "NewRequest: ")
		return
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, "Do: ")
		return
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var Record songProperties

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&Record); err != nil {
		log.Println(err)
	}

	fmt.Println(Record[0].TrackID)

	//S3 downloader
	AWSAuth := aws.Auth{
		AccessKey: os.Getenv("AWS_ID"),
		SecretKey: os.Getenv("AWS_SECRET"),
	}

	region := aws.EUWest
	// change this to your AWS region
	// click on the bucketname in AWS control panel and click Properties
	// the region for your bucket should be under "Static Website Hosting" tab

	connection := s3.New(AWSAuth, region)

	bucket := connection.Bucket(os.Getenv("S3_BUCKET")) // change this your bucket name

	//path := ""  // this is the target file and location in S3



	for  R := range Record {

		// Download(GET)
		fmt.Println("  ============================= Fetching: " + Record[R].TrackID + ".mp4 ==========================")

		downloadBytes, err := bucket.Get(Record[R].TrackID + ".mp4")

		if err != nil {
			fmt.Println(err)
			log_missing_keys(Record[R].TrackID)
			panicAndRecover("CANT FETCH ::: #*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*  404 TRACK: "+ Record[R].TrackID +" MISSING #*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#* ")
			//os.Exit(1)

		}

		//downloadFile, err := os.Create(Record[R].SongID + ".mp4")

		downloadFile, err := os.Create(strings.ToLower(Record[R].Genres) + " " + Record[R].SongGeneratedID + ".mp4")

		if err != nil {
			fmt.Println(err)
			panicAndRecover("CANT WRITE ::: #*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*  404 TRACK: "+ Record[R].TrackID +" MISSING #*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#*#* ")
			//os.Exit(1)
		}

		defer downloadFile.Close()

		downloadBuffer := bufio.NewWriter(downloadFile)
		downloadBuffer.Write(downloadBytes)

		io.Copy(downloadBuffer, downloadFile)

		fmt.Printf("Downloaded from S3 and saved to: " + Record[R].SongID + ".mp4 \n\n")

	}


}

func log_missing_keys(assetKey string){
	f, err := os.OpenFile("missingkey.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(assetKey+"\n"); err != nil {
		log.Println(err)
	}
}

func panicAndRecover(message string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic(message)
}