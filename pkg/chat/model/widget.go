package model

import "fmt"

// ImageType is the shape used to crop an image.
type ImageType string

const (
	ImageTypeSquare ImageType = "SQUARE"
	ImageTypeCircle ImageType = "CIRCLE"
)

// HorizontalAlignment specifies where widgets are aligned horizontally.
type HorizontalAlignment string

const (
	HorizontalAlignmentUnspecified HorizontalAlignment = "HORIZONTAL_ALIGNMENT_UNSPECIFIED"
	HorizontalAlignmentStart       HorizontalAlignment = "START"
	HorizontalAlignmentCenter      HorizontalAlignment = "CENTER"
	HorizontalAlignmentEnd         HorizontalAlignment = "END"
)

// VerticalAlignment specifies where widgets are aligned vertically.
type VerticalAlignment string

const (
	VerticalAlignmentUnspecified VerticalAlignment = "VERTICAL_ALIGNMENT_UNSPECIFIED"
	VerticalAlignmentTop         VerticalAlignment = "TOP"
	VerticalAlignmentMiddle      VerticalAlignment = "MIDDLE"
	VerticalAlignmentBottom      VerticalAlignment = "BOTTOM"
)

// Widget is a UI element contained in a card section. Only one of the data
// fields may be set.
type Widget struct {
	HorizontalAlignment HorizontalAlignment `json:"horizontalAlignment,omitempty"`

	TextParagraph  *TextParagraph  `json:"textParagraph,omitempty"`
	Image          *Image          `json:"image,omitempty"`
	DecoratedText  *DecoratedText  `json:"decoratedText,omitempty"`
	ButtonList     *ButtonList     `json:"buttonList,omitempty"`
	TextInput      *TextInput      `json:"textInput,omitempty"`
	SelectionInput *SelectionInput `json:"selectionInput,omitempty"`
	DateTimePicker *DateTimePicker `json:"dateTimePicker,omitempty"`
	Divider        *Divider        `json:"divider,omitempty"`
	Grid           *Grid           `json:"grid,omitempty"`
	Columns        *Columns        `json:"columns,omitempty"`
	Carousel       *Carousel       `json:"carousel,omitempty"`
	ChipList       *ChipList       `json:"chipList,omitempty"`
}

// Validate checks that at most one widget type is set, and validates it.
func (w *Widget) Validate() error {
	if err := oneOf("Widget data", map[string]bool{
		"textParagraph":  w.TextParagraph != nil,
		"image":          w.Image != nil,
		"decoratedText":  w.DecoratedText != nil,
		"buttonList":     w.ButtonList != nil,
		"textInput":      w.TextInput != nil,
		"selectionInput": w.SelectionInput != nil,
		"dateTimePicker": w.DateTimePicker != nil,
		"divider":        w.Divider != nil,
		"grid":           w.Grid != nil,
		"columns":        w.Columns != nil,
		"carousel":       w.Carousel != nil,
		"chipList":       w.ChipList != nil,
	}); err != nil {
		return fmt.Errorf("validating widget: %w", err)
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
	case w.Grid != nil:
		return wrap("grid", w.Grid.Validate())
	case w.Columns != nil:
		return wrap("columns", w.Columns.Validate())
	case w.Carousel != nil:
		return wrap("carousel", w.Carousel.Validate())
	case w.ChipList != nil:
		return wrap("chipList", w.ChipList.Validate())
	}
	return nil
}

func wrap(what string, err error) error {
	if err != nil {
		return fmt.Errorf("validating %s: %w", what, err)
	}
	return nil
}

// TextSyntax is the syntax used to format text.
type TextSyntax string

const (
	TextSyntaxUnspecified TextSyntax = "TEXT_SYNTAX_UNSPECIFIED"
	TextSyntaxHTML        TextSyntax = "HTML"
	TextSyntaxMarkdown    TextSyntax = "MARKDOWN"
)

// TextParagraph is a paragraph of text that supports formatting.
type TextParagraph struct {
	Text       string     `json:"text"`
	MaxLines   int32      `json:"maxLines,omitempty"`
	TextSyntax TextSyntax `json:"textSyntax,omitempty"`
}

// Image is an image specified by a URL that can have an onClick action.
type Image struct {
	ImageURL string   `json:"imageUrl"`
	OnClick  *OnClick `json:"onClick,omitempty"`
	AltText  string   `json:"altText,omitempty"`
}

// Divider displays a horizontal line divider between widgets.
type Divider struct{}

// SwitchControlType is the way a switch appears in the UI.
type SwitchControlType string

const (
	SwitchControlTypeSwitch   SwitchControlType = "SWITCH"
	SwitchControlTypeCheckbox SwitchControlType = "CHECKBOX"
	// SwitchControlTypeCheckBox is an alias accepted by the API for
	// SwitchControlTypeCheckbox.
	SwitchControlTypeCheckBox SwitchControlType = "CHECK_BOX"
)

// SwitchControl is a toggle-style switch or checkbox inside a DecoratedText.
type SwitchControl struct {
	Name           string            `json:"name,omitempty"`
	Value          string            `json:"value,omitempty"`
	Selected       bool              `json:"selected,omitempty"`
	OnChangeAction *Action           `json:"onChangeAction,omitempty"`
	ControlType    SwitchControlType `json:"controlType,omitempty"`
}

// DecoratedText displays text with optional decorations such as a label above
// or below the text, an icon in front of the text, a selection widget, or a
// button after the text. At most one of Button, SwitchControl and EndIcon may
// be set.
type DecoratedText struct {
	// Icon is deprecated; use StartIcon.
	Icon                   *Icon             `json:"icon,omitempty"`
	StartIcon              *Icon             `json:"startIcon,omitempty"`
	StartIconVerticalAlign VerticalAlignment `json:"startIconVerticalAlignment,omitempty"`
	TopLabel               string            `json:"topLabel,omitempty"`
	TopLabelText           *TextParagraph    `json:"topLabelText,omitempty"`
	Text                   string            `json:"text,omitempty"`
	ContentText            *TextParagraph    `json:"contentText,omitempty"`
	WrapText               bool              `json:"wrapText,omitempty"`
	BottomLabel            string            `json:"bottomLabel,omitempty"`
	BottomLabelText        *TextParagraph    `json:"bottomLabelText,omitempty"`
	OnClick                *OnClick          `json:"onClick,omitempty"`
	Button                 *Button           `json:"button,omitempty"`
	SwitchControl          *SwitchControl    `json:"switchControl,omitempty"`
	EndIcon                *Icon             `json:"endIcon,omitempty"`
}

// Validate checks the oneof constraints of the decorated text and its children.
func (d *DecoratedText) Validate() error {
	if err := oneOf("DecoratedText control", map[string]bool{
		"button":        d.Button != nil,
		"switchControl": d.SwitchControl != nil,
		"endIcon":       d.EndIcon != nil,
	}); err != nil {
		return err
	}
	if err := d.OnClick.Validate(); err != nil {
		return fmt.Errorf("validating onClick: %w", err)
	}
	if err := d.Button.Validate(); err != nil {
		return fmt.Errorf("validating button: %w", err)
	}
	if err := d.StartIcon.Validate(); err != nil {
		return fmt.Errorf("validating startIcon: %w", err)
	}
	if err := d.Icon.Validate(); err != nil {
		return fmt.Errorf("validating icon: %w", err)
	}
	if err := d.EndIcon.Validate(); err != nil {
		return fmt.Errorf("validating endIcon: %w", err)
	}
	return nil
}

// ButtonList is a list of buttons laid out horizontally.
type ButtonList struct {
	Buttons []Button `json:"buttons,omitempty"`
}

// Validate checks every button in the list.
func (b *ButtonList) Validate() error {
	for i := range b.Buttons {
		if err := b.Buttons[i].Validate(); err != nil {
			return fmt.Errorf("validating button %d: %w", i, err)
		}
	}
	return nil
}

// ChipListLayout is the layout of a ChipList.
type ChipListLayout string

const (
	ChipListLayoutUnspecified          ChipListLayout = "LAYOUT_UNSPECIFIED"
	ChipListLayoutWrapped              ChipListLayout = "WRAPPED"
	ChipListLayoutHorizontalScrollable ChipListLayout = "HORIZONTAL_SCROLLABLE"
)

// ChipList is a list of chips laid out horizontally, which can either scroll
// horizontally or wrap to the next line.
type ChipList struct {
	Layout ChipListLayout `json:"layout,omitempty"`
	Chips  []Chip         `json:"chips,omitempty"`
}

// Validate checks every chip in the list.
func (c *ChipList) Validate() error {
	for i := range c.Chips {
		if err := c.Chips[i].OnClick.Validate(); err != nil {
			return fmt.Errorf("validating chip %d onClick: %w", i, err)
		}
	}
	return nil
}

// Chip is a text, icon, or text and icon chip that users can click.
type Chip struct {
	Icon    *Icon    `json:"icon,omitempty"`
	Label   string   `json:"label,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
	// Enabled is deprecated; use Disabled.
	Enabled  bool   `json:"enabled,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	AltText  string `json:"altText,omitempty"`
}

// Carousel is a carousel of nested widgets.
type Carousel struct {
	CarouselCards []CarouselCard `json:"carouselCards,omitempty"`
}

// Validate checks every widget in every card.
func (c *Carousel) Validate() error {
	for i := range c.CarouselCards {
		if err := c.CarouselCards[i].Validate(); err != nil {
			return fmt.Errorf("validating carousel card %d: %w", i, err)
		}
	}
	return nil
}

// CarouselCard is a card that can be displayed as a carousel item.
type CarouselCard struct {
	Widgets       []NestedWidget `json:"widgets,omitempty"`
	FooterWidgets []NestedWidget `json:"footerWidgets,omitempty"`
}

// Validate checks the widgets and footer widgets.
func (c *CarouselCard) Validate() error {
	for i := range c.Widgets {
		if err := c.Widgets[i].Validate(); err != nil {
			return fmt.Errorf("validating widget %d: %w", i, err)
		}
	}
	for i := range c.FooterWidgets {
		if err := c.FooterWidgets[i].Validate(); err != nil {
			return fmt.Errorf("validating footer widget %d: %w", i, err)
		}
	}
	return nil
}
