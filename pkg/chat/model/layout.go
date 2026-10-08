package model

import "fmt"

// ImageCropType is the crop style applied to an image.
type ImageCropType string

const (
	ImageCropTypeUnspecified     ImageCropType = "IMAGE_CROP_TYPE_UNSPECIFIED"
	ImageCropTypeSquare          ImageCropType = "SQUARE"
	ImageCropTypeCircle          ImageCropType = "CIRCLE"
	ImageCropTypeRectangleCustom ImageCropType = "RECTANGLE_CUSTOM"
	ImageCropTypeRectangle43     ImageCropType = "RECTANGLE_4_3"
)

// ImageCropStyle represents the crop style applied to an image.
type ImageCropStyle struct {
	Type ImageCropType `json:"type,omitempty"`
	// AspectRatio is used when Type is ImageCropTypeRectangleCustom.
	AspectRatio float64 `json:"aspectRatio,omitempty"`
}

// BorderType is the type of border applied to a component.
type BorderType string

const (
	BorderTypeUnspecified BorderType = "BORDER_TYPE_UNSPECIFIED"
	BorderTypeNoBorder    BorderType = "NO_BORDER"
	BorderTypeStroke      BorderType = "STROKE"
)

// BorderStyle represents the style of a border.
type BorderStyle struct {
	Type         BorderType `json:"type,omitempty"`
	StrokeColor  *Color     `json:"strokeColor,omitempty"`
	CornerRadius int32      `json:"cornerRadius,omitempty"`
}

// ImageComponent represents an image within a grid item.
type ImageComponent struct {
	ImageURI    string          `json:"imageUri,omitempty"`
	AltText     string          `json:"altText,omitempty"`
	CropStyle   *ImageCropStyle `json:"cropStyle,omitempty"`
	BorderStyle *BorderStyle    `json:"borderStyle,omitempty"`
}

// GridItemLayout is how a grid item's text is positioned.
type GridItemLayout string

const (
	GridItemLayoutUnspecified GridItemLayout = "GRID_ITEM_LAYOUT_UNSPECIFIED"
	GridItemLayoutTextBelow   GridItemLayout = "TEXT_BELOW"
	GridItemLayoutTextAbove   GridItemLayout = "TEXT_ABOVE"
)

// GridItem represents an item in a grid layout.
type GridItem struct {
	ID       string          `json:"id,omitempty"`
	Image    *ImageComponent `json:"image,omitempty"`
	Title    string          `json:"title,omitempty"`
	Subtitle string          `json:"subtitle,omitempty"`
	Layout   GridItemLayout  `json:"layout,omitempty"`
}

// Grid displays a grid with a collection of items.
type Grid struct {
	Title       string       `json:"title,omitempty"`
	Items       []GridItem   `json:"items,omitempty"`
	BorderStyle *BorderStyle `json:"borderStyle,omitempty"`
	ColumnCount int32        `json:"columnCount,omitempty"`
	OnClick     *OnClick     `json:"onClick,omitempty"`
}

// Validate checks the grid's click handler.
func (g *Grid) Validate() error {
	if err := g.OnClick.Validate(); err != nil {
		return fmt.Errorf("validating onClick: %w", err)
	}
	return nil
}

// HorizontalSizeStyle specifies how a column fills the width of the card.
type HorizontalSizeStyle string

const (
	HorizontalSizeStyleUnspecified        HorizontalSizeStyle = "HORIZONTAL_SIZE_STYLE_UNSPECIFIED"
	HorizontalSizeStyleFillAvailableSpace HorizontalSizeStyle = "FILL_AVAILABLE_SPACE"
	HorizontalSizeStyleFillMinimumSpace   HorizontalSizeStyle = "FILL_MINIMUM_SPACE"
)

// ColumnVerticalAlignment specifies whether widgets align to the top, bottom
// or center of a column.
type ColumnVerticalAlignment string

const (
	ColumnVerticalAlignmentUnspecified ColumnVerticalAlignment = "VERTICAL_ALIGNMENT_UNSPECIFIED"
	ColumnVerticalAlignmentCenter      ColumnVerticalAlignment = "CENTER"
	ColumnVerticalAlignmentTop         ColumnVerticalAlignment = "TOP"
	ColumnVerticalAlignmentBottom      ColumnVerticalAlignment = "BOTTOM"
)

// Columns displays up to 2 columns in a card or dialog.
type Columns struct {
	ColumnItems []Column `json:"columnItems,omitempty"`
}

// Validate checks every column.
func (c *Columns) Validate() error {
	for i := range c.ColumnItems {
		if err := c.ColumnItems[i].Validate(); err != nil {
			return fmt.Errorf("validating column %d: %w", i, err)
		}
	}
	return nil
}

// Column is a column within a Columns widget.
type Column struct {
	HorizontalSizeStyle HorizontalSizeStyle     `json:"horizontalSizeStyle,omitempty"`
	HorizontalAlignment HorizontalAlignment     `json:"horizontalAlignment,omitempty"`
	VerticalAlignment   ColumnVerticalAlignment `json:"verticalAlignment,omitempty"`
	Widgets             []ColumnWidget          `json:"widgets,omitempty"`
}

// Validate checks every widget in the column.
func (c *Column) Validate() error {
	for i := range c.Widgets {
		if err := c.Widgets[i].Validate(); err != nil {
			return fmt.Errorf("validating widget %d: %w", i, err)
		}
	}
	return nil
}

// ColumnWidget is a widget that can be placed in a column. Only one field may
// be set.
type ColumnWidget struct {
	TextParagraph  *TextParagraph  `json:"textParagraph,omitempty"`
	Image          *Image          `json:"image,omitempty"`
	DecoratedText  *DecoratedText  `json:"decoratedText,omitempty"`
	ButtonList     *ButtonList     `json:"buttonList,omitempty"`
	TextInput      *TextInput      `json:"textInput,omitempty"`
	SelectionInput *SelectionInput `json:"selectionInput,omitempty"`
	DateTimePicker *DateTimePicker `json:"dateTimePicker,omitempty"`
	ChipList       *ChipList       `json:"chipList,omitempty"`
}

// Validate checks that at most one widget type is set, and validates it.
func (w *ColumnWidget) Validate() error {
	if err := oneOf("ColumnWidget data", map[string]bool{
		"textParagraph":  w.TextParagraph != nil,
		"image":          w.Image != nil,
		"decoratedText":  w.DecoratedText != nil,
		"buttonList":     w.ButtonList != nil,
		"textInput":      w.TextInput != nil,
		"selectionInput": w.SelectionInput != nil,
		"dateTimePicker": w.DateTimePicker != nil,
		"chipList":       w.ChipList != nil,
	}); err != nil {
		return err
	}
	switch {
	case w.Image != nil:
		return wrap("image", w.Image.OnClick.Validate())
	case w.DecoratedText != nil:
		return wrap("decoratedText", w.DecoratedText.Validate())
	case w.ButtonList != nil:
		return wrap("buttonList", w.ButtonList.Validate())
	case w.SelectionInput != nil:
		return wrap("selectionInput", w.SelectionInput.Validate())
	case w.ChipList != nil:
		return wrap("chipList", w.ChipList.Validate())
	}
	return nil
}
