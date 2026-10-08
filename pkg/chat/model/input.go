package model

// TextInputType is how a text input appears in the UI.
type TextInputType string

const (
	TextInputTypeSingleLine   TextInputType = "SINGLE_LINE"
	TextInputTypeMultipleLine TextInputType = "MULTIPLE_LINE"
)

// TextInput is a field into which users can enter text.
type TextInput struct {
	Name               string        `json:"name"`
	Label              string        `json:"label,omitempty"`
	HintText           string        `json:"hintText,omitempty"`
	Value              string        `json:"value,omitempty"`
	Type               TextInputType `json:"type,omitempty"`
	OnChangeAction     *Action       `json:"onChangeAction,omitempty"`
	InitialSuggestions *Suggestions  `json:"initialSuggestions,omitempty"`
	AutoCompleteAction *Action       `json:"autoCompleteAction,omitempty"`
	Validation         *Validation   `json:"validation,omitempty"`
	PlaceholderText    string        `json:"placeholderText,omitempty"`
}

// Suggestions are values users can enter in a text input.
type Suggestions struct {
	Items []SuggestionItem `json:"items,omitempty"`
}

// SuggestionItem is a single suggested input value.
type SuggestionItem struct {
	Text string `json:"text,omitempty"`
}

// InputType is the type of a validated input.
type InputType string

const (
	InputTypeUnspecified InputType = "INPUT_TYPE_UNSPECIFIED"
	InputTypeText        InputType = "TEXT"
	InputTypeInteger     InputType = "INTEGER"
	InputTypeFloat       InputType = "FLOAT"
	InputTypeEmail       InputType = "EMAIL"
	InputTypeEmojiPicker InputType = "EMOJI_PICKER"
)

// Validation holds the validation data for input widgets.
type Validation struct {
	CharacterLimit int32     `json:"characterLimit,omitempty"`
	InputType      InputType `json:"inputType,omitempty"`
}

// SelectionType is the format of the items users can select.
type SelectionType string

const (
	SelectionTypeCheckBox    SelectionType = "CHECK_BOX"
	SelectionTypeRadioButton SelectionType = "RADIO_BUTTON"
	SelectionTypeSwitch      SelectionType = "SWITCH"
	SelectionTypeDropdown    SelectionType = "DROPDOWN"
	SelectionTypeMultiSelect SelectionType = "MULTI_SELECT"
)

// SelectionItem is an item that users can select in a selection input.
type SelectionItem struct {
	Text         string `json:"text"`
	Value        string `json:"value"`
	Selected     bool   `json:"selected,omitempty"`
	StartIconURI string `json:"startIconUri,omitempty"`
	BottomText   string `json:"bottomText,omitempty"`
}

// CommonDataSource is a data source shared by all Google Workspace apps.
type CommonDataSource string

const (
	CommonDataSourceUnknown CommonDataSource = "UNKNOWN"
	CommonDataSourceUser    CommonDataSource = "USER"
)

// PlatformDataSource is a data source for a multi-select menu. Only one
// field may be set.
type PlatformDataSource struct {
	CommonDataSource CommonDataSource `json:"commonDataSource,omitempty"`
}

// SelectionInput is a widget that creates one or more UI items that users can
// select, such as checkboxes, radio buttons, switches or dropdown menus. At
// most one of ExternalDataSource and PlatformDataSource may be set.
type SelectionInput struct {
	Name           string          `json:"name"`
	Label          string          `json:"label,omitempty"`
	Type           SelectionType   `json:"type,omitempty"`
	Items          []SelectionItem `json:"items,omitempty"`
	OnChangeAction *Action         `json:"onChangeAction,omitempty"`

	MultiSelectMaxSelectedItems *int32 `json:"multiSelectMaxSelectedItems,omitempty"`
	MultiSelectMinQueryLength   int32  `json:"multiSelectMinQueryLength,omitempty"`

	ExternalDataSource *Action             `json:"externalDataSource,omitempty"`
	PlatformDataSource *PlatformDataSource `json:"platformDataSource,omitempty"`
}

// Validate checks the multi-select data source constraint.
func (s *SelectionInput) Validate() error {
	return oneOf("SelectionInput multiSelectDataSource", map[string]bool{
		"externalDataSource": s.ExternalDataSource != nil,
		"platformDataSource": s.PlatformDataSource != nil,
	})
}

// DateTimePickerType is the format of the date and time in a DateTimePicker.
type DateTimePickerType string

const (
	DateTimePickerTypeDateAndTime DateTimePickerType = "DATE_AND_TIME"
	DateTimePickerTypeDateOnly    DateTimePickerType = "DATE_ONLY"
	DateTimePickerTypeTimeOnly    DateTimePickerType = "TIME_ONLY"
)

// DateTimePicker lets users specify a date, a time, or both.
type DateTimePicker struct {
	Name  string             `json:"name"`
	Label string             `json:"label,omitempty"`
	Type  DateTimePickerType `json:"type,omitempty"`
	// ValueMsEpoch is the default value in milliseconds since the Unix epoch.
	// It is encoded as a JSON string, as for all protobuf int64 values.
	ValueMsEpoch *int64 `json:"valueMsEpoch,string,omitempty"`
	// TimezoneOffsetDate is the offset from UTC in minutes.
	TimezoneOffsetDate int32   `json:"timezoneOffsetDate,omitempty"`
	OnChangeAction     *Action `json:"onChangeAction,omitempty"`
}
