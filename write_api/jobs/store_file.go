package jobs

import (
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
	"io"
)

const JobName_StoreFile = "store_file"

type StoreFileJob struct {
	fid     int64
	reader  io.Reader
	storage store.Storage
}

func NewStoreFileJob(fid int64, reader io.Reader, storage store.Storage) *StoreFileJob {
	j := &StoreFileJob{
		fid:     fid,
		reader: reader,
		storage: storage,
	}
	return j
}

func (j *StoreFileJob) GetName() string {
	return JobName_SaveFileMeta
}

func (j *StoreFileJob) Run() (interface{}, error) {
	if err := j.storage.Write(j.fid, j.reader); err != nil {
		logrus.Errorf("Persistent File Error: %s", err)
		return nil, err
	}
	return nil, nil
}
