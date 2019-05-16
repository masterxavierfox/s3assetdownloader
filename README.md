# s3assetdownloader

Simple s3 client to download assets from s3 bucket.

###Usage

set the below `env varibles` or use `.env` file

    JSON_URL = https://<endpoint>.json
    AWS_ID = <aws_access_key_id>
    AWS_SECRET = <aws_secret_access_key>
    S3_BUCKET = <bucket_name>
    
The code uses a custom json stracture like below:

    TrackID          string `json:"track_id"`
    TrackTitle       string `json:"track_title"`
    TrackStreetTitle string `json:"track_street_title"`
    AlbumTitle       string `json:"album_title"`
    ArtistName       string `json:"artist_name"`
    ProviderName     string `json:"provider_name"`
    StreamURL        string `json:"stream_url"`
    Genres           string `json:"genres"`
    SongGeneratedID  string `json:"song_generated_id"`
    SongID           string `json:"song_id"`    


You can change this from the `main.go` file.


### Build

make sure you have golang dep installed, then run:

    dep ensure
    
then run the build script:

    bash ./build.sh
    
 ### TO DO
 
1.Standardize input struct/source