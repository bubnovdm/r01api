package internal

type Domains struct {
	Content struct {
		Data []struct {
			ID     int    `json:"id"`
			Domain string `json:"domain"`
		} `json:"data"`
	} `json:"content"`
}

type RRecords struct {
	Content struct {
		Data []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			Ttl  int    `json:"ttl"`
		}
	}
}

type AddRecordResponse struct {
	Content struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	} `json:"content"`
}
