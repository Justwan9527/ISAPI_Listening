package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 获取大小的接口
type Sizer interface {
	Size() int64
}

type EventInfo struct {
	XMLName          xml.Name `xml:"EventNotificationAlert"`
	LicensePlateInfo ANPRInfo `xml:"ANPR"`
}

type ANPRInfo struct {
	XMLName      xml.Name `xml:"ANPR"`
	LicensePlate string   `xml:"licensePlate"`
}

var EventCount int

var EventCount4ANPR int

func parseMultipartFormFile(r *http.Request, formFiles map[string][]*multipart.FileHeader) {

	var strLicensePlate string
	var bHasLicensePlate bool

	var strEventCount string = strconv.Itoa(EventCount)

	//save background picture data
	var b bytes.Buffer
	var bHasBackgroundPicture bool

	for formName := range formFiles {
		// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
		// FormFile returns the first file for the provided form key
		formFile, fileHeader, _ := r.FormFile(formName)

		//log.Printf("File formname: %s, filename: %s, file length: %d\n", formName, formFileHeader.Filename, formFileHeader.Size)

		//parse license plate number
		log.Println("fileHeader:", fileHeader.Header)
		if strings.Contains(formName, "anpr") {
			log.Println("anpr xml")

			var struEventData EventInfo

			var b bytes.Buffer
			_, _ = io.Copy(&b, formFile)
			err := xml.Unmarshal(b.Bytes(), &struEventData)
			if err != nil {
				log.Println("Unmarshal fail:", err, "XML", b.String())
			}

			bHasLicensePlate = true
			strLicensePlate = struEventData.LicensePlateInfo.LicensePlate
			EventCount4ANPR++
			log.Println("Parse license plate:", strLicensePlate)
		}

		//only save background picture
		if strings.Contains(formName, "detection") {
			log.Println("detection picture")

			//copy file data
			_, _ = io.Copy(&b, formFile)
			bHasBackgroundPicture = true

		}

	}

	//wirte the background picture data to local file
	if bHasLicensePlate {
		strEventCount = strconv.Itoa(EventCount4ANPR) + "_" + strLicensePlate
	} else {
		strEventCount = strconv.Itoa(EventCount4ANPR)
	}

	if bHasBackgroundPicture {

		f, err := os.Create(strEventCount + ".jpg")

		if err != nil {
			log.Println("Create file fail, err:", err)
		} else {

			iWritelen, err := f.Write(b.Bytes())
			if err != nil {
				log.Println("Write to file fail, err:", err)
			} else {
				log.Println("Write:", iWritelen, "to file")

				f.Close()
			}

			log.Println("Write:", b.Len(), "detection picture to file")
			//fmt.Println("Write:", b.Len(), "to file")
		}

	} else {
		log.Println("!!!!!!!!!!!!!!!!!!!", strEventCount, "!!!!!!!!!!!!!!!!!!!")
	}

}

func parseMultipartForm(r *http.Request, formFiles *multipart.Form) {

	//var strLicensePlate string
	//var bHasLicensePlate bool

	var strEventCount string = strconv.Itoa(EventCount)

	//save background picture data
	//var b bytes.Buffer
	//var bHasBackgroundPicture bool

	//EventCount++

	//parse json/xml
	for formName, formContent := range formFiles.Value {

		fmt.Println("key:", formName)
		fmt.Println("value:", formContent)

		//	var b bytes.Buffer
		//	_, _ = io.Copy(&b, formFile)
		//err := xml.Unmarshal(b.Bytes(), &struEventData)
		/*
			if err != nil {
				log.Println("Unmarshal fail:", err, "XML", string(b.Bytes()))
			}
		*/

		//bHasLicensePlate = true
		//strLicensePlate = struEventData.LicensePlateInfo.LicensePlate

		var fileName string = strEventCount + "." + formName + ".json"

		f, err := os.Create(fileName)
		if err != nil {
			log.Println("Create file fail, err:", err)
		} else {

			iWritelen, err := f.Write([]byte(formContent[0]))
			if err != nil {
				log.Println("Write to file fail, err:", err)
			} else {
				log.Println("Write:", iWritelen, "to file")

				f.Close()
			}

			log.Println("Write:", len(formContent[0]), fileName, "json/xml to file")
			//fmt.Println("Write:", b.Len(), "to file")
		}
	}

	//parse file
	for formName := range formFiles.File {
		// func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
		// FormFile returns the first file for the provided form key
		formFile, fileHeader, _ := r.FormFile(formName)

		//log.Printf("File formname: %s, filename: %s, file length: %d\n", formName, formFileHeader.Filename, formFileHeader.Size)

		//parse license plate number
		log.Println("fileHeader:", fileHeader.Header["Content-Type"])
		contentType := fileHeader.Header["Content-Type"][0]
		log.Println("Content-Type:", contentType)

		if strings.Contains(contentType, "image/jpeg") {

			var b bytes.Buffer
			_, _ = io.Copy(&b, formFile)
			//err := xml.Unmarshal(b.Bytes(), &struEventData)
			/*
				if err != nil {
					log.Println("Unmarshal fail:", err, "XML", string(b.Bytes()))
				}
			*/

			//bHasLicensePlate = true
			//strLicensePlate = struEventData.LicensePlateInfo.LicensePlate
			var fileName string = strEventCount + "." + formName + ".jpg"

			f, err := os.Create(fileName)
			if err != nil {
				log.Println("Create file fail, err:", err)
			} else {

				iWritelen, err := f.Write(b.Bytes())
				if err != nil {
					log.Println("Write to file fail, err:", err)
				} else {
					log.Println("Write:", iWritelen, "to file")

					f.Close()
				}

				log.Println("Write:", b.Len(), fileName, " to file")
				//fmt.Println("Write:", b.Len(), "to file")
			}

		} else if strings.Contains(contentType, "text/json") || strings.Contains(contentType, "application/json") {
			log.Println("anpr xml")

			//var struEventData EventInfo

			var b bytes.Buffer
			_, _ = io.Copy(&b, formFile)
			//err := xml.Unmarshal(b.Bytes(), &struEventData)
			/*
				if err != nil {
					log.Println("Unmarshal fail:", err, "XML", string(b.Bytes()))
				}
			*/

			//bHasLicensePlate = true
			//strLicensePlate = struEventData.LicensePlateInfo.LicensePlate

			var fileName string = strEventCount + "." + formName + ".json"
			f, err := os.Create(fileName)
			if err != nil {
				log.Println("Create file fail, err:", err)
			} else {

				iWritelen, err := f.Write(b.Bytes())
				if err != nil {
					log.Println("Write to file fail, err:", err)
				} else {
					log.Println("Write:", iWritelen, "to file")

					f.Close()
				}

				log.Println("Write:", b.Len(), "json to file")
				//fmt.Println("Write:", b.Len(), "to file")
			}
		} else if strings.Contains(contentType, "text/xml") || strings.Contains(contentType, "application/xml") {
			log.Println("anpr xml")

			//var struEventData EventInfo

			var b bytes.Buffer
			_, _ = io.Copy(&b, formFile)
			//err := xml.Unmarshal(b.Bytes(), &struEventData)
			/*
				if err != nil {
					log.Println("Unmarshal fail:", err, "XML", string(b.Bytes()))
				}
			*/

			//bHasLicensePlate = true
			//strLicensePlate = struEventData.LicensePlateInfo.LicensePlate

			var fileName string = strEventCount + "." + formName + ".xml"
			f, err := os.Create(fileName)
			if err != nil {
				log.Println("Create file fail, err:", err)
			} else {

				iWritelen, err := f.Write(b.Bytes())
				if err != nil {
					log.Println("Write to file fail, err:", err)
				} else {
					log.Println("Write:", iWritelen, "to file")

					f.Close()
				}

				log.Println("Write:", b.Len(), "xml to file")
				//fmt.Println("Write:", b.Len(), "to file")
			}
		} else {
			log.Println("Unknown http Content-Type")
		}

	}

}

// hello world, the web server
func HelloServer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("receive http request from client")
	if r.Method == "POST" {

		err := r.ParseMultipartForm(20 * 1024 * 1024)
		if err != nil {
			log.Fatal("Parse multipart form fail, err", err)
		} else {

			//save file to local disk.
			EventCount++
			fmt.Println("File:", r.MultipartForm.File)
			fmt.Println("Value:", r.MultipartForm.Value)
			parseMultipartForm(r, r.MultipartForm)

		}

	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Read http body fail, err", err)
		} else {
			fmt.Println("requst data:", string(body))
		}

	}

	defer r.Body.Close()

	// 上传页面
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
<form enctype="multipart/form-data" action="/test" method="POST">
    Send this file: <input name="file" type="file" />
    <input type="submit" value="Send File" />
</form>
`
	io.WriteString(w, html)
}

// hello world, the web server
func LicensePlate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("receive http request from client")
	if r.Method == "POST" {

		err := r.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			log.Fatal("Parse multipart form fail, err", err)
		} else {
			log.Println("parse multipart form:", r.MultipartForm)

			//save file to local disk.
			if r.MultipartForm.File != nil {

				EventCount++
				parseMultipartFormFile(r, r.MultipartForm.File)
			}

		}

	}

	defer r.Body.Close()

	// 上传页面
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
<form enctype="multipart/form-data" action="/test" method="POST">
    Send this file: <input name="file" type="file" />
    <input type="submit" value="Send File" />
</form>
`
	io.WriteString(w, html)
}

func main() {

	var args []string = os.Args

	var strIP string
	var strURL string

	if args == nil || len(args) < 2 {

		fmt.Println("Please input port and url, like 8080 /test")
		return
	}

	if strings.Contains(args[1], ":") {
		strIP = args[1]

	} else {
		strIP = ":" + args[1]
	}

	if !strings.HasPrefix(args[2], "/") {
		strURL = "/" + args[2]

	} else {
		strURL = args[2]
	}

	log.Println("IP:", strIP, "URL:", strURL)

	http.HandleFunc(strURL, HelloServer)
	http.HandleFunc("/license", LicensePlate)
	err := http.ListenAndServe(strIP, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
