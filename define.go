package tupu

const (
	VERSION = "1.0"
)

type ResponseCapsule struct {
	Json string `json:"json"`
}

type Response struct {
	Code         int            `json:"code"`
	Message      string         `json:"message"`
	Timestamp    string         `json:"timestamp"`
	Nonce        string         `json:"nonce"`
	FileList     []ResponseFile `json:"fileList"`
	Statistic    []int          `json:"statistic"`
	CallRecordId string         `json:"callRecordId"`
	Signature    string         `json:"signature"`
}

type ResponseFile struct {
	Rate  float64 `json:"rate"`
	Label int     `json:"label"`
	Name  string  `json:"name"`
}
