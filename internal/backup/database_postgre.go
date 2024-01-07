package backup

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"cubeflow/pkg/config"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func BackupDB(pgDatabase string) (string, error) {

	isBucketExist, err := DoesBucketExist(config.Variable.GCP.BucketName, config.Variable.GCP.CredentialsPath)
	if !isBucketExist {
		if err != nil {
			return config.Variable.GCP.BucketName, err
		}
		err := CreateBucket(config.Variable.GCP.ProjectID, config.Variable.GCP.BucketName, config.Variable.GCP.CredentialsPath)
		if err != nil {
			return config.Variable.GCP.BucketName, err
		}
	}

	dumpPath, dumpFile, err := CreateDump(config.Variable.Database.Host, config.Variable.Database.Port, config.Variable.Database.Username, pgDatabase, config.Variable.Environment)
	if err != nil {
		return config.Variable.GCP.BucketName, err
	}

	err = UploadToGCS(dumpPath, config.Variable.GCP.CredentialsPath, config.Variable.GCP.BucketName, dumpFile)
	if err != nil {
		return config.Variable.GCP.BucketName, err
	}

	err = CleanDump(dumpPath)
	if err != nil {
		return config.Variable.GCP.BucketName, err
	}
	return config.Variable.GCP.BucketName, nil
}

func DoesBucketExist(bucketName string, gcsCredentialsFile string) (bool, error) {
	// Set up a Google Cloud Storage client.
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(gcsCredentialsFile))
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return false, err
	}
	defer client.Close()

	// Check if the bucket exists.
	_, err = client.Bucket(bucketName).Attrs(context.Background())
	if err == storage.ErrBucketNotExist {
		// Bucket does not exist.
		log.Printf("This bucket does not exist: %v", bucketName)
		return false, nil
	} else if err != nil {
		// An error occurred other than the bucket not existing.
		log.Printf("Error while getting bucket check.: %v", err)
		return false, err
	}

	// Bucket exists.
	log.Printf("%v bucket is existed", bucketName)
	return true, nil
}

func CreateBucket(projectID string, bucketName string, gcsCredentialsFile string) error {
	// Set up a Google Cloud Storage client.
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(gcsCredentialsFile))
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return err
	}
	defer client.Close()

	// Create a new bucket.
	if err := client.Bucket(bucketName).Create(context.Background(), projectID, nil); err != nil {
		log.Printf("Failed to create bucket: %v", err)
		return err
	}

	fmt.Printf("Bucket '%s' created.\n", bucketName)
	return nil
}

func CreateDump(pgHost string, pgPort string, pgUser string, pgDatabase string, env string) (string, string, error) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02_15-04-05")

	// Path for the temporary dump file
	dumpFile := pgDatabase + "_" + env + "_" + formattedTime + ".sql"
	dumpPath := "/tmp/" + dumpFile

	// Construct the pg_dump command
	cmd := exec.Command(
		"pg_dump",
		"-h", pgHost,
		"-p", pgPort,
		"-U", pgUser,
		"-d", pgDatabase,
		"--file", dumpPath,
	)

	// Run the pg_dump command
	err := cmd.Run()
	if err != nil {
		println("Creating dump error ", err)
		return dumpPath, dumpFile, err
	}

	return dumpPath, dumpFile, nil
}

func UploadToGCS(dumpFilePath string, gcsCredentialsFile string, gcsBucketName string, dumpFile string) error {
	// Upload the dump to Google Cloud Storage
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsCredentialsFile))
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return err
	}

	bucket := client.Bucket(gcsBucketName)
	obj := bucket.Object(dumpFile)
	wc := obj.NewWriter(ctx)

	defer wc.Close()

	// Upload the dump file
	OpenedDumpFile, err := os.Open(dumpFilePath)
	if err != nil {
		log.Printf("Failed to open dump file: %v", err)
		return err
	}
	defer OpenedDumpFile.Close()

	if _, err = io.Copy(wc, OpenedDumpFile); err != nil {
		log.Printf("Failed to copy dump file to GCS: %v", err)
		return err
	}

	log.Printf("Dump uploaded to GCS: gs://%s/%s", gcsBucketName, dumpFile)
	return nil
}

func CleanDump(dumpPath string) error {
	cmd := exec.Command(
		"rm", dumpPath,
	)
	err := cmd.Run()
	if err != nil {
		log.Println("Error when cleaning backup dump: ", err)
		return err
	}
	return nil
}
