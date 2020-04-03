//Schema for representing OCR results exported from FineReader 10.0 SDK. Copyright 2001-2011 ABBYY, Inc.
package abbyyxml

import (
	"encoding/xml"
	"io"
)

func Decode(r io.Reader) (*Document, error) {
	var doc Document
	if err := xml.NewDecoder(r).Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func Unmarshal(data []byte) (*Document, error) {
	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// ParagraphStyles ...
type ParagraphStyles struct {
	XMLName        xml.Name              `xml:"paragraphStyles"`
	ParagraphStyle []*ParagraphStyleType `xml:"paragraphStyle"`
}

// Sections ...
type Sections struct {
	XMLName xml.Name       `xml:"sections"`
	Section []*SectionType `xml:"section"`
}

// DocumentData ...
type DocumentData struct {
	XMLName         xml.Name           `xml:"documentData"`
	ParagraphStyles []*ParagraphStyles `xml:"paragraphStyles"`
	Sections        []*Sections        `xml:"sections"`
}

// Page ...
type Page struct {
	XMLName            xml.Name           `xml:"page"`
	WidthAttr          int                `xml:"width,attr"`
	HeightAttr         int                `xml:"height,attr"`
	ResolutionAttr     int                `xml:"resolution,attr"`
	OriginalCoordsAttr bool               `xml:"originalCoords,attr,omitempty"`
	RotationAttr       string             `xml:"rotation,attr,omitempty"`
	Block              []*BlockType       `xml:"block"`
	PageSection        []*PageSectionType `xml:"pageSection"`
	PageStream         []*PageStreamType  `xml:"pageStream"`
}

// Document ...
type Document struct {
	XMLName          xml.Name        `xml:"document"`
	VersionAttr      string          `xml:"version,attr"`
	ProducerAttr     string          `xml:"producer,attr"`
	PagesCountAttr   int             `xml:"pagesCount,attr,omitempty"`
	MainLanguageAttr string          `xml:"mainLanguage,attr,omitempty"`
	LanguagesAttr    string          `xml:"languages,attr,omitempty"`
	DocumentData     []*DocumentData `xml:"documentData"`
	Page             []*Page         `xml:"page"`
}

// ParagraphStyleType ...
type ParagraphStyleType struct {
	IdAttr              string           `xml:"id,attr"`
	NameAttr            string           `xml:"name,attr"`
	MainFontStyleIdAttr string           `xml:"mainFontStyleId,attr"`
	RoleAttr            string           `xml:"role,attr"`
	FontStyle           []*FontStyleType `xml:"fontStyle"`
}

// RoleLevel ...
type RoleLevel int

// Align ...
type Align string

// Rtl ...
type Rtl bool

// Before ...
type Before int

// After ...
type After int

// StartIndent ...
type StartIndent int

// LeftIndent ...
type LeftIndent int

// RightIndent ...
type RightIndent int

// LineSpacing ...
type LineSpacing int

// LineSpacingRatio ...
type LineSpacingRatio *float64

// FixedLineSpacing ...
type FixedLineSpacing bool

// FontStyleType ...
type FontStyleType struct {
	IdAttr              string   `xml:"id,attr"`
	BaseFontAttr        bool     `xml:"baseFont,attr,omitempty"`
	ItalicAttr          bool     `xml:"italic,attr,omitempty"`
	BoldAttr            bool     `xml:"bold,attr,omitempty"`
	UnderlineAttr       bool     `xml:"underline,attr,omitempty"`
	StrikeoutAttr       bool     `xml:"strikeout,attr,omitempty"`
	SmallcapsAttr       bool     `xml:"smallcaps,attr,omitempty"`
	ScalingAttr         int      `xml:"scaling,attr,omitempty"`
	SpacingAttr         int      `xml:"spacing,attr,omitempty"`
	ColorAttr           int      `xml:"color,attr,omitempty"`
	BackgroundColorAttr int      `xml:"backgroundColor,attr,omitempty"`
	FfAttr              string   `xml:"ff,attr"`
	FsAttr              *float64 `xml:"fs,attr"`
}

// PageSectionType ...
type PageSectionType struct {
	PageStream []*PageStreamType `xml:"pageStream"`
}

// PageStreamType ...
type PageStreamType struct {
	StreamTypeAttr string             `xml:"streamType,attr"`
	PageElement    []*PageElementType `xml:"pageElement"`
}

// PageElementType ...
type PageElementType struct {
	PageElemIdAttr string         `xml:"pageElemId,attr"`
	Text           []*TextType    `xml:"text"`
	Table          []*TableType   `xml:"table"`
	Barcode        []*BarcodeType `xml:"barcode"`
	Picture        []*PictureType `xml:"picture"`
}

// TableCell ...
type TableCell struct {
	XMLName               xml.Name           `xml:"tableCell"`
	TopPosAttr            int                `xml:"topPos,attr"`
	BottomPosAttr         int                `xml:"bottomPos,attr"`
	LeftPosAttr           int                `xml:"leftPos,attr"`
	RightPosAttr          int                `xml:"rightPos,attr"`
	VerticalAlignmentAttr string             `xml:"VerticalAlignment,attr"`
	Text                  []*PageElementType `xml:"text"`
}

// TableType ...
type TableType struct {
	IdAttr    string         `xml:"id,attr"`
	Caption   []*CaptionType `xml:"caption"`
	TableCell []*TableCell   `xml:"tableCell"`
}

// PictureType ...
type PictureType struct {
	IdAttr  string         `xml:"id,attr"`
	Caption []*CaptionType `xml:"caption"`
}

// BarcodeType ...
type BarcodeType struct {
	BarcodeValueAttr string `xml:"BarcodeValue,attr"`
}

// CaptionType ...
type CaptionType struct {
	PageElement []*PageElementType `xml:"pageElement"`
}

// SectionType ...
type SectionType struct {
	Stream []*TextStreamType `xml:"stream"`
}

// MainText ...
type MainText struct {
	XMLName         xml.Name `xml:"mainText"`
	RtlAttr         bool     `xml:"rtl,attr,omitempty"`
	ColumnCountAttr int      `xml:"columnCount,attr"`
}

// ElemId ...
type ElemId struct {
	XMLName xml.Name `xml:"elemId"`
	IdAttr  string   `xml:"id,attr"`
}

// TextStreamType ...
type TextStreamType struct {
	RoleAttr string      `xml:"role,attr,omitempty"`
	MainText []*MainText `xml:"mainText"`
	ElemId   []*ElemId   `xml:"elemId"`
}

// VertCjk ...
type VertCjk bool

// BeginPage ...
type BeginPage int

// EndPage ...
type EndPage int

// Rect ...
type Rect struct {
	XMLName xml.Name `xml:"rect"`
	LAttr   int      `xml:"l,attr"`
	TAttr   int      `xml:"t,attr"`
	RAttr   int      `xml:"r,attr"`
	BAttr   int      `xml:"b,attr"`
}

// Region ...
type Region struct {
	XMLName xml.Name `xml:"region"`
	Rect    []*Rect  `xml:"rect"`
}

// SeparatorsBox ...
type SeparatorsBox struct {
	XMLName   xml.Name              `xml:"separatorsBox"`
	Separator []*SeparatorBlockType `xml:"separator"`
}

// BlockType ...
type BlockType struct {
	BlockTypeAttr string                `xml:"blockType,attr"`
	Region        []*Region             `xml:"region"`
	Text          []*TextType           `xml:"text"`
	Row           []*TableRowType       `xml:"row"`
	SeparatorsBox []*SeparatorsBox      `xml:"separatorsBox"`
	Separator     []*SeparatorBlockType `xml:"separator"`
	BarcodeInfo   []*BarcodeInfoType    `xml:"barcodeInfo"`
}

// PageElemId ...
type PageElemId string

// BlockName ...
type BlockName string

// IsHidden ...
type IsHidden bool

// L ...
type L int

// T ...
type T int

// R ...
type R int

// B ...
type B int

// TextType ...
type TextType struct {
	IdAttr          string           `xml:"id,attr,omitempty"`
	OrientationAttr string           `xml:"orientation,attr,omitempty"`
	Par             []*ParagraphType `xml:"par"`
}

// BackgroundColor ...
type BackgroundColor int

// Mirrored ...
type Mirrored bool

// Inverted ...
type Inverted bool

// Cell ...
type Cell struct {
	XMLName          xml.Name    `xml:"cell"`
	ColSpanAttr      int         `xml:"colSpan,attr,omitempty"`
	RowSpanAttr      int         `xml:"rowSpan,attr,omitempty"`
	AlignAttr        string      `xml:"align,attr,omitempty"`
	PictureAttr      bool        `xml:"picture,attr,omitempty"`
	LeftBorderAttr   string      `xml:"leftBorder,attr,omitempty"`
	TopBorderAttr    string      `xml:"topBorder,attr,omitempty"`
	RightBorderAttr  string      `xml:"rightBorder,attr,omitempty"`
	BottomBorderAttr string      `xml:"bottomBorder,attr,omitempty"`
	WidthAttr        int         `xml:"width,attr"`
	HeightAttr       int         `xml:"height,attr"`
	Text             []*TextType `xml:"text"`
}

// TableRowType ...
type TableRowType struct {
	Cell []*Cell `xml:"cell"`
}

// ParagraphType ...
type ParagraphType struct {
	DropCapCharsCountAttr int         `xml:"dropCapCharsCount,attr,omitempty"`
	DropCaplAttr          int         `xml:"dropCap-l,attr,omitempty"`
	DropCaptAttr          int         `xml:"dropCap-t,attr,omitempty"`
	DropCaprAttr          int         `xml:"dropCap-r,attr,omitempty"`
	DropCapbAttr          int         `xml:"dropCap-b,attr,omitempty"`
	AlignAttr             string      `xml:"align,attr,omitempty"`
	LeftIndentAttr        int         `xml:"leftIndent,attr,omitempty"`
	RightIndentAttr       int         `xml:"rightIndent,attr,omitempty"`
	StartIndentAttr       int         `xml:"startIndent,attr,omitempty"`
	LineSpacingAttr       int         `xml:"lineSpacing,attr,omitempty"`
	IdAttr                string      `xml:"id,attr,omitempty"`
	StyleAttr             string      `xml:"style,attr,omitempty"`
	HasOverflowedHeadAttr bool        `xml:"hasOverflowedHead,attr,omitempty"`
	HasOverflowedTailAttr bool        `xml:"hasOverflowedTail,attr,omitempty"`
	IsListItemAttr        bool        `xml:"isListItem,attr,omitempty"`
	LstLvlAttr            int         `xml:"lstLvl,attr,omitempty"`
	LstNumAttr            int         `xml:"lstNum,attr,omitempty"`
	Line                  []*LineType `xml:"line"`
}

// ParagraphAlignment ...
type ParagraphAlignment string

// LineType ...
type LineType struct {
	BaselineAttr int               `xml:"baseline,attr"`
	LAttr        int               `xml:"l,attr"`
	TAttr        int               `xml:"t,attr"`
	RAttr        int               `xml:"r,attr"`
	BAttr        int               `xml:"b,attr"`
	Formatting   []*FormattingType `xml:"formatting"`
}

// WordRecVariants ...
type WordRecVariants struct {
	XMLName        xml.Name                  `xml:"wordRecVariants"`
	WordRecVariant []*WordRecognitionVariant `xml:"wordRecVariant"`
}

// FormattingType ...
type FormattingType struct {
	LangAttr          string           `xml:"lang,attr"`
	FfAttr            string           `xml:"ff,attr,omitempty"`
	FsAttr            *float64         `xml:"fs,attr,omitempty"`
	BoldAttr          bool             `xml:"bold,attr,omitempty"`
	ItalicAttr        bool             `xml:"italic,attr,omitempty"`
	SubscriptAttr     bool             `xml:"subscript,attr,omitempty"`
	SuperscriptAttr   bool             `xml:"superscript,attr,omitempty"`
	SmallcapsAttr     bool             `xml:"smallcaps,attr,omitempty"`
	UnderlineAttr     bool             `xml:"underline,attr,omitempty"`
	StrikeoutAttr     bool             `xml:"strikeout,attr,omitempty"`
	ColorAttr         int              `xml:"color,attr,omitempty"`
	ScalingAttr       int              `xml:"scaling,attr,omitempty"`
	SpacingAttr       int              `xml:"spacing,attr,omitempty"`
	StyleAttr         string           `xml:"style,attr,omitempty"`
	Base64encodedAttr bool             `xml:"base64encoded,attr,omitempty"`
	CharParams        *CharParamsType  `xml:"charParams"`
	WordRecVariants   *WordRecVariants `xml:"wordRecVariants"`
	Value             string           `xml:",innerxml"`
}

// VariantText ...
type VariantText struct {
	XMLName    xml.Name          `xml:"variantText"`
	CharParams []*CharParamsType `xml:"charParams"`
}

// WordRecognitionVariant ...
type WordRecognitionVariant struct {
	WordFromDictionaryAttr bool           `xml:"wordFromDictionary,attr,omitempty"`
	WordNormalAttr         bool           `xml:"wordNormal,attr,omitempty"`
	WordNumericAttr        bool           `xml:"wordNumeric,attr,omitempty"`
	WordIdentifierAttr     bool           `xml:"wordIdentifier,attr,omitempty"`
	WordPenaltyAttr        int            `xml:"wordPenalty,attr,omitempty"`
	MeanStrokeWidthAttr    int            `xml:"meanStrokeWidth,attr,omitempty"`
	VariantText            []*VariantText `xml:"variantText"`
}

// CharRecognitionVariant ...
type CharRecognitionVariant struct {
	CharConfidenceAttr   int `xml:"charConfidence,attr,omitempty"`
	SerifProbabilityAttr int `xml:"serifProbability,attr,omitempty"`
}

// CharRecVariants ...
type CharRecVariants struct {
	XMLName        xml.Name                  `xml:"charRecVariants"`
	CharRecVariant []*CharRecognitionVariant `xml:"charRecVariant"`
}

// CharParamsType ...
type CharParamsType struct {
	LAttr                  int              `xml:"l,attr"`
	TAttr                  int              `xml:"t,attr"`
	RAttr                  int              `xml:"r,attr"`
	BAttr                  int              `xml:"b,attr"`
	SuspiciousAttr         bool             `xml:"suspicious,attr,omitempty"`
	ProofedAttr            bool             `xml:"proofed,attr,omitempty"`
	WordStartAttr          bool             `xml:"wordStart,attr,omitempty"`
	WordFirstAttr          bool             `xml:"wordFirst,attr,omitempty"`
	WordLeftMostAttr       bool             `xml:"wordLeftMost,attr,omitempty"`
	WordFromDictionaryAttr bool             `xml:"wordFromDictionary,attr,omitempty"`
	WordNormalAttr         bool             `xml:"wordNormal,attr,omitempty"`
	WordNumericAttr        bool             `xml:"wordNumeric,attr,omitempty"`
	WordIdentifierAttr     bool             `xml:"wordIdentifier,attr,omitempty"`
	CharConfidenceAttr     int              `xml:"charConfidence,attr,omitempty"`
	SerifProbabilityAttr   int              `xml:"serifProbability,attr,omitempty"`
	WordPenaltyAttr        int              `xml:"wordPenalty,attr,omitempty"`
	MeanStrokeWidthAttr    int              `xml:"meanStrokeWidth,attr,omitempty"`
	CharacterHeightAttr    int              `xml:"characterHeight,attr,omitempty"`
	HasUncertainHeightAttr bool             `xml:"hasUncertainHeight,attr,omitempty"`
	BaseLineAttr           int              `xml:"baseLine,attr,omitempty"`
	IsTabAttr              bool             `xml:"isTab,attr,omitempty"`
	TabLeaderCountAttr     int              `xml:"tabLeaderCount,attr,omitempty"`
	CharRecVariants        *CharRecVariants `xml:"charRecVariants"`
}

// TableCellBorderType ...
type TableCellBorderType string

// SeparatorBlockType ...
type SeparatorBlockType struct {
	ThicknessAttr int      `xml:"thickness,attr"`
	TypeAttr      string   `xml:"type,attr"`
	Start         []*Point `xml:"start"`
	End           []*Point `xml:"end"`
}

// Point ...
type Point struct {
	XAttr int `xml:"x,attr"`
	YAttr int `xml:"y,attr"`
}

// BarcodeInfoType ...
type BarcodeInfoType struct {
	TypeAttr       string `xml:"type,attr"`
	SupplementAttr string `xml:"supplement,attr,omitempty"`
}

// BarcodeTypeEnum ...
type BarcodeTypeEnum string

// BarcodeSupplementEnum ...
type BarcodeSupplementEnum string
