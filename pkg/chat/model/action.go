package model

import "fmt"

// OnClick represents how to respond when users click an interactive element.
// Only one field may be set.
type OnClick struct {
	Action                *Action       `json:"action,omitempty"`
	OpenLink              *OpenLink     `json:"openLink,omitempty"`
	OpenDynamicLinkAction *Action       `json:"openDynamicLinkAction,omitempty"`
	Card                  *Card         `json:"card,omitempty"`
	OverflowMenu          *OverflowMenu `json:"overflowMenu,omitempty"`
}

// Validate checks that at most one handler is set. It is safe to call on a nil
// receiver.
func (o *OnClick) Validate() error {
	if o == nil {
		return nil
	}
	if err := oneOf("OnClick data", map[string]bool{
		"action":                o.Action != nil,
		"openLink":              o.OpenLink != nil,
		"openDynamicLinkAction": o.OpenDynamicLinkAction != nil,
		"card":                  o.Card != nil,
		"overflowMenu":          o.OverflowMenu != nil,
	}); err != nil {
		return err
	}
	if err := o.Card.Validate(); err != nil {
		return fmt.Errorf("validating nested card: %w", err)
	}
	if err := o.OverflowMenu.Validate(); err != nil {
		return fmt.Errorf("validating overflow menu: %w", err)
	}
	return nil
}

// OpenAs is how a link opens.
type OpenAs string

const (
	OpenAsFullSize OpenAs = "FULL_SIZE"
	OpenAsOverlay  OpenAs = "OVERLAY"
)

// OnClose is what the app does when a link opened by an OnClick closes.
type OnClose string

const (
	OnCloseNothing OnClose = "NOTHING"
	OnCloseReload  OnClose = "RELOAD"
)

// OpenLink represents an onClick event that opens a hyperlink.
type OpenLink struct {
	URL     string  `json:"url"`
	OpenAs  OpenAs  `json:"openAs,omitempty"`
	OnClose OnClose `json:"onClose,omitempty"`
}

// LoadIndicator specifies the loading indicator shown while an action runs.
type LoadIndicator string

const (
	LoadIndicatorSpinner LoadIndicator = "SPINNER"
	LoadIndicatorNone    LoadIndicator = "NONE"
)

// Interaction is what an action does in addition to running its function.
type Interaction string

const (
	InteractionUnspecified Interaction = "INTERACTION_UNSPECIFIED"
	InteractionOpenDialog  Interaction = "OPEN_DIALOG"
)

// ActionParameter is a list of string parameters supplied when an action
// method is invoked.
type ActionParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Action describes the behaviour triggered by a user interaction.
type Action struct {
	Function              string            `json:"function,omitempty"`
	Parameters            []ActionParameter `json:"parameters,omitempty"`
	LoadIndicator         LoadIndicator     `json:"loadIndicator,omitempty"`
	PersistValues         bool              `json:"persistValues,omitempty"`
	Interaction           Interaction       `json:"interaction,omitempty"`
	RequiredWidgets       []string          `json:"requiredWidgets,omitempty"`
	AllWidgetsAreRequired bool              `json:"allWidgetsAreRequired,omitempty"`
}

// OverflowMenu is a widget that presents a popup menu of actions.
type OverflowMenu struct {
	Items []OverflowMenuItem `json:"items,omitempty"`
}

// Validate checks every item's click handler.
func (m *OverflowMenu) Validate() error {
	if m == nil {
		return nil
	}
	for i := range m.Items {
		if err := m.Items[i].OnClick.Validate(); err != nil {
			return fmt.Errorf("validating item %d onClick: %w", i, err)
		}
	}
	return nil
}

// OverflowMenuItem is an option users can invoke in an overflow menu.
type OverflowMenuItem struct {
	StartIcon *Icon    `json:"startIcon,omitempty"`
	Text      string   `json:"text,omitempty"`
	OnClick   *OnClick `json:"onClick,omitempty"`
	Disabled  bool     `json:"disabled,omitempty"`
}

// ButtonType is the style of a button.
type ButtonType string

const (
	ButtonTypeUnspecified ButtonType = "TYPE_UNSPECIFIED"
	ButtonTypeOutlined    ButtonType = "OUTLINED"
	ButtonTypeFilled      ButtonType = "FILLED"
	ButtonTypeFilledTonal ButtonType = "FILLED_TONAL"
	ButtonTypeBorderless  ButtonType = "BORDERLESS"
)

// Button is a text, icon, or text and icon button that users can click.
type Button struct {
	Text     string     `json:"text,omitempty"`
	Icon     *Icon      `json:"icon,omitempty"`
	Color    *Color     `json:"color,omitempty"`
	OnClick  *OnClick   `json:"onClick,omitempty"`
	Disabled bool       `json:"disabled,omitempty"`
	AltText  string     `json:"altText,omitempty"`
	Type     ButtonType `json:"type,omitempty"`
}

// Validate checks the button's icon and click handler. It is safe to call on a
// nil receiver.
func (b *Button) Validate() error {
	if b == nil {
		return nil
	}
	if err := b.Icon.Validate(); err != nil {
		return fmt.Errorf("validating icon: %w", err)
	}
	if err := b.OnClick.Validate(); err != nil {
		return fmt.Errorf("validating onClick: %w", err)
	}
	return nil
}

// Icon is an icon displayed in a widget on a card. At most one of KnownIcon,
// IconURL and MaterialIcon may be set.
type Icon struct {
	KnownIcon    string        `json:"knownIcon,omitempty"`
	IconURL      string        `json:"iconUrl,omitempty"`
	MaterialIcon *MaterialIcon `json:"materialIcon,omitempty"`
	AltText      string        `json:"altText,omitempty"`
	ImageType    ImageType     `json:"imageType,omitempty"`
}

// Validate checks that at most one icon source is set. It is safe to call on
// a nil receiver.
func (i *Icon) Validate() error {
	if i == nil {
		return nil
	}
	return oneOf("Icon icons", map[string]bool{
		"knownIcon":    i.KnownIcon != "",
		"iconUrl":      i.IconURL != "",
		"materialIcon": i.MaterialIcon != nil,
	})
}

// MaterialIcon is a Google Material Icon.
type MaterialIcon struct {
	Name   string `json:"name"`
	Fill   bool   `json:"fill,omitempty"`
	Weight int32  `json:"weight,omitempty"`
	Grade  int32  `json:"grade,omitempty"`
}

// Color is a color in the RGBA color space (google.type.Color). Each channel
// is in the range [0, 1].
type Color struct {
	Red   float32 `json:"red,omitempty"`
	Green float32 `json:"green,omitempty"`
	Blue  float32 `json:"blue,omitempty"`
	// Alpha is nil for fully opaque; the zero value of a pointer differs from
	// an explicit 0 (fully transparent).
	Alpha *float32 `json:"alpha,omitempty"`
}
