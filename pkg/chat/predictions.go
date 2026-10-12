package chat

import (
	"fmt"
	"strconv"
	"strings"

	chatmodel "github.com/alanraison/fleps/pkg/chat/model"
	"github.com/alanraison/fleps/pkg/model"
)

type predictionsDialogBuilder struct {
	teamRepository model.TeamRepository
	actionUrl      string
}

func PredictionsDialogBuilder(
	teamRepository model.TeamRepository,
	actionUrl string,
) *predictionsDialogBuilder {
	return &predictionsDialogBuilder{
		teamRepository: teamRepository,
		actionUrl:      actionUrl,
	}
}

func (pdb *predictionsDialogBuilder) MakePredictionsDialog(
	roundID model.RoundID,
	fixtures []model.Fixture,
	predictions model.GamePredictions,
) *chatmodel.RenderActions {
	var sections []chatmodel.Section
	for _, fixture := range fixtures {
		sections = append(sections, pdb.matchSection(fixture, predictions))
	}
	return &chatmodel.RenderActions{
		Action: &chatmodel.RenderAction{
			Navigations: []chatmodel.Navigation{
				{
					PushCard: &chatmodel.Card{
						Header: &chatmodel.CardHeader{
							Title:    "Enter your predictions",
							Subtitle: fmt.Sprintf("Round %v", roundID),
						},
						Sections: sections,
						FixedFooter: &chatmodel.CardFixedFooter{
							PrimaryButton: &chatmodel.Button{
								Text: "Submit Predictions",
								OnClick: &chatmodel.OnClick{
									Action: &chatmodel.Action{
										Function: pdb.actionUrl,
										Parameters: []chatmodel.ActionParameter{
											{
												Key:   "actionName",
												Value: SubmitPredictionsAction,
											},
											{
												Key:   "roundId",
												Value: string(roundID),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (pdb *predictionsDialogBuilder) matchSection(fixture model.Fixture, prediction model.GamePredictions) chatmodel.Section {
	homeTeam, err := pdb.teamRepository.FindTeamByKey(fixture.HomeTeam)
	if err != nil {
		panic(fmt.Sprintf("unknown home team: %v", fixture.HomeTeam))
	}
	awayTeam, err := pdb.teamRepository.FindTeamByKey(fixture.AwayTeam)
	if err != nil {
		panic(fmt.Sprintf("unknown away team: %v", fixture.AwayTeam))
	}
	var homeGoals, awayGoals int
	if prediction != nil {
		prd := prediction[model.Game{
			HomeTeam: fixture.HomeTeam,
			AwayTeam: fixture.AwayTeam,
		}]
		homeGoals = prd.HomeGoals
		awayGoals = prd.AwayGoals
	}
	return chatmodel.Section{
		Header: fmt.Sprintf("%v v %v", homeTeam.FullName, awayTeam.FullName),
		Widgets: []chatmodel.Widget{
			{
				Columns: &chatmodel.Columns{
					ColumnItems: []chatmodel.Column{
						{
							HorizontalSizeStyle: chatmodel.HorizontalSizeStyleFillMinimumSpace,
							HorizontalAlignment: chatmodel.HorizontalAlignmentStart,
							VerticalAlignment:   chatmodel.ColumnVerticalAlignmentCenter,
							Widgets: []chatmodel.ColumnWidget{
								{
									TextInput: &chatmodel.TextInput{
										Name:  GoalsInputName(fixture.HomeTeam, fixture.AwayTeam, HomeSide),
										Label: string(homeTeam.ShortName),
										Value: strconv.Itoa(homeGoals),
									},
								},
							},
						},
						{
							HorizontalSizeStyle: chatmodel.HorizontalSizeStyleFillMinimumSpace,
							HorizontalAlignment: chatmodel.HorizontalAlignmentStart,
							VerticalAlignment:   chatmodel.ColumnVerticalAlignmentCenter,
							Widgets: []chatmodel.ColumnWidget{
								{
									TextInput: &chatmodel.TextInput{
										Name:  GoalsInputName(fixture.HomeTeam, fixture.AwayTeam, AwaySide),
										Label: string(awayTeam.ShortName),
										Value: strconv.Itoa(awayGoals),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Side identifies which team's goals a prediction input holds.
type Side string

const (
	HomeSide Side = "home"
	AwaySide Side = "away"
)

// GoalsInputName is the form input name for one side of a game's prediction.
func GoalsInputName(home, away model.TeamKey, side Side) string {
	return fmt.Sprintf("%s_%s_%s", home, away, side)
}

// ParseGoalsInputName is the inverse of GoalsInputName. ok is false if name
// is not a prediction input.
func ParseGoalsInputName(name string) (game model.Game, side Side, ok bool) {
	parts := strings.Split(name, "_")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return model.Game{}, "", false
	}
	side = Side(parts[2])
	if side != HomeSide && side != AwaySide {
		return model.Game{}, "", false
	}
	return model.Game{HomeTeam: model.TeamKey(parts[0]), AwayTeam: model.TeamKey(parts[1])}, side, true
}
