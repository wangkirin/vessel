package pipeline

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"
)

var (
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d/stage/"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/plv-%d/stagev-%d/"
)

// RunPipeline : run pipeline generate pipelineVersion
func RunPipeline(pl *models.Pipeline) (*models.PipelineVersion, error) {
	// first test is pipeline legal if not return err
	relationMap, err := isPipelineLegal(pl)
	if err != nil {
		return nil, err
	}

	etcd.SavePipelineInfo(pl)
	for _, stage := range pl.Stages {
		if relationMap[stage.Name][0] != "" {
			stage.From = strings.Split(relationMap[stage.Name][0], ",")
		}
		if relationMap[stage.Name][1] != "" {
			stage.To = strings.Split(relationMap[stage.Name][1], ",")
		}
		pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pl.WorkspaceId, pl.ProjectId, pl.Id)
		stagePath := pipelinePath + stage.Name
		etcd.SaveStageInfo(stage, stagePath)
	}

	pipelineVersion := new(models.PipelineVersion)
	pipelineVersion.Id = time.Now().UnixNano()
	pipelineVersion.WorkspaceId = pl.WorkspaceId
	pipelineVersion.ProjectId = pl.ProjectId
	pipelineVersion.PipelineId = pl.Id
	pipelineVersion.Namespace = "plv" + "-" + strconv.FormatInt(pipelineVersion.Id, 10)
	pipelineVersion.SelfLink = ""
	pipelineVersion.Created = time.Now().Unix()
	pipelineVersion.Updated = time.Now().Unix()
	pipelineVersion.Labels = pl.Labels
	pipelineVersion.Annotations = pl.Annotations
	pipelineVersion.Detail = pl.Detail
	pipelineVersion.StageVersions = []string{strconv.FormatInt(pipelineVersion.Id, 10)}
	pipelineVersion.Status = 0

	stageVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stageVersionPath = stageVersionPath[:len(stageVersionPath)-1]
	stageVersionPath = stageVersionPath[:strings.LastIndex(stageVersionPath, "/")] + "/pipelineId"
	etcd.SavePipelineId(stageVersionPath, strconv.FormatInt(pl.Id, 10))

	return pipelineVersion, nil
}

// test is the given pipeline is legal ,if legal return pipeline's stage relationMap if not return error
func isPipelineLegal(pipeline *models.Pipeline) (map[string][]string, error) {
	stageMap := make(map[string]*models.Stage, 0)
	dependenceCount := make(map[string]int, 0)
	stageRelationMap := make(map[string][]string, 0)

	// regist all stage,and check repeat/nil stage name
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("stage has a nil name")
		}
		if _, exist := stageMap[stage.Name]; !exist {
			stageMap[stage.Name] = stage
		} else {
			// has a repeat stage name ,return
			return nil, errors.New("stage has repeat name:" + stage.Name)
		}
	}

	// init stage dependence count
	for stageName, _ := range stageMap {
		dependenceCount[stageName] = 0
	}

	// count stage dependence
	for _, stage := range stageMap {
		for _, from := range stage.From {
			dependenceCount[from]++
		}
	}

	// check DAG
	//if AnnulusTag == nowReleaseStageCount or nowReleaseStageCount == len(dependenceCount) then exit for,if nowReleaseStageCount == len(dependenceCount) then isDAG,else isNotDAG
	nowReleaseStageCount := 0
	for true {

		annulusTag := 0
		for stageName, stage := range stageMap {
			if dependenceCount[stageName] == 0 {
				nowReleaseStageCount++
				for _, from := range stage.From {
					dependenceCount[from]--
				}

				dependenceCount[stage.Name] = -1
			} else if dependenceCount[stageName] == -1 {
				annulusTag++
			}
		}

		if annulusTag == nowReleaseStageCount || nowReleaseStageCount == len(dependenceCount) {
			break
		}
	}

	if nowReleaseStageCount != len(dependenceCount) {
		return nil, errors.New("given pipeline's stage can't create a DAG")
	}

	// generate stage relationMap
	// stageRelationMap := map[stageName]{"stage.From","stage.To"}
	for stageName, stage := range stageMap {
		if _, exist := stageRelationMap[stageName]; !exist {
			stageRelationMap[stageName] = make([]string, 2)
		}
		stageRelationMap[stageName][0] = strings.Join(stage.From, ",")

		for _, from := range stage.From {
			if _, exist := stageRelationMap[from]; !exist {
				stageRelationMap[from] = make([]string, 2)
			}
			stageRelationMap[from][1] = strings.Join(append(strings.Split(stageRelationMap[from][1], ","), stageName), ",")
		}
	}

	return stageRelationMap, nil
}