package loadtest

import (
	"bedouin/bedouin/session"
	stats "bedouin/bedouin/tracing"
	"fmt"
	"net/http"
)

type SimpleGetConstantLoadTest struct {
	hs          *session.HttpSession
	endPointUrl string
}

func NewSimpleGetConstantLoadTest(endPointUrl string) (*SimpleGetConstantLoadTest, error) {
	return &SimpleGetConstantLoadTest{
		hs:          session.DefaultHttpSession,
		endPointUrl: endPointUrl,
	}, nil
}

func (t *SimpleGetConstantLoadTest) Send() {
	req, err := http.NewRequest("GET", t.endPointUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := t.hs.Submit(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body:", err)
	//	return
	//}
	//fmt.Println("Body:", string(body))

	defer resp.Body.Close()
}

func (t *SimpleGetConstantLoadTest) GetAggStats() *stats.AggStats {
	return t.hs.GetAggStats()
}

func (t *SimpleGetConstantLoadTest) GetPrintableAggStats() map[string]any {
	aggStats := *t.GetAggStats()
	return structToMap(aggStats)
}
