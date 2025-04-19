package masscan

type Report struct {
	Ranges  []string `json:"ranges"`
	Ports   []string `json:"ports"`
	MaxRate int      `json:"max_rate"`

	Results    map[string]Results `json:"results"`
	RawResults RawResults         `json:"raw_results"`
}

type Results struct {
	IP    string `json:"ip"`
	Ports Ports  `json:"ports"`
}

type RawResults []RawResult

type RawResult struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     Ports  `json:"ports"`
}

type Ports []Port

type Port struct {
	Port   int    `json:"port"`
	Proto  string `json:"proto"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	TTL    int    `json:"ttl"`
}
