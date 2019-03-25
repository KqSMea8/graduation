package jobs

import (
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
)

const JobName_StoreFile = "store_file"

type StoreFileJob struct {
	fid     int64
	content []byte
	storage store.Storage
}

func NewStoreFileJob(fid int64, content []byte, storage store.Storage) *StoreFileJob {
	j := &StoreFileJob{
		fid:     fid,
		content: content,
		storage: storage,
	}
	return j
}

func (j *StoreFileJob) GetName() string {
	return JobName_SaveFileMeta
}

func (j *StoreFileJob) Run() (interface{}, error) {
	if err := j.storage.Write(j.fid, j.content); err != nil {
		logrus.Errorf("Persistent File Error: %s", err)
		return nil, err
	}
	return nil, nil
}
