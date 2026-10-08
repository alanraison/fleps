package model

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestCardJSONRoundTrip(t *testing.T) {
	epoch := int64(1700000000000)
	alpha := float32(0.5)
	maxSel := int32(3)
	card := Card{
		Header:              &CardHeader{Title: "Title", ImageType: ImageTypeCircle},
		SectionDividerStyle: DividerStyleSolid,
		Sections: []Section{{
			Header: "s",
			Widgets: []Widget{
				{TextParagraph: &TextParagraph{Text: "hi", TextSyntax: TextSyntaxMarkdown}},
				{Divider: &Divider{}},
				{DecoratedText: &DecoratedText{
					Text:   "x",
					Button: &Button{Text: "go", Color: &Color{Red: 1, Alpha: &alpha}, OnClick: &OnClick{OpenLink: &OpenLink{URL: "https://example.com"}}},
				}},
				{SelectionInput: &SelectionInput{Name: "n", Type: SelectionTypeMultiSelect, MultiSelectMaxSelectedItems: &maxSel,
					PlatformDataSource: &PlatformDataSource{CommonDataSource: CommonDataSourceUser}}},
				{DateTimePicker: &DateTimePicker{Name: "d", ValueMsEpoch: &epoch}},
				{Columns: &Columns{ColumnItems: []Column{{Widgets: []ColumnWidget{{TextParagraph: &TextParagraph{Text: "c"}}}}}}},
				{Carousel: &Carousel{CarouselCards: []CarouselCard{{Widgets: []NestedWidget{{Image: &Image{ImageURL: "u"}}}}}}},
			},
		}},
		FixedFooter: &CardFixedFooter{PrimaryButton: &Button{Text: "ok"}},
	}
	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"valueMsEpoch":"1700000000000"`) {
		t.Errorf("int64 should be a JSON string: %s", b)
	}
	if !strings.Contains(string(b), `"divider":{}`) {
		t.Errorf("empty divider must be emitted: %s", b)
	}
	var got Card
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(card, got) {
		t.Errorf("round trip mismatch:\nwant %+v\ngot  %+v", card, got)
	}
	if err := card.Validate(); err != nil {
		t.Errorf("valid card failed validation: %v", err)
	}
}

func TestUnmarshalAPIJSON(t *testing.T) {
	in := `{"header":{"title":"T"},"sections":[{"widgets":[{"buttonList":{"buttons":[{"text":"b","onClick":{"action":{"function":"f","parameters":[{"key":"k","value":"v"}]}}}]}}]}]}`
	var c Card
	if err := json.Unmarshal([]byte(in), &c); err != nil {
		t.Fatal(err)
	}
	a := c.Sections[0].Widgets[0].ButtonList.Buttons[0].OnClick.Action
	if a.Function != "f" || a.Parameters[0].Key != "k" {
		t.Errorf("unexpected action: %+v", a)
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		card    Card
		wantErr string
	}{
		"two widget types": {
			card:    Card{Sections: []Section{{Widgets: []Widget{{TextParagraph: &TextParagraph{}, Divider: &Divider{}}}}}},
			wantErr: "[divider textParagraph]",
		},
		"two onClick handlers": {
			card:    Card{CardActions: []CardAction{{OnClick: &OnClick{Action: &Action{}, OpenLink: &OpenLink{}}}}},
			wantErr: "OnClick",
		},
		"two decorated text controls": {
			card:    Card{Sections: []Section{{Widgets: []Widget{{DecoratedText: &DecoratedText{Button: &Button{}, EndIcon: &Icon{}}}}}}},
			wantErr: "DecoratedText control",
		},
		"two icon sources": {
			card:    Card{FixedFooter: &CardFixedFooter{PrimaryButton: &Button{Icon: &Icon{KnownIcon: "STAR", IconURL: "u"}}}},
			wantErr: "Icon",
		},
		"two data sources": {
			card:    Card{Sections: []Section{{Widgets: []Widget{{SelectionInput: &SelectionInput{ExternalDataSource: &Action{}, PlatformDataSource: &PlatformDataSource{}}}}}}},
			wantErr: "multiSelectDataSource",
		},
		"nested card": {
			card:    Card{CardActions: []CardAction{{OnClick: &OnClick{Card: &Card{Sections: []Section{{Widgets: []Widget{{Divider: &Divider{}, Grid: &Grid{}}}}}}}}}},
			wantErr: "validating nested card",
		},
		"column widget": {
			card:    Card{Sections: []Section{{Widgets: []Widget{{Columns: &Columns{ColumnItems: []Column{{Widgets: []ColumnWidget{{ChipList: &ChipList{}, Image: &Image{}}}}}}}}}}},
			wantErr: "ColumnWidget",
		},
		"carousel": {
			card:    Card{Sections: []Section{{Widgets: []Widget{{Carousel: &Carousel{CarouselCards: []CarouselCard{{FooterWidgets: []NestedWidget{{Image: &Image{}, ButtonList: &ButtonList{}}}}}}}}}}},
			wantErr: "NestedWidget",
		},
		"empty card": {card: Card{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.card.Validate()
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatal("expected error, got nil")
			case err != nil && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("error %q does not contain %q", err, tc.wantErr)
			}
		})
	}
}
