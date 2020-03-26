package models

/**
 * The task can have one of the following statuses
 */
type TaskStatus string

const (
	/**
	 * The task has been registered in the system, but has not yet been passed for processing.
	 */
	TaskStatusSubmitted TaskStatus = "Submitted"

	/**
	 * The task has been placed in the processing queue and is waiting to be processed.
	 */
	TaskStatusQueued TaskStatus = "Queued"

	/**
	 * The task is being processed.
	 */
	TaskStatusInProgress TaskStatus = "InProgress"

	/**
	 * The task has been processed successfully. For a task with this status, the URL for
	 * downloading the result of processing is available in the response.
	 */
	TaskStatusCompleted TaskStatus = "Completed"

	/**
	 * The task has not been processed because an error occurred. You can find the description of the
	 * error in the XML response.
	 *
	 * <b>Note:</b> We do not recommend sending the same file for processing repeatedly, if the first task
	 * failed. However, if you have created several tasks for the same file, and all of them have
	 * failed, the fullest and most specific error description can be received by calling
	 * {@link IOcrClient#getTaskStatusAsync(UUID)} for the first of those tasks.
	 */
	TaskStatusProcessingFailed TaskStatus = "ProcessingFailed"

	/**
	 * The task has been deleted.
	 */
	TaskStatusDeleted TaskStatus = "Deleted"

	/**
	 * You do not have enough money on your account to process the task.
	 */
	TaskStatusNotEnoughCredits TaskStatus = "NotEnoughCredits"
)

/**
 * Format of exporting recognized text.
 * Several formats may be provided at once
 */
type ExportFormat string

const (
	/**
	 * TXT
	 */
	ExportFormatTxt ExportFormat = "txt"

	/**
	 * TXT. The exported file contains the text that was saved according to the order of the original blocks.
	 */
	ExportFormatTxtUnstructured ExportFormat = "txtUnstructured"

	/**
	 * RTF
	 */
	ExportFormatRtf ExportFormat = "rtf"

	/**
	 * DOCX
	 */
	ExportFormatDocx ExportFormat = "docx"

	/**
	 * XLSX
	 */
	ExportFormatXlsx ExportFormat = "xlsx"

	/**
	 * PPTX
	 */
	ExportFormatPptx ExportFormat = "pptx"

	/**
	 * PDF. The entire image is saved as a picture, with recognized text put under the image.
	 */
	ExportFormatPdfSearchable ExportFormat = "pdfSearchable"

	/**
	 * PDF. The recognized text is saved as text, and the pictures are embedded as images.
	 */
	ExportFormatPdfTextAndImages ExportFormat = "pdfTextAndImages"

	/**
	 * PDF/A-1b. The file is saved in PDF/A-1b-compliant format, with the entire image saved
	 * as a picture and recognized text put under it.
	 */
	ExportFormatPdfA ExportFormat = "pdfA"

	/**
	 * XML. All coordinates are saved relative to the original image.
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/xml-scheme-recognized-document">
	 * Output XML Document</a> for the description of tags. If you select this export format, barcodes
	 * are recognized on the image and saved to output XML no matter which profile is used for recognition.
	 */
	ExportFormatXml ExportFormat = "xml"

	/**
	 * XML. All coordinates are saved relative to the image after geometry correction.
	 *
	 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/xml-scheme-recognized-document">
	 * Output XML Document</a> for the description of tags. If you select this export format, barcodes
	 * are recognized on the image and saved to output XML no matter which profile is used for recognition.
	 */
	ExportFormatXmlForCorrectedImage ExportFormat = "xmlForCorrectedImage"

	/**
	 * ALTO
	 */
	ExportFormatAlto ExportFormat = "alto"
)

/**
 * ABBYY Cloud OCR SDK provides a set of processing profiles which are designed for the main usage
 * scenarios. Choose the profile which fits your situation, because selecting the wrong profile can
 *significantly reduce processing quality.
 */
type ProcessingProfile string

const (
	/**
	 * Suitable for converting documents into an editable format such as RTF or DOCX. This profile is
	 * focused on retaining the original document structure and appearance, including font styles,
	 * pictures, background color, etc.
	 */
	ProcessingProfileDocumentConversion ProcessingProfile = "documentConversion"

	/**
	 * Suitable for creating a digital archive of searchable documents. This profile is not intended for
	 * converting documents into an editable format. It is best optimized for creating searchable PDFs
	 * (with the recognized text saved under the image). All possible text is found and recognized on the
	 * input image, including text embedded into pictures.
	 */
	ProcessingProfileDocumentArchiving ProcessingProfile = "documentArchiving"

	/**
	 * Suitable for extracting all text from the input image, including small text areas of low quality.
	 * The document appearance and structure are ignored, pictures and tables are not detected.
	 *
	 * <b>Note:</b> This profile is not intended for converting documents into an editable format. It is designed for
	 * the situations when you need to retrieve the data from the image for some further processing on
	 * your side, such as extracting data from bills, receipts or invoices for use in accounting systems.
	 * Consider using XML export format, which allows you to access the recognized text coordinates.
	 */
	ProcessingProfileTextExtraction ProcessingProfile = "textExtraction"

	/**
	 * Suitable for barcode extraction. Extracts only barcodes (texts, pictures, or tables are not detected)
	 */
	ProcessingProfileBarcodeRecognition ProcessingProfile = "barcodeRecognition"
)

/**
 * Types of text supported by ABBYY Cloud OCR SDK.
 * Several types may be provided at once
 *
 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/text-types/"/>
 */
type TextType string

const (
	/**
	 * Common typographic type of text.
	 */
	TextTypeNormal TextType = "normal"

	/**
	 * The text is typed on a typewriter.
	 */
	TextTypeTypewriter TextType = "typewriter"

	/**
	 * The text is printed on a dot-matrix printer.
	 */
	TextTypeMatrix TextType = "matrix"

	/**
	 * A special set of characters including only
	 * digits written in ZIP-code style.
	 */
	TextTypeIndex TextType = "index"

	/**
	 * Handprinted text.
	 */
	TextTypeHandprinted TextType = "handprinted"

	/**
	 * A monospaced font, designed for Optical Character
	 * Recognition. Largely used by banks, credit card
	 * companies and similar businesses.
	 */
	TextTypeOcrA TextType = "ocrA"

	/**
	 * A font designed for Optical Character Recognition.
	 */
	TextTypeOcrB TextType = "ocrB"

	/**
	 * A special set of characters including only digits and
	 * A, B, C, D characters printed in magnetic ink.
	 * MICR (Magnetic Ink Character Recognition) characters
	 * are found in a variety of places, including personal
	 * checks.
	 */
	TextTypeE13b TextType = "e13b"

	/**
	 * A special set of characters, which includes only digits
	 * and A, B, C, D, E characters, written in MICR barcode
	 * font (CMC-7).
	 */
	TextTypeCmc7 TextType = "cmc7"

	/**
	 * Text printed in Gothic type.
	 */
	TextTypeGothic TextType = "gothic"
)

/**
 * Specifies the source of the image. It can be either a scanned image,
 * or a photograph created with a digital camera. Special preprocessing
 * operations can be performed with the image depending on the selected
 * source. For example, the system can automatically correct distorted
 * text lines, poor focus and lighting on photos.
 */
type ImageSource string

const (
	/**
	 * The image source is detected automatically.
	 */
	ImageSourceAuto ImageSource = "auto"

	/**
	 * Photo
	 */
	ImageSourcePhoto ImageSource = "photo"

	/**
	 * Scanner
	 */
	ImageSourceScanner ImageSource = "scanner"
)

/**
 * Specifies whether the result must be written as tagged PDF.
 */
type WriteTags string

const (
	/**
	 * Automatic selection: the tags are written into the output PDF file
	 * if it must comply with PDF/A-1a standard, and are not written otherwise.
	 */
	WriteTagsAuto WriteTags = "auto"

	/**
	 * Write tags
	 */
	WriteTagsWrite WriteTags = "write"

	/**
	 * Don't write tags
	 */
	WriteTagsDontWrite WriteTags = "dontWrite"
)

/**
 * Specifies the export format.
 */
type BusinessCardExportFormat string

const (
	/**
	 * XML
	 */
	BusinessCardExportFormatXml BusinessCardExportFormat = "xml"

	/**
	 * vCard
	 */
	BusinessCardExportFormatVCard BusinessCardExportFormat = "vCard"

	/**
	 * CSV
	 */
	BusinessCardExportFormatCsv BusinessCardExportFormat = "csv"
)

/**
 * Specifies the type of marking around letters
 *
 * <b>Note:</b> See <a href="https://www.ocrsdk.com/documentation/specifications/field-marking/"/>
 */
type MarkingType string

const (
	/**
	 * This value denotes plain text.
	 */
	MarkingTypeSimpleText MarkingType = "simpleText"

	/**
	 * This value specifies that the text is underlined.
	 */
	MarkingTypeUnderlinedText MarkingType = "underlinedText"

	/**
	 * This value specifies that the text is enclosed in a
	 * frame.
	 */
	MarkingTypeTextInFrame MarkingType = "textInFrame"

	/**
	 * This value specifies that the text is located in white
	 * fields on a gray background.
	 */
	MarkingTypeGreyBoxes MarkingType = "greyBoxes"

	/**
	 * This value specifies that the field where the text is
	 * located is a set of separate boxes.
	 */
	MarkingTypeCharBoxSeries MarkingType = "charBoxSeries"

	/**
	 * This value specifies that the field where the text is
	 * located is a comb.
	 */
	MarkingTypeSimpleComb MarkingType = "simpleComb"

	/**
	 * This value specifies that the field where the text is
	 * located is a comb and that this comb is also the bottom
	 * line of a frame.
	 */
	MarkingTypeCombInFrame MarkingType = "combInFrame"

	/**
	 * This value specifies that the field where the text is
	 * located is a frame and this frame is split by vertical
	 * lines.
	 */
	MarkingTypePartitionedFrame MarkingType = "partitionedFrame"
)

/**
 * Provides additional information about handprinted
 * letters writing style.
 */
type WritingStyle string

const (
	WritingStyleDefault    WritingStyle = "default"
	WritingStyleAmerican   WritingStyle = "american"
	WritingStyleGerman     WritingStyle = "german"
	WritingStyleRussian    WritingStyle = "russian"
	WritingStylePolish     WritingStyle = "polish"
	WritingStyleThai       WritingStyle = "thai"
	WritingStyleJapanese   WritingStyle = "japanese"
	WritingStyleArabic     WritingStyle = "arabic"
	WritingStyleBaltic     WritingStyle = "baltic"
	WritingStyleBritish    WritingStyle = "british"
	WritingStyleBulgarian  WritingStyle = "bulgarian"
	WritingStyleCanadian   WritingStyle = "canadian"
	WritingStyleCzech      WritingStyle = "czech"
	WritingStyleCroatian   WritingStyle = "croatian"
	WritingStyleFrench     WritingStyle = "french"
	WritingStyleGreek      WritingStyle = "greek"
	WritingStyleHungarian  WritingStyle = "hungarian"
	WritingStyleItalian    WritingStyle = "italian"
	WritingStyleRomanian   WritingStyle = "romanian"
	WritingStyleSlovak     WritingStyle = "slovak"
	WritingStyleSpanish    WritingStyle = "spanish"
	WritingStyleTurkish    WritingStyle = "turkish"
	WritingStyleUkrainian  WritingStyle = "ukrainian"
	WritingStyleCommon     WritingStyle = "common"
	WritingStyleChinese    WritingStyle = "chinese"
	WritingStyleAzerbaijan WritingStyle = "azerbaijan"
	WritingStyleKazakh     WritingStyle = "kazakh"
	WritingStyleKirgiz     WritingStyle = "kirgiz"
	WritingStyleLatvian    WritingStyle = "latvian"
)

/**
 * Specifies the type of the barcode recognizable by ABBYY Cloud OCR SDK.
 * Several barcode types may be provided.
 */
type BarcodeType string

const (
	/**
	 * Aztec is a high density two-dimensional matrix style bar code symbology
	 * that can encode up to 3832 digits, 3067 alphanumeric characters, or
	 * 1914 bytes of data. The symbol is built on a square grid with a bulls-eye
	 * pattern at its center.
	 */
	BarcodeTypeAztec BarcodeType = "aztec"

	/**
	 * Codabar is a self-checking, variable length barcode that can encode
	 * 16 data characters. It is used primarily for numeric data, but also
	 * encodes six special characters. Codabar is useful for encoding dollar
	 * and mathematical figures because a decimal point, plus sign, and minus
	 * sign can be encoded.
	 */
	BarcodeTypeCodabar BarcodeType = "codabar"

	/**
	 * Code 128 is an alphanumeric, very high-density, compact, variable length
	 * barcode scheme that can encode the full 128 ASCII character set. Each
	 * character is represented by three bars and three spaces totaling 11 modules.
	 * Each bar or space is one, two, three, or four modules wide with the total
	 * number of modules representing bars an even number and the total number of
	 * modules representing a space an odd number. Three different start characters
	 * are used to select one of three character sets.
	 */
	BarcodeTypeCode128 BarcodeType = "code128"

	/**
	 * Code 39, also referred to as Code 3 of 9, is an alphanumeric, self-checking,
	 * variable length barcode that uses five black bars and four spaces to define
	 * a character. Three of the elements are wide and six are narrow.
	 */
	BarcodeTypeCode39 BarcodeType = "code39"

	/**
	 * Code 93 is a variable length bar code that encodes 47 characters. It is named
	 * Code 93 because every character is constructed from nine elements arranged
	 * into three bars with their adjacent spaces. Code 93 is a compressed version
	 * of Code 39 and was designed to complement Code 39.
	 */
	BarcodeTypeCode93 BarcodeType = "code93"

	/**
	 * Data Matrix is a two-dimensional matrix barcode consisting of black and white
	 * modules arranged in either a square or rectangular pattern. Every Data Matrix
	 * is composed of two solid adjacent borders in an “L” shape and two other
	 * borders consisting of alternating dark and light modules. Within these
	 * borders are rows and columns of cells encoding information. A Data Matrix
	 * barcode can store up to 2335 alphanumeric characters.
	 */
	BarcodeTypeDataMatrix BarcodeType = "dataMatrix"

	/**
	 * The European Article Numbering (EAN) system is used for products that require
	 * a country origin. EAN8 and EAN13 are fixed-length barcodes used to encode
	 * either eight or thirteen characters. The first two characters identify the
	 * country of origin, the next characters are data characters, and the last
	 * character is the checksum. These barcodes may include an additional barcode
	 * to the right of the main barcode. This second barcode, which is usually not
	 * as tall as the primary barcode, is used to encode additional information for
	 * newspapers, books, and other periodicals. The supplemental barcode may
	 * either encoded 2 or 5 digits of information.
	 */
	BarcodeTypeEan8 BarcodeType = "ean8"

	/**
	 * The European Article Numbering (EAN) system is used for products that require
	 * a country origin. EAN8 and EAN13 are fixed-length barcodes used to encode
	 * either eight or thirteen characters. The first two characters identify the
	 * country of origin, the next characters are data characters, and the last
	 * character is the checksum. These barcodes may include an additional barcode
	 * to the right of the main barcode. This second barcode, which is usually not
	 * as tall as the primary barcode, is used to encode additional information for
	 * newspapers, books, and other periodicals. The supplemental barcode may
	 * either encoded 2 or 5 digits of information.
	 */
	BarcodeTypeEan13 BarcodeType = "ean13"

	/**
	 * IATA 2 of 5 is a barcode standard designed by the IATA (International Air
	 * Transport Association). This standard is used for all boarding passes.
	 */
	BarcodeTypeIata25 BarcodeType = "iata25"

	/**
	 * Industrial 2 of 5 is numeric-only barcode that has been in use a int time.
	 * Unlike Interleaved 2 of 5, all of the information is encoded in the bars
	 * the spaces are fixed width and are used only to separate the bars. The code
	 * is self-checking and does not include a checksum.
	 */
	BarcodeTypeIndustrial25 BarcodeType = "industrial25"

	/**
	 * Interleaved 2 of 5 is a variable length (must be a multiple of two),
	 * high-density, self-checking, numeric barcode that uses five black bars
	 * and five white bars to define a character. Two digits are encoded in every
	 * character one in the black bars and one in the white bars. Two of the
	 * black bars and two of the white bars are wide. The other bars are narrow.
	 */
	BarcodeTypeInterleaved25 BarcodeType = "interleaved25"

	/**
	 * Matrix 2 of 5 is self-checking numeric-only barcode. Unlike Interleaved
	 * 2 of 5, all of the information is encoded in the bars the spaces are
	 * fixed width and are used only to separate the bars. Matrix 2 of 5 is used
	 * primarily for warehouse sorting, photo finishing, and airline ticket marking.
	 */
	BarcodeTypeMatrix25 BarcodeType = "matrix25"

	/**
	 * Patch code contains a pattern of horizontal black bars separated by spaces.
	 * Typically, a patch code is placed near the top center of a paper document
	 * to be scanned and used as a document separator.
	 */
	BarcodeTypePatch BarcodeType = "patch"

	/**
	 * PDF417 is a variable length, two-dimensional (2D), stacked symbology that
	 * can store up to 2710 digits, 1850 printable ASCII characters or 1108 binary
	 * characters per symbol. PDF417 is designed with selectable levels of error
	 * correction. Its high data capacity can be helpful in applications where a
	 * large amount of data must travel with a labeled document or item.
	 */
	BarcodeTypePdf417 BarcodeType = "pdf417"

	/**
	 * The Postnet (Postal Numeric Encoding Technique) is a fixed length symbology
	 * (5, 6, 9, or 11 characters) which uses constant bar and space width.
	 * Information is encoded by varying the bar height between the two values.
	 * Postnet barcodes are placed on the lower right of envelopes or postcards,
	 * and are used to expedite the processing of mail with automatic equipment
	 * and provide reduced postage rates.
	 */
	BarcodeTypePostNet BarcodeType = "postNet"

	/**
	 * QR Code is a two-dimensional matrix barcode. The barcode has 3 large
	 * squares (registration marks) in the corners which define the top
	 * of the barcode. The black and white squares in the area between the
	 * registration marks are the encoded data and error correction keys. QR
	 * Codes can encode over 4000 ASCII characters.
	 */
	BarcodeTypeQrCode BarcodeType = "qrCode"

	/**
	 * This type of barcode is a 19 digit barcode with a 20th check digit. For a
	 * total of 20 digits. It typically is used for carton identification. Both for
	 * internal carton numbering and also for using the UCC-128 barcode on your
	 * cartons being shipped out to your customers.
	 */
	BarcodeTypeUcc128 BarcodeType = "ucc128"

	/**
	 * The UPC-A (Universal Product Code) barcode is 12 digits int, including its
	 * checksum. Each digit is represented by a seven-bit sequence, encoded by a
	 * series of alternating bars and spaces. UPC-A is used for marking products
	 * which are sold at retail in the USA.
	 */
	BarcodeTypeUpcA BarcodeType = "upcA"

	/**
	 * The UPC-E barcode is a shortened version of UPC-A barcode. It compresses
	 * the data characters and the checksum into six characters. This bar code is
	 * ideal for small packages because it is the smallest bar code.
	 */
	BarcodeTypeUpcE BarcodeType = "upcE"
)

/**
 * Specifies the type of the checkmark.
 */
type CheckmarkType string

const (
	/**
	 * Checkmark against an empty background
	 */
	CheckmarkTypeEmpty CheckmarkType = "empty"

	/**
	 * Checkmark in a circle
	 */
	CheckmarkTypeCircle CheckmarkType = "circle"

	/**
	 * Checkmark in a square
	 */
	CheckmarkTypeSquare CheckmarkType = "square"
)

/**
 * Specifies the country where the receipt was printed.
 * This parameter can contain several names of countries.
 */
type ReceiptRecognizingCountry string

const (
	ReceiptRecognizingCountryUk          ReceiptRecognizingCountry = "uk"
	ReceiptRecognizingCountryUsa         ReceiptRecognizingCountry = "usa"
	ReceiptRecognizingCountryAustralia   ReceiptRecognizingCountry = "australia"
	ReceiptRecognizingCountryCanada      ReceiptRecognizingCountry = "canada"
	ReceiptRecognizingCountryJapan       ReceiptRecognizingCountry = "japan"
	ReceiptRecognizingCountryGermany     ReceiptRecognizingCountry = "germany"
	ReceiptRecognizingCountryItaly       ReceiptRecognizingCountry = "italy"
	ReceiptRecognizingCountryFrance      ReceiptRecognizingCountry = "france"
	ReceiptRecognizingCountryBrazil      ReceiptRecognizingCountry = "brazil"
	ReceiptRecognizingCountryRussia      ReceiptRecognizingCountry = "russia"
	ReceiptRecognizingCountryChina       ReceiptRecognizingCountry = "china"
	ReceiptRecognizingCountryKorea       ReceiptRecognizingCountry = "korea"
	ReceiptRecognizingCountryNetherlands ReceiptRecognizingCountry = "netherlands"
	ReceiptRecognizingCountrySpain       ReceiptRecognizingCountry = "spain"
	ReceiptRecognizingCountrySingapore   ReceiptRecognizingCountry = "singapore"
	ReceiptRecognizingCountryTaiwan      ReceiptRecognizingCountry = "taiwan"
	ReceiptRecognizingCountryTurkey      ReceiptRecognizingCountry = "turkey"
	ReceiptRecognizingCountryPoland      ReceiptRecognizingCountry = "poland"
)

/**
 * Specifies if the coordinates of field regions should be
 * saved to the resulting XML file, and how the coordinates
 * should be specified: on the original or on the corrected
 * image
 */
type FieldRegionExportMode string

const (
	FieldRegionExportModeDoNotExport       FieldRegionExportMode = "doNotExport"
	FieldRegionExportModeForOriginalImage  FieldRegionExportMode = "forOriginalImage"
	FieldRegionExportModeForCorrectedImage FieldRegionExportMode = "forCorrectedImage"
)
