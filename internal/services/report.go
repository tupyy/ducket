package services

import "git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"

type ReportService struct {
	dt *pg.Datastore
}

func NewReportService(dt *pg.Datastore) *ReportService {
	return &ReportService{dt: dt}
}

