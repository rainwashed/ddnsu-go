package global

type ManagedRecord struct {
	Record DDNSURecord
	Action string
}

type DDNSURecord struct {
	Name    string
	Comment string
	Ttl     int
	Content string
	Type    string
	Id      string
}

type CloudflareZoneResultArray struct {
	Id   string
	Name string
}

type CloudflareZoneResponse struct {
	Success bool
	Result  []CloudflareZoneResultArray
}

type CloudflareZoneRecordResult struct {
	Comment    string                 `json:"comment,omitempty"`
	Content    string                 `json:"content"`
	CreatedOn  string                 `json:"created_on"`
	Id         string                 `json:"id"`
	Meta       map[string]interface{} `json:"meta"`
	ModifiedOn string                 `json:"modified_on"`
	Name       string                 `json:"name"`
	Proxiable  bool                   `json:"proxiable"`
	Proxied    bool                   `json:"proxied"`
	Setttings  map[string]interface{} `json:"settings,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	Ttl        int                    `json:"ttl"`
	Type       string                 `json:"type"`
	ZoneId     string                 `json:"zone_id"`
	ZoneName   string                 `json:"zone_name"`
	Priority   *int                   `json:"priority,omitempty"`
}

type CloudflareZoneRecordResponseMulti struct {
	Errors     []interface{}                `json:"errors"`
	Messages   []interface{}                `json:"messages"`
	Result     []CloudflareZoneRecordResult `json:"result"`
	ResultInfo map[string]interface{}       `json:"result_info"`
	Success    bool                         `json:"success"`
}
type CloudflareZoneRecordResponseSingle struct {
	Errors     []interface{}              `json:"errors"`
	Messages   []interface{}              `json:"messages"`
	Result     CloudflareZoneRecordResult `json:"result"`
	ResultInfo map[string]interface{}     `json:"result_info"`
	Success    bool                       `json:"success"`
}

type SerializedRecords struct {
	Comment string
	Content string
	Id      string
	Name    string
	Type    string
	Ttl     interface{}
	Managed bool
}
type SerializedDNSState struct {
	Records   []SerializedRecords
	Provider  string
	PrimaryIP string
}
