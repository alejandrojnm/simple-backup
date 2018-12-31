package main

import (
	"flag"
	"github.com/alejandrojnm/simple-backup/core"
	"github.com/minio/minio-go"
	"github.com/vjeantet/jodaTime"
	"log"
	"log/syslog"
	"path/filepath"
	"time"
)

func init() {
	// Init log
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "simple-backup")
	if e == nil {
		log.SetOutput(logwriter)
	}
}

func main() {

	//Flags for the program
	//S3 data
	endPoint := flag.String("endpoint", "localhost:9000", "Endpoint to S3")
	ak := flag.String("ak", "accesskey", "Access key")
	sk := flag.String("sk", "secretkey", "Secret key")
	secureEndPoint := flag.Bool("secureendpoint", false, "Enabled secure endpoint")
	// other options
	bucket := flag.String("bucket", "alamesa-db", "Buket Name")
	bucketLocation := flag.String("bucketlocation", "us-east-1", "Bucket Zone")
	backupDir := flag.String("backupdir", "backup", "Backup directory location")
	flag.Parse()
	//

	// Create a new client for the storage server
	s3Client, err := minio.New(*endPoint, *ak, *sk, *secureEndPoint)
	if err != nil {
		log.Print(err)
	}

	//We send to create the bucket, if we do not do anything
	err = s3Client.MakeBucket(*bucket, *bucketLocation)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := s3Client.BucketExists(*bucket)
		if err == nil && exists {
			log.Printf("We already own %s\n", *bucket)
		} else {
			log.Print(err)
		}
	}

	todayRegx := []string{"*"}
	listFile, _ := core.FindFile(*backupDir+"/", todayRegx)

	//We verify that there is something to upload and that there is no error
	if len(listFile) != 0 {
		// for all file in listFile we send to storage server
		for _, file := range listFile {
			fileName := filepath.Base(file)
			//log.Printf(file_name)
			n, err := s3Client.FPutObject(*bucket, jodaTime.Format("YYYYMMdd", time.Now())+"/"+fileName, file, minio.PutObjectOptions{ContentType: ""})
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Successfully uploaded %s of size %d\n", fileName, n)
		}
	}

	objectsCh := make(chan string)
	// Send object names that are needed to be removed to objectsCh in subroutine
	go func() {
		defer close(objectsCh)
		doneCh := make(chan struct{})

		// Indicate to our routine to exit cleanly upon return.
		defer close(doneCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range s3Client.ListObjects(*bucket, jodaTime.Format("YYYYMMdd", time.Now().Add(-144*time.Hour)), true, doneCh) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object.Key
		}
	}()

	// Call RemoveObjects API
	errorCh := s3Client.RemoveObjects(*bucket, objectsCh)
	// Print errors received from RemoveObjects API
	for e := range errorCh {
		log.Fatalln("Failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}

}
