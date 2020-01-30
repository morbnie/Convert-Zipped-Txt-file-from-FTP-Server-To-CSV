package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

var (
	currentTime = time.Now()
	ftpHost     = "HOSTNAME:PORT"
	ftpUser     = "USERNAME"
	ftpPassword = "PASSWORD"
	ftpDir      = flag.String("directory", "INSERT/DIRECTORY/HERE", "Directory to where the Zip file is located")
	ftpDownload = flag.String("filename", "LOGFILE-"+currentTime.Format("02")+"-"+currentTime.Format("01")+"-"+currentTime.Format("06")+".zip", "The name of the Zip file")
)

func main() {

	fmt.Println("Job started...")

	//Get "ftpDownload" file from the FTP server
	getFileFromFTP()

	//Unzip the zip file and save the txt file as .csv
	unzipAndSave()

	//Delete the temp.zip file
	firstfile := "temp.zip"
	err := os.Remove(firstfile)
	if err != nil {
		log.Fatal(err)
	}
}

func getFileFromFTP() {

	flag.Parse()

	ftpFile, err := ftp.Dial(ftpHost, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = ftpFile.Login(ftpUser, ftpPassword)
	if err != nil {
		log.Fatal(err)
	}

	ftpFile.ChangeDir(*ftpDir)

	res, err := ftpFile.Retr(*ftpDownload)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close()

	outFile, err := os.Create("temp.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res)
	if err != nil {
		log.Fatal(err)
	}

}

func unzipAndSave() {
	//unzip file
	files, err := unzip("temp.zip", "output")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done. The following txt files has been unzipped and converted to CSV:\n" + strings.Join(files, "\n"))

}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {

		if filepath.Ext(f.Name) == ".txt" {

			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name+".csv")
			filenames = append(filenames, fpath)

			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return filenames, fmt.Errorf("%s: illegal file path", fpath)
			}

			if f.FileInfo().IsDir() {
				// Make Folder
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				log.Fatal(err)
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatal(err)
			}

			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			defer outFile.Close()
			defer rc.Close()

			if err != nil {
				return filenames, err
			}
		}
	}
	return filenames, nil

}
