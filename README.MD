# Cloud OCR SDK
This Go library is intended to simplify the development of image recognition solutions with ABBYY Cloud OCR SDK API. 
This repo contains the client library enabling simple access to high-quality recognition technologies. 
For more information about the product visit [ABBYY Cloud OCR SDK website](https://www.ocrsdk.com/).

## Install
```
go get github.com/microo8/cloudsdk-client-go
```

## Usage
The following sample demonstrates how to process image with the specified parameters and export the result to TXT and DOCX 
using the ABBYY Cloud OCR SDK client library. 
Start your developement with the following steps.

Provide your authentication info and the HTTP server URL address, used for API calls (see Data processing location for details) and create `OcrClient`:
```
client := abbyysdk.NewOcrClient(host, applicationId, password)
```

Set the processing and result options (you can find all the settings 
[here](https://www.ocrsdk.com/documentation/api-reference/process-image-method/):
```
imageProcessingParams := &abbyysdk.ImageProcessingParams{
	ExportFormats: []abbyysdk.ExportFormat{abbyysdk.ExportFormatDocx, abbyysdk.ExportFormatTxt},
	Language: "English,French",
}
```

Open file as input stream:
```
f, err := os.Open(filepath)
if err != nil {
	log.Fatal(err)
}
```

Pass the parameters to a new task. 
Method `ProcessImage` of `OcrClient` get's the task info:
```
taskInfo, err := ocrClient.ProcessImage(imageProcessingParams, fileStream, filepath)
if err != nil {
	log.Fatal(err)
}
```
Waiting for the task to be completed is done with the `WaitForTask` method:
```
taskInfo, err = ocrClient.WaitForTask(taskInfo)
if err != nil {
	log.Fatal(err)
}
```

Display the URL addresses of the result documents:
```
for _, resultUrl := range taskInfo.ResultUrls {
    fmt.Println(resultUrl)
}
```

See more examples in tests `ocr_client_test.go`.

## License
Copyright © 2019 microo8

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
