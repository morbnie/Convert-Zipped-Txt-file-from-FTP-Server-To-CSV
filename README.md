# Convert-ZIPPED-TXT-file-from-FTP-Server-To-CSV
This small Go application will download a ZIP compressed file from an FTP server. Every TXT files in this ZIP compressed file will be extraced and saved as CSV files.

This application can be used for downloading ZIP compressed log files in TXT format from on an FTP Server. The TXT files will be extracted and saved as CSV files. These CSV files are then ready to be sent to visualization tool such as Tableau or Kibana/Elasticsearch (Elastic Stack). 


## How to use the application
* Set the following values in the main.go file

| Variable Name  | Description | Example value | 
| ------------- | ------------- | ------------- |
| ftpHost  | The value should include the hostname of the FTP server including the port.  | ftp.example.com:21  |
| ftpUser  | The username of the user that has read access to the file.  | myuser |
| ftpPassword  | The password of the user.  | mypasswprd  |
| directory  | The directory path of the user to where the file is located. The value can either be set as a varible in the main.go file or as a Command-line flag with the use of **-directory**.  | folder1/folder2  |
| filename  | The name of the file that should be downloaded. If the application should run everyday and download a daily file in a format like "Logfile-30-01-2020.zip", use currentTime to set the day, month and year values.<br><br>Use currentTime.Format("02") for day, currentTime.Format("01") for month and currentTime.Format("06") for year. For example:<br> `Logfile"+currentTime.Format("02")+"-"+currentTime.Format("01")+"-"+currentTime.Format("06")+".zip`<br><br>The value can either be set as a varible in the main.go file or as a Command-line flag with the use of **-filename**.  | File.zip  |

* Run the go file, and change set the flags if necessary.
`go run main.go -directory=folder1/folder2 -filename=File.zip` or `go run main.go -directory folder1/folder2 -filename File.zip`
* The CSV file will be saved in the output folder
