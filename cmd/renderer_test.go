package cmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
)

func TestRendererJSONRenderJobs(t *testing.T) {
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	j := cltest.NewJobSpec()
	jobSpecs := []models.JobSpec{j}
	assert.Nil(t, r.Render(&jobSpecs))
}

func TestRendererTableRenderJobs(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	j := cltest.NewJobSpec()
	jobSpecs := []models.JobSpec{j}
	assert.Nil(t, r.Render(&jobSpecs))
}

func TestRendererTableRenderShowJob(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	j, initr := cltest.NewJobSpecWithWebInitiator()
	run := j.NewRun(initr)
	p := presenters.JobSpec{JobSpec: j, Runs: []models.JobRun{run}}
	assert.Nil(t, r.Render(&p))
}

func TestRendererTableRenderUnknown(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.NotNil(t, r.Render(&anon))
}
