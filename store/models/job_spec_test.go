package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestJobSpec_Save(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1, initr := cltest.NewJobSpecWithSchedule("* * * * 7")
	assert.Nil(t, store.SaveJob(&j1))

	store.Save(j1)
	j2, err := store.FindJob(j1.ID)
	assert.Nil(t, err)
	assert.Equal(t, initr.Schedule, j2.Initiators[0].Schedule)
}

func TestJobSpec_NewRun(t *testing.T) {
	t.Parallel()

	j, initr := cltest.NewJobSpecWithSchedule("1 * * * *")
	j.Tasks = []models.TaskSpec{cltest.NewTask("NoOp", `{"a":1}`)}

	run := j.NewRun(initr)

	assert.Equal(t, j.ID, run.JobID)
	assert.Equal(t, 1, len(run.TaskRuns))

	taskRun := run.TaskRuns[0]
	assert.Equal(t, "NoOp", taskRun.Task.Type)
	adapter, _ := adapters.For(taskRun.Task, nil)
	assert.NotNil(t, adapter)
	assert.JSONEq(t, `{"type":"NoOp","a":1}`, taskRun.Task.Params.String())

	assert.Equal(t, initr, run.Initiator)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.Time{Valid: false}, endAt.Time, false},
		{"before end at", endAt, endAt.Time.Add(-time.Nanosecond), false},
		{"at end at", endAt, endAt.Time, false},
		{"after end at", endAt, endAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j := cltest.NewJobSpec()
			j.EndAt = test.endAt

			assert.Equal(t, test.want, j.Ended(test.current))
		})
	}
}

func TestJobSpec_Started(t *testing.T) {
	t.Parallel()

	startAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		startAt null.Time
		current time.Time
		want    bool
	}{
		{"no start at", null.Time{Valid: false}, startAt.Time, true},
		{"before start at", startAt, startAt.Time.Add(-time.Nanosecond), false},
		{"at start at", startAt, startAt.Time, true},
		{"after start at", startAt, startAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j := cltest.NewJobSpec()
			j.StartAt = test.startAt

			assert.Equal(t, test.want, j.Started(test.current))
		})
	}
}

func TestTask_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		taskType      string
		confirmations uint64
		json          string
	}{
		{"noop", "noop", 0, `{"type":"NoOp"}`},
		{"httpget", "httpget", 0, `{"type":"httpget","url":"http://www.no.com"}`},
		{"with confirmations", "noop", 10, `{"type":"noOp","confirmations":10}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var task models.TaskSpec
			err := json.Unmarshal([]byte(test.json), &task)
			assert.Nil(t, err)
			assert.Equal(t, test.confirmations, task.Confirmations)

			assert.Equal(t, test.taskType, task.Type)
			_, err = adapters.For(task, nil)
			assert.Nil(t, err)

			s, err := json.Marshal(task)
			assert.Nil(t, err)
			assert.Equal(t, test.json, string(s))
		})
	}
}
