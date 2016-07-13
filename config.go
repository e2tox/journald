package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type Config struct {
	AWSRegion      string
	EC2InstanceId  string
	LogGroupName   string
	LogStreamName  string
	LogPriority    Priority
	StateFilename  string
	BufferSize     int
}

type fileConfig struct {
	AWSRegion     string `hcl:"aws_region"`
	EC2InstanceId string `hcl:"ec2_instance_id"`
	LogGroupName  string `hcl:"log_group"`
	LogStreamName string `hcl:"log_stream"`
	LogPriority   string `hcl:"log_priority"`
	StateFilename string `hcl:"state_file"`
	BufferSize    int    `hcl:"buffer_size"`
}

func getLogLevel(priority string) (Priority, error) {

	logLevels := map[Priority][]string{
		EMERGENCY: {"0", "emerg"},
		ALERT: {"1", "alert"},
		CRITICAL: {"2", "crit"},
		ERROR: {"3", "err"},
		WARNING: {"4", "warning"},
		NOTICE: {"5", "notice"},
		INFO: {"6", "info"},
		DEBUG: {"7", "debug"},
	}

	for i, s := range logLevels {
		if s[0] == priority || s[1] == priority {
			return i, nil
		}
	}

	return DEBUG, fmt.Errorf("'%s' is unsupported log priority", priority)
}

func LoadConfig(filename string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var fConfig fileConfig
	err = hcl.Decode(&fConfig, string(configBytes))
	if err != nil {
		return nil, err
	}

	if fConfig.LogGroupName == "" {
		return nil, fmt.Errorf("log_group is required")
	}
	if fConfig.StateFilename == "" {
		return nil, fmt.Errorf("state_file is required")
	}

	config := &Config{}

	config.LogGroupName = fConfig.LogGroupName

	if fConfig.LogStreamName != "" {
		config.LogStreamName = fConfig.LogStreamName
	} else {
		// By default we use the instance id as the stream name.
		config.LogStreamName = config.EC2InstanceId
	}

	config.StateFilename = fConfig.StateFilename

	if fConfig.BufferSize != 0 {
		config.BufferSize = fConfig.BufferSize
	} else {
		config.BufferSize = 100
	}

	return config, nil
}
