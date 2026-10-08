package model

import "fmt"

// DividerStyle is the divider style of a card, applied between sections.
type DividerStyle string

const (
	DividerStyleUnspecified DividerStyle = "DIVIDER_STYLE_UNSPECIFIED"
	DividerStyleSolid       DividerStyle = "SOLID_DIVIDER"
	DividerStyleNone        DividerStyle = "NO_DIVIDER"
)

// DisplayStyle determines how a card is displayed (Google Workspace add-ons only).
type DisplayStyle string

const (
	DisplayStyleUnspecified DisplayStyle = "DISPLAY_STYLE_UNSPECIFIED"
	DisplayStylePeek        DisplayStyle = "PEEK"
	DisplayStyleReplace     DisplayStyle = "REPLACE"
)

// Card is a card interface displayed in a Google Chat message or Google
// Workspace add-on.
type Card struct {
	Header              *CardHeader      `json:"header,omitempty"`
	Sections            []Section        `json:"sections,omitempty"`
	SectionDividerStyle DividerStyle     `json:"sectionDividerStyle,omitempty"`
	CardActions         []CardAction     `json:"cardActions,omitempty"`
	Name                string           `json:"name,omitempty"`
	FixedFooter         *CardFixedFooter `json:"fixedFooter,omitempty"`
	DisplayStyle        DisplayStyle     `json:"displayStyle,omitempty"`
	PeekCardHeader      *CardHeader      `json:"peekCardHeader,omitempty"`
}

// Validate checks the card and all of its descendants.
func (c *Card) Validate() error {
	if c == nil {
		return nil
	}
	for i := range c.Sections {
		if err := c.Sections[i].Validate(); err != nil {
			return fmt.Errorf("validating section %d: %w", i, err)
		}
	}
	for i := range c.CardActions {
		if err := c.CardActions[i].Validate(); err != nil {
			return fmt.Errorf("validating card action %d: %w", i, err)
		}
	}
	if err := c.FixedFooter.Validate(); err != nil {
		return fmt.Errorf("validating fixed footer: %w", err)
	}
	return nil
}

// CardHeader represents a card header.
type CardHeader struct {
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle,omitempty"`
	ImageType    ImageType `json:"imageType,omitempty"`
	ImageURL     string    `json:"imageUrl,omitempty"`
	ImageAltText string    `json:"imageAltText,omitempty"`
}

// Section contains a collection of widgets rendered vertically in the order
// they are specified.
type Section struct {
	Header                    string           `json:"header,omitempty"`
	Widgets                   []Widget         `json:"widgets,omitempty"`
	Collapsible               bool             `json:"collapsible,omitempty"`
	UncollapsibleWidgetsCount int32            `json:"uncollapsibleWidgetsCount,omitempty"`
	CollapseControl           *CollapseControl `json:"collapseControl,omitempty"`
}

// Validate checks the section's widgets.
func (s *Section) Validate() error {
	for i := range s.Widgets {
		if err := s.Widgets[i].Validate(); err != nil {
			return fmt.Errorf("validating widget %d: %w", i, err)
		}
	}
	return nil
}

// CardAction is an action associated with the card, such as a menu item in the
// card toolbar.
type CardAction struct {
	ActionLabel string   `json:"actionLabel,omitempty"`
	OnClick     *OnClick `json:"onClick,omitempty"`
}

// Validate checks the action's click handler.
func (a *CardAction) Validate() error {
	if err := a.OnClick.Validate(); err != nil {
		return fmt.Errorf("validating onClick: %w", err)
	}
	return nil
}

// NestedWidget is a list of widgets that can be displayed in a containing
// layout, such as a CarouselCard. Only one field may be set.
type NestedWidget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	ButtonList    *ButtonList    `json:"buttonList,omitempty"`
	Image         *Image         `json:"image,omitempty"`
}

// Validate checks that at most one widget type is set.
func (n *NestedWidget) Validate() error {
	return oneOf("NestedWidget data", map[string]bool{
		"textParagraph": n.TextParagraph != nil,
		"buttonList":    n.ButtonList != nil,
		"image":         n.Image != nil,
	})
}

// CardFixedFooter is a persistent (sticky) footer displayed at the bottom of
// the card.
type CardFixedFooter struct {
	PrimaryButton   *Button `json:"primaryButton,omitempty"`
	SecondaryButton *Button `json:"secondaryButton,omitempty"`
}

// Validate checks the footer's buttons.
func (f *CardFixedFooter) Validate() error {
	if f == nil {
		return nil
	}
	if err := f.PrimaryButton.Validate(); err != nil {
		return fmt.Errorf("validating primary button: %w", err)
	}
	if err := f.SecondaryButton.Validate(); err != nil {
		return fmt.Errorf("validating secondary button: %w", err)
	}
	return nil
}

// CollapseControl represents an expand and collapse control.
type CollapseControl struct {
	HorizontalAlignment HorizontalAlignment `json:"horizontalAlignment,omitempty"`
	ExpandButton        *Button             `json:"expandButton,omitempty"`
	CollapseButton      *Button             `json:"collapseButton,omitempty"`
}
