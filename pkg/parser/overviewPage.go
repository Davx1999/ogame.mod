package parser

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
)

func (p OverviewPage) ExtractActiveItems() ([]ogame.ActiveItem, error) {
	return p.e.ExtractActiveItems(p.content)
}

func (p OverviewPage) ExtractDMCosts() (ogame.DMCosts, error) {
	return p.e.ExtractDMCosts(p.content)
}

func (p OverviewPage) ExtractConstructions() (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64, ogame.ID, int64) {
	return p.e.ExtractConstructions(p.content)
}

func (p OverviewPage) ExtractUserInfos() (ogame.UserInfos, error) {
	return p.e.ExtractUserInfos(p.content)
}

func (p OverviewPage) ExtractCancelResearchInfos() (token string, techID, listID int64, err error) {
	return p.e.ExtractCancelResearchInfos(p.content)
}

func (p OverviewPage) ExtractCancelBuildingInfos() (token string, techID, listID int64, err error) {
	return p.e.ExtractCancelBuildingInfos(p.content)
}

func (p OverviewPage) ExtractCancelLfBuildingInfos() (token string, id, listID int64, err error) {
	return p.e.ExtractCancelLfBuildingInfos(p.content)
}

func (p OverviewPage) ExtractOverviewProduction() ([]ogame.Quantifiable, int64, error) {
	return p.e.ExtractOverviewProduction(p.content)
}
