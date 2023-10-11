package main

import (
	"encoding/xml"
)

type GilabVersion struct {
	Version string `json:"version"`
}

type GitlabUpdateStatus struct {
	XMLName xml.Name `xml:"svg"`
	Text    struct {
		Text string `xml:",chardata"`
	} `xml:"text"`
}

type DiskInfo struct {
	Size uint64 `json:"Size"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

type AllFiles struct {
	Name string
	Date string
}

type LivenessProbe struct {
	Status string `json:"status"`
}

type ReadinessProbe struct {
	Status      string `json:"status"`
	MasterCheck []struct {
		Status string `json:"status"`
	} `json:"master_check"`
	DbCheck []struct {
		Status string `json:"status"`
	} `json:"db_check"`
	CacheCheck []struct {
		Status string `json:"status"`
	} `json:"cache_check"`
	QueuesCheck []struct {
		Status string `json:"status"`
	} `json:"queues_check"`
	RateLimitingCheck []struct {
		Status string `json:"status"`
	} `json:"rate_limiting_check"`
	SessionsCheck []struct {
		Status string `json:"status"`
	} `json:"sessions_check"`
	SharedStateCheck []struct {
		Status string `json:"status"`
	} `json:"shared_state_check"`
	TraceChunksCheck []struct {
		Status string `json:"status"`
	} `json:"trace_chunks_check"`
	GitalyCheck []struct {
		Status string `json:"status"`
	} `json:"gitaly_check"`
}

type Labels struct {
	Trigger        string `json:"trigger"`
	Project        string `json:"project"`
	Hostname       string `json:"hostname"`
	IP             string `json:"ip"`
	Severity_level string `json:"severity_level"`
}

type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type Alert struct {
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}
