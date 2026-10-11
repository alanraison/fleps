package chat

import (
	"fmt"
	"strconv"

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
	email string,
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
							Subtitle: fmt.Sprintf("Season %v, Round %v", "2026/27", "26/27R5"),
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
												Value: "submitPredictions",
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
										Name:  string(fixture.HomeTeam),
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
										Name:  string(fixture.AwayTeam),
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
