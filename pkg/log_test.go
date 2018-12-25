package root_test

import (
	"testing"
	"xiaoyun/pkg"
)

func TestLog_StdOut(t *testing.T) {

	log := root.NewLogStdOut()

	log.Logger.Info("std_Out_test")

}

func TestLog_FileOut(t *testing.T) {

	log := root.NewLogFileOut("test.log")
	log.Logger.Info("file_out_test")

}

func TestLog_MultiOut(t *testing.T) {

	log := root.NewLogMultiOut("test.log")
	log.Logger.Info("multi_out_test")

}
