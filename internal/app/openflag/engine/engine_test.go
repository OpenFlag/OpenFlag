package engine_test

import (
	"errors"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/stretchr/testify/suite"
)

type fakeLogger struct {
	called bool
}

func (f *fakeLogger) Log(result engine.Result) {
	f.called = true
}

type fakeFlagRepo struct {
	model.FlagRepo
	repoError bool
}

func (f *fakeFlagRepo) FindAll() ([]model.Flag, error) {
	if f.repoError {
		return nil, errors.New("fake flag repo error")
	}

	return []model.Flag{
		{
			ID:   10,
			Flag: "flag1",
			Segments: `
			[
				{
					"description": "segment 1",
					"constraints": {
						"A": {
							"name": "<",
							"parameters": {
								"value": 10
							}
						},
						"B": {
							"name": ">",
							"parameters": {
								"value": 5
							}
						}
					},
					"expression": "A ∩ B",
					"variant": {
						"variant_key": "on1"
					}
				},
				{
					"description": "segment 2",
					"constraints": {
						"A": {
							"name": "<",
							"parameters": {
								"value": 20
							}
						},
						"B": {
							"name": ">",
							"parameters": {
								"value": 15
							}
						}
					},
					"expression": "A ∩ B",
					"variant": {
						"variant_key": "on2"
					}
				}
			]
		`,
		},
		{
			ID:   11,
			Flag: "flag2",
			Segments: `
			[
				{
					"description": "segment 1",
					"constraints": {
						"A": {
							"name": "<",
							"parameters": {
								"value": 20
							}
						},
						"B": {
							"name": ">",
							"parameters": {
								"value": 15
							}
						}
					},
					"expression": "A ∩ B",
					"variant": {
						"variant_key": "on3"
					}
				}
			]
		`,
		},
	}, nil
}

type EngineSuite struct {
	suite.Suite
}

func (suite *EngineSuite) TestEngine() {
	cases := []struct {
		name        string
		repoError   bool
		evaluations []struct {
			flags  []string
			entity model.Entity
			result engine.Result
		}
	}{
		{
			name:      "successfully create engine and evaluate entities 1",
			repoError: false,
			evaluations: []struct {
				flags  []string
				entity model.Entity
				result engine.Result
			}{
				{
					flags: []string{"flag1", "flag2"},
					entity: model.Entity{
						EntityID: 7,
					},
					result: engine.Result{
						Entity: model.Entity{
							EntityID: 7,
						},
						Evaluations: []engine.Evaluation{
							{
								Flag: "flag1",
								Variant: model.Variant{
									VariantKey: "on1",
								},
							},
						},
					},
				},
				{
					entity: model.Entity{
						EntityID: 7,
					},
					result: engine.Result{
						Entity: model.Entity{
							EntityID: 7,
						},
						Evaluations: []engine.Evaluation{
							{
								Flag: "flag1",
								Variant: model.Variant{
									VariantKey: "on1",
								},
							},
						},
					},
				},
				{
					flags: []string{"flag1", "flag2"},
					entity: model.Entity{
						EntityID: 17,
					},
					result: engine.Result{
						Entity: model.Entity{
							EntityID: 17,
						},
						Evaluations: []engine.Evaluation{
							{
								Flag: "flag1",
								Variant: model.Variant{
									VariantKey: "on2",
								},
							},
							{
								Flag: "flag2",
								Variant: model.Variant{
									VariantKey: "on3",
								},
							},
						},
					},
				},
				{
					flags: []string{"flag1", "flag2"},
					entity: model.Entity{
						EntityID: 27,
					},
					result: engine.Result{
						Entity: model.Entity{
							EntityID: 27,
						},
						Evaluations: []engine.Evaluation{},
					},
				},
			},
		},
		{
			name:      "failed to create engine and evaluate entities 1",
			repoError: true,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			logger := &fakeLogger{}
			flagRepo := &fakeFlagRepo{}

			flagRepo.repoError = tc.repoError

			eng := engine.New(logger, flagRepo)

			err := eng.Fetch()
			if tc.repoError {
				suite.Error(err)
				return
			}

			suite.NoError(err)

			for _, evaluation := range tc.evaluations {
				result, err := eng.Evaluate(evaluation.flags, evaluation.entity)
				suite.NoError(err)
				suite.Equal(evaluation.result.Entity, result.Entity)
				suite.Equal(evaluation.result.Evaluations, result.Evaluations)
			}

			suite.Equal(true, logger.called)
		})
	}
}

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}
