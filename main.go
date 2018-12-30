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

	//Flags
	//S3
	endpoint := flag.String("endpoint", "localhost:9000", "Endpoint to S3")
	ak := flag.String("ak", "accesskey", "Access key")
	sk := flag.String("sk", "secretkey", "Secret key")
	secureendpoint := flag.Bool("secureendpoint", false, "Enabled secure endpoint")
	//
	bucket := flag.String("bucket", "alamesa-db", "Buket Name")
	bucketlocation := flag.String("bucketlocation", "us-east-1", "Bucket Zone")
	backupdir := flag.String("backupdir", "backup", "Backup directory location")
	flag.Parse()
	//

	s3Client, err := minio.New(*endpoint, *ak, *sk, *secureendpoint)
	if err != nil {
		log.Print(err)
	}

	//We send to create the bucket, if we do not do anything
	err = s3Client.MakeBucket(*bucket, *bucketlocation)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := s3Client.BucketExists(*bucket)
		if err == nil && exists {
			log.Printf("We already own %s\n", *bucket)
		} else {
			log.Print(err)
		}
	}

	today_regex := []string{"*"}
	listfile, _ := core.FindFile(*backupdir+"/", today_regex)

	//We verify that there is something to upload and that there is no error
	if len(listfile) != 0 {
		for _, file := range listfile {
			file_name := filepath.Base(file)
			//log.Printf(file_name)
			n, err := s3Client.FPutObject(*bucket, jodaTime.Format("YYYYMMdd", time.Now())+"/"+file_name, file, minio.PutObjectOptions{ContentType: ""})
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Successfully uploaded %s of size %d\n", file_name, n)
		}
	}

	objectsCh := make(chan string)
	// Send object names that are needed to be removed to objectsCh
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
