package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Params interface {
	Params() map[string]string
}

func unionMaps(a, b map[string]string) map[string]string {
	c := make(map[string]string)
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

type TaskInfo struct {
	/**
	 * Task identifier
	 */
	TaskId uuid.UUID

	/**
	 * Task creation time
	 */
	RegistrationTime time.Time

	/**
	 * Last Task modification time
	 */
	StatusChangeTime time.Time

	/**
	 * The task can have one of the following statuses
	 */
	Status TaskStatus

	/**
	 * Description of the processing error. Specified only with
	 * ProcessingFailed Task status
	 */
	Error Error

	/**
	 * Number of files added to a Task
	 */
	FilesCount int

	/**
	 * Recommended delay before request for new Task Status in milliseconds
	 */
	RequestStatusDelay int

	/**
	 * The hyperlink collection with recognition results.
	 * The links have limited lifetime. If you want to download the
	 * result after the time has passed, call the
	 * {@link IOcrClient#getTaskStatusAsync(UUID)}
	 * or {@link IOcrClient#listTasksAsync(TasksListingParams)} method
	 * to obtain the new hyperlink collection.
	 */
	ResultUrls []string

	/**
	 * Task description specified when the Task is created
	 */
	Description string
}

func (ti TaskInfo) IsInProcess() bool {
	switch ti.Status {
	case TaskStatusQueued:
	case TaskStatusInProgress:
		return true
	default:
		return false
	}
	return false
}

type ImageProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default is {@link ExportFormat#Rtf}. Specifies the export format.
	 */
	ExportFormats []ExportFormat `json:"exportFormat"`

	/**
	 * Optional. Default is {@link ProcessingProfile#DocumentConversion}. Specifies a profile with predefined processing settings.
	 */
	Profile ProcessingProfile

	/**
	 * Optional. Default is {@link TextType#Normal}. Specifies the type of the text on a page.
	 */
	TextTypes []TextType `json:"textType"`

	/**
	 * Optional. Default is {@link ImageSource#Auto}. Specifies the source of the image.
	 */
	ImageSource ImageSource

	/**
	 * Optional. Default "true". Specifies whether the orientation of the image should be automatically detected and corrected.
	 * <ul>
	 *   <li><b>true</b></li> The page orientation is automatically detected, and if it differs from normal the image is rotated.
	 *   <li><b>false</b></li> The page orientation detection and correction is not performed.
	 * </ul>
	 */
	CorrectOrientation bool

	/**
	 * Optional. Default "true". Specifies whether the skew of the image should be automatically detected and corrected.
	 */
	CorrectSkew bool

	/**
	 * Optional. Default "English". Specifies recognition language of the document.
	 * This parameter can contain several language names separated with commas, for example
	 * "English,French,German".
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/recognition-languages/"/>
	 */
	Language string

	/**
	 * Optional. Default is {@link WriteTags#Auto}. Specifies whether the result must be written as tagged PDF.
	 * his parameter can be used only if the {@link ExportFormat} parameter contains one of the
	 * values for export to PDF.
	 */
	WriteTags WriteTags `json:"pdf:writeTags"`

	/**
	 * Optional. Default "false". Specifies whether the variants of characters recognition
	 * should be written to an output file in XML format. This parameter can be used only
	 * if the {@link ExportFormat} parameter contains xml or xmlForCorrectedImage value.
	 */
	WriteRecognitionVariants bool `json:"xml:writeRecognitionVariants"`

	/**
	 * Optional. Default "false". Specifies whether the paragraph and character styles
	 * should be written to an output file in XML format. This parameter can be
	 * used only if the {@link ExportFormat} parameter contains xml or
	 * xmlForCorrectedImage value.
	 */
	WriteFormatting bool `json:"xml:writeFormatting"`

	/**
	 * Optional. Default "true" for xml export format and "false" in other cases.
	 * Specifies whether barcodes must be detected on the image, recognized and exported
	 * to the result file.
	 */
	ReadBarcodes bool
}

func (p *ImageProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.Language != "" {
		params["language"] = p.Language
	}
	if p.Profile != "" {
		params["profile"] = string(p.Profile)
	}
	if len(p.TextTypes) != 0 {
		s := make([]string, len(p.TextTypes))
		for i, v := range p.TextTypes {
			s[i] = string(v)
		}
		params["textType"] = strings.Join(s, ",")
	}
	if p.ImageSource != "" {
		params["imageSource"] = string(p.ImageSource)
	}
	params["correctOrientation"] = strconv.FormatBool(p.CorrectOrientation)
	params["correctSkew"] = strconv.FormatBool(p.CorrectSkew)
	params["readBarcodes"] = strconv.FormatBool(p.ReadBarcodes)
	if len(p.ExportFormats) != 0 {
		s := make([]string, len(p.ExportFormats))
		for i, v := range p.ExportFormats {
			s[i] = string(v)
		}
		params["exportFormat"] = strings.Join(s, ",")
	}
	params["xml:writeFormatting"] = strconv.FormatBool(p.WriteFormatting)
	params["xml:writeRecognitionVariants"] = strconv.FormatBool(p.WriteRecognitionVariants)
	if p.WriteTags != "" {
		params["pdf:writeTags"] = string(p.WriteTags)
	}
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Image Submitting request
 */
type ImageSubmittingParams struct {
	TaskInfo

	/**
	 * Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string
}

func (p *ImageSubmittingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.TaskId.String() != "00000000-0000-0000-0000-000000000000" {
		params["taskId"] = p.TaskId.String()
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Document Processing request
 */
type DocumentProcessingParams struct {
	TaskInfo
	/**
	 * Required. Specifies the identifier of the task. If the task with the
	 * specified identifier does not exist or has been deleted, an error is
	 * returned.
	 */
	TaskId uuid.UUID

	/**
	 * Optional. Contains the description of the processing task. Cannot
	 * contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default is {@link ExportFormat#Rtf}. Specifies the export format.
	 */
	ExportFormats []ExportFormat `json:"exportFormat"`

	/**
	 *  Optional. Default is {@link ProcessingProfile#DocumentConversion}. Specifies a profile with predefined processing settings.
	 */
	Profile ProcessingProfile

	/**
	 * Optional. Default is {@link TextType#Normal}. Specifies the type of the text on a page.
	 */
	TextTypes []TextType `json:"textType"`

	/**
	 * Optional. Default "English". Specifies recognition language of the document.
	 * This parameter can contain several language names separated with commas, for example
	 * "English,French,German".
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/recognition-languages/"/>
	 */
	Language string

	/**
	 * Optional. Default is {@link ImageSource#Auto}. Specifies the source of the image.
	 */
	ImageSource ImageSource

	/**
	 * Optional. Default "true". Specifies whether the orientation of the image should be automatically detected and corrected.
	 * <ul>
	 *   <li><b>true</b></li> The page orientation is automatically detected, and if it differs from normal the image is rotated.
	 *   <li><b>false</b></li> The page orientation detection and correction is not performed.
	 * </ul>
	 */
	CorrectOrientation bool

	/**
	 * Optional. Default "true". Specifies whether the skew of the image should be automatically detected and corrected.
	 */
	CorrectSkew bool

	/**
	 * Optional. Default is {@link WriteTags#Auto}. Specifies whether the result must be written as tagged PDF.
	 * This parameter can be used only if the {@link ExportFormat} parameter contains one of the
	 * values for export to PDF.
	 */
	WriteTags WriteTags `json:"pdf:writeTags"`

	/**
	 * Optional. Default "false". Specifies whether the variants of characters recognition
	 * should be written to an output file in XML format. This parameter can be used only
	 * if the {@link ExportFormat} parameter contains xml or xmlForCorrectedImage value.
	 */
	WriteRecognitionVariants bool `json:"xml:writeRecognitionVariants"`

	/**
	 * Optional. Default "false". Specifies whether the paragraph and character styles
	 * should be written to an output file in XML format. This parameter can be
	 * used only if the {@link ExportFormat} parameter contains xml or
	 * xmlForCorrectedImage value.
	 */
	WriteFormatting bool `json:"xml:writeFormatting"`

	/**
	 * Optional. Default "true" for xml export format and "false" in other cases.
	 * Specifies whether barcodes must be detected on the image, recognized and exported
	 * to the result file.
	 */
	ReadBarcodes bool
}

func (p *DocumentProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	params["taskId"] = p.TaskId.String()
	if p.Language != "" {
		params["language"] = p.Language
	}
	if p.Profile != "" {
		params["profile"] = string(p.Profile)
	}
	if len(p.TextTypes) != 0 {
		s := make([]string, len(p.TextTypes))
		for i, v := range p.TextTypes {
			s[i] = string(v)
		}
		params["textType"] = strings.Join(s, ",")
	}
	if p.ImageSource != "" {
		params["imageSource"] = string(p.ImageSource)
	}
	params["correctOrientation"] = strconv.FormatBool(p.CorrectOrientation)
	params["correctSkew"] = strconv.FormatBool(p.CorrectSkew)
	params["readBarcodes"] = strconv.FormatBool(p.ReadBarcodes)
	if len(p.ExportFormats) != 0 {
		s := make([]string, len(p.ExportFormats))
		for i, v := range p.ExportFormats {
			s[i] = string(v)
		}
		params["exportFormat"] = strings.Join(s, ",")
	}
	params["xml:writeFormatting"] = strconv.FormatBool(p.WriteFormatting)
	params["xml:writeRecognitionVariants"] = strconv.FormatBool(p.WriteRecognitionVariants)
	if p.WriteTags != "" {
		params["pdf:writeTags"] = string(p.WriteTags)
	}
	if p.Description != "" {
		params["description"] = p.Description
	}
	return params
}

/**
 * Parameters for Business Card Processing request
 */
type BusinessCardProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default is <see cref="Enums.BusinessCardExportFormat.VCard"/>. Specifies the export format.
	 */
	ExportFormat BusinessCardExportFormat `json:"exportFormat"`

	/**
	 * Optional. Default "English". Specifies recognition language of the document.
	 * This parameter can contain several language names separated with commas, for example
	 * "English,French,German".
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/recognition-languages/"/>
	 */
	Language string

	/**
	 * Optional. Default is {@link ImageSource#Auto}. Specifies the source of the image.
	 */
	ImageSource ImageSource

	/**
	 * Optional. Default "true". Specifies whether the orientation of the image should be automatically detected and corrected.
	 * <ul>
	 *   <li><b>true</b></li> The page orientation is automatically detected, and if it differs from normal the image is rotated.
	 *   <li><b>false</b></li> The page orientation detection and correction is not performed.
	 * </ul>
	 */
	CorrectOrientation bool

	/**
	 * Optional. Default "true". Specifies whether the skew of the image should be automatically detected and corrected.
	 */
	CorrectSkew bool

	/**
	 * Optional. Default "false". Specifies whether the additional information
	 * on the recognized characters (e.g. whether the character is recognized
	 * uncertainly) should be written to an output file in XML format. This
	 * parameter can be used only if the {@link #exportFormats} parameter
	 * is set to {@link ExportFormat#Xml}.
	 */
	WriteExtendedCharacterInfo bool `json:"xml:writeExtendedCharacterInfo"`

	/**
	 * Optional. Default "false". Specifies whether the field components should
	 * be written to an output file in XML format. For example, for the Name
	 * field the components can include first name and last name, returned separately. This
	 * parameter can be used only if the {@link #exportFormats} parameter
	 * is set to {@link ExportFormat#Xml}.
	 */
	WriteFieldComponents bool `json:"xml:writeFieldComponents"`
}

func (p *BusinessCardProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.Language != "" {
		params["language"] = p.Language
	}
	if p.ImageSource != "" {
		params["imageSource"] = string(p.ImageSource)
	}
	params["correctOrientation"] = strconv.FormatBool(p.CorrectOrientation)
	params["correctSkew"] = strconv.FormatBool(p.CorrectSkew)
	if p.ExportFormat != "" {
		params["exportFormat"] = string(p.ExportFormat)
	}
	params["xml:writeExtendedCharacterInfo"] = strconv.FormatBool(p.WriteExtendedCharacterInfo)
	params["xml:writeFieldComponents"] = strconv.FormatBool(p.WriteFieldComponents)
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Text Field Processing request.
 */
type TextFieldProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default "-1,-1,-1,-1". Specifies the region of the text field on the image.
	 * The coordinates of the region are measured in pixels relative to the left top corner of the image and
	 * are specified in the following order: left, top, right, bottom. By default, the region of the whole image is used.
	 */
	Region string

	/**
	 * Optional. Default "English". Specifies recognition language of the document.
	 * This parameter can contain several language names separated with commas, for example
	 * "English,French,German".
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/recognition-languages/"/>
	 */
	Language string

	/**
	 * Optional. Default "". Specifies the letter set, which should be used during recognition.
	 * Contains a string with the letter set characters. For example, "ABCDabcd'-.".
	 * By default, the letter set of the language, specified in the {@link #language} parameter, is used.
	 */
	LetterSet string

	/**
	 * Optional. Default "". Specifies the regular expression which defines which words are allowed in the field
	 * and which are not. By default, the set of allowed words is defined by the dictionary of the language,
	 * specified in the language parameter.
	 *
	 * <b>Note:</b> See the <a href="https://www.ocrsdk.com/documentation/specifications/regular-expressions/"> description of regular expressions</a>.
	 * Note that regular expressions do not strictly limit the set of characters of the output result,
	 * i.e. the recognized value may contain characters which are not included into the regular expression.
	 * During recognition all hypotheses of a word recognition are checked against the specified regular expression.
	 * If a given recognition variant conforms to the expression, it has higher probability of being selected
	 * as final recognition output. But if there is no variant that matches regular expression,
	 * the result will not conform to the expression. If you want to limit the set of characters, which can be recognized,
	 * the best way to do it is to use letterSet parameter.
	 */
	RegExp string

	/**
	 * Optional. Default is {@link TextType#Normal}. Specifies the type of the text on a page.
	 */
	TextTypes []TextType `json:"textType"`

	/**
	 * Optional. Default "false". Specifies whether the field contains only one text line.
	 * The value should be true, if there is one text line in the field otherwise it should be false.
	 */
	OneTextLine bool

	/**
	 * Optional. Default "false". Specifies whether the field contains only one word in each text line.
	 * The value should be true, if no text line contains more than one word (so the lines of text will be recognized
	 * as a single word) otherwise it should be false.
	 */
	OneWordPerTextLine bool

	/**
	 * Optional. Default is {@link MarkingType#SimpleText}. This property is valid only
	 * for the {@link TextType#Handprinted} recognition. Specifies the type of marking around letters
	 * (for example, underline, frame, box, etc.). By default, there is no marking around letters.
	 */
	MarkingType MarkingType

	/**
	 * Optional. Default "1". Specifies the number of character cells for the field.
	 * This property has a sense only for the field marking types(the markingType parameter) that imply splitting the text in cells.
	 * Default value for this property is 1, but you should set the appropriate value to recognize the text correctly.
	 */
	PlaceholdersCount int

	/**
	 * Optional. Default "default". Provides additional information about handprinted letters writing style.
	 */
	WritingStyle WritingStyle
}

func (p *TextFieldProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.Region != "" {
		params["region"] = p.Region
	}
	if p.Language != "" {
		params["language"] = p.Language
	}
	if p.LetterSet != "" {
		params["letterSet"] = p.LetterSet
	}
	if p.RegExp != "" {
		params["regExp"] = p.RegExp
	}
	if len(p.TextTypes) > 0 {
		s := make([]string, len(p.TextTypes))
		for i, v := range p.TextTypes {
			s[i] = string(v)
		}
		params["textType"] = strings.Join(s, ",")
	}
	params["oneTextLine"] = strconv.FormatBool(p.OneTextLine)
	params["oneWordPerTextLine"] = strconv.FormatBool(p.OneWordPerTextLine)
	if p.MarkingType != "" {
		params["markingType"] = string(p.MarkingType)
	}
	if p.PlaceholdersCount != 0 {
		params["placeholdersCount"] = strconv.Itoa(p.PlaceholdersCount)
	}
	if p.WritingStyle != "" {
		params["writingStyle"] = string(p.WritingStyle)
	}
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Barcode Field Processing
 */
type BarcodeFieldProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default "-1,-1,-1,-1". Specifies the region of the text field on the image. The coordinates of the region
	 * are measured in pixels relative to the left top corner of the image and are specified in the following order:
	 * left, top, right, bottom. By default, the region of the whole image is used.
	 */
	Region string

	/**
	 * Optional. Default is Autodetect. Specifies the type of the barcode.
	 * This parameter may also contain several barcode types.
	 */
	BarcodeTypes []BarcodeType `json:"barcodeType"`

	/**
	 * Optional. Default "false". This parameter makes sense only for {@link BarcodeType#Pdf417}
	 * and {@link BarcodeType#Aztec} barcodes, which encode some binary data.
	 * If this parameter is set to true, the binary data encoded in a barcode are saved as a sequence of hexadecimal
	 * values for corresponding bytes.
	 */
	ContainsBinaryData bool
}

func (p *BarcodeFieldProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.Region != "" {
		params["region"] = p.Region
	}
	if len(p.BarcodeTypes) > 0 {
		s := make([]string, len(p.BarcodeTypes))
		for i, v := range p.BarcodeTypes {
			s[i] = string(v)
		}
		params["barcodeType"] = strings.Join(s, ",")
	}
	params["containsBinaryData"] = strconv.FormatBool(p.ContainsBinaryData)
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Checkmark Field Processing request
 */
type CheckmarkFieldProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default "-1,-1,-1,-1". Specifies the region of the text field on the image.
	 * The coordinates of the region are measured in pixels relative to the left top corner of the image and
	 * are specified in the following order: left, top, right, bottom. By default, the region of the whole image is used.
	 */
	Region string

	/**
	 * Optional. Default is {@link CheckmarkType#Empty}. Specifies the type of the checkmark.
	 */
	CheckmarkType CheckmarkType

	/**
	 * Optional. Default "false". This property set to true means that checkmark block can be selected and then corrected.
	 */
	CorrectionAllowed bool
}

func (p *CheckmarkFieldProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	if p.Region != "" {
		params["region"] = p.Region
	}
	if p.CheckmarkType != "" {
		params["checkmarkType"] = string(p.CheckmarkType)
	}
	params["correctionAllowed"] = strconv.FormatBool(p.CorrectionAllowed)
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Fields Processing request
 */
type FieldsProcessingParams struct {
	TaskInfo
	/**
	 * Required. Specifies the identifier of the task. If the task with the specified identifier does not exist or
	 * has been deleted, an error is returned.
	 */
	TaskId uuid.UUID

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default "false". Specifies whether the recognition variants should be written to the result.
	 * If you set this value to true, additional recognition variants (charRecVariants) appear in the XML result file.
	 */
	WriteRecognitionVariants bool
}

func (p *FieldsProcessingParams) Params() map[string]string {
	params := make(map[string]string)
	params["taskId"] = p.TaskId.String()
	params["writeRecognitionVariants"] = strconv.FormatBool(p.WriteRecognitionVariants)
	if p.Description != "" {
		params["description"] = p.Description
	}
	return params
}

/**
 * Parameters for MRZ Processing request
 */
type MrzProcessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 * Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string
}

func (p *MrzProcessingParams) Params() map[string]string {
	return make(map[string]string)
}

/**
 * Parameters for Receipt Processing request
 */
type ReceiptProccessingParams struct {
	TaskInfo
	/**
	 * Optional. Contains a password for accessing password-protected images in PDF format.
	 */
	PdfPassword string

	/**
	 *  Optional. Contains the description of the processing task. Cannot contain more than 255 characters.
	 */
	Description string

	/**
	 * Optional. Default is {@link ImageSource#Auto}. Specifies the source of the image.
	 */
	ImageSource ImageSource

	/**
	 * Optional. Default "true". Specifies whether the orientation of the image should be automatically detected and corrected.
	 * <ul>
	 *   <li><b>true</b></li> The page orientation is automatically detected, and if it differs from normal the image is rotated.
	 *   <li><b>false</b></li> The page orientation detection and correction is not performed.
	 * </ul>
	 */
	CorrectOrientation bool

	/**
	 * Optional. Default "true". Specifies whether the skew of the image should be automatically detected and corrected.
	 */
	CorrectSkew bool

	/**
	 * Optional. Default is {@link ReceiptRecognizingCountry#Usa}.
	 * Important! The technology fully supports the receipts from USA and France, other countries
	 * are currently supported in beta mode. Specifies the country where the receipt was printed.
	 * This parameter can contain several names of countries.
	 */
	Countries []ReceiptRecognizingCountry `json:"country"`

	/**
	 * Optional. Default "false". Specifies whether the additional information on the recognized characters
	 * (e.g. whether the character is recognized uncertainly) should be written to an output file in XML format.
	 */
	WriteExtendedCharacterInfo bool `json:"xml:writeExtendedCharacterInfo"`

	/**
	 * Optional. Default is {@link FieldRegionExportMode#DoNotExport}. Specifies if the coordinates of field regions
	 * should be saved to the resulting XML file, and how the coordinates should be specified:
	 * on the original or on the corrected image.
	 */
	FieldRegionExportMode FieldRegionExportMode `json:"xml:fieldRegionExportMode"`
}

func (p *ReceiptProccessingParams) Params() map[string]string {
	params := make(map[string]string)
	if len(p.Countries) > 0 {
		s := make([]string, len(p.Countries))
		for i, v := range p.Countries {
			s[i] = string(v)
		}
		params["country"] = strings.Join(s, ",")
	}
	if p.ImageSource != "" {
		params["imageSource"] = string(p.ImageSource)
	}
	params["correctOrientation"] = strconv.FormatBool(p.CorrectOrientation)
	params["correctSkew"] = strconv.FormatBool(p.CorrectSkew)
	params["xml:writeExtendedCharacterInfo"] = strconv.FormatBool(p.WriteExtendedCharacterInfo)
	if p.FieldRegionExportMode != "" {
		params["xml:fieldRegionExportMode"] = string(p.FieldRegionExportMode)
	}
	if p.Description != "" {
		params["description"] = p.Description
	}
	if p.PdfPassword != "" {
		params["pdfPassword"] = p.PdfPassword
	}
	return params
}

/**
 * Parameters for Tasks Listing request
 */
type TasksListingParams struct {
	TaskInfo
	/**
	 * Optional. Default is the current date minus 7 days. Specifies the date to list tasks from.
	 */
	FromDate time.Time

	/**
	 * Optional. Default is the current date. Specifies the date to list tasks to.
	 */
	ToDate time.Time

	/**
	 * Optional. Default is "false". Specifies if the tasks that have already been deleted must be excluded from the listing.
	 */
	ExcludeDeleted bool
}

func (p *TasksListingParams) Params() map[string]string {
	params := make(map[string]string)
	var zero time.Time
	if p.FromDate != zero {
		//yyyy-mm-ddThh:mm:ssZ
		params["fromDate"] = p.FromDate.Format("2006-01-02T15:04:05-07")
	}
	if p.ToDate != zero {
		//yyyy-mm-ddThh:mm:ssZ
		params["toDate"] = p.ToDate.Format("2006-01-02T15:04:05-07")
	}
	params["excludeDeleted"] = strconv.FormatBool(p.ExcludeDeleted)
	return params
}

/**
 * Parameters for Task Deletion request
 */
type TaskDeletionParams struct {
	/**
	 * Required. Specifies the identifier of the task. If the task with the specified identifier does not exist, an error is returned.
	 */
	TaskId uuid.UUID
}

func (p *TaskDeletionParams) Params() map[string]string {
	params := make(map[string]string)
	params["taskId"] = p.TaskId.String()
	return params
}
