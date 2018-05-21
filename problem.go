package moonapi

type MoonBoardConfiguration struct {
	ID          int         `json:"Id"`
	Description string      `json:"Description"`
	LowGrade    interface{} `json:"LowGrade"`
	HighGrade   interface{} `json:"HighGrade"`
}

type Setter struct {
	ID              string `json:"Id"`
	Nickname        string `json:"Nickname"`
	Firstname       string `json:"Firstname"`
	Lastname        string `json:"Lastname"`
	City            string `json:"City"`
	Country         string `json:"Country"`
	ProfileImageURL string `json:"ProfileImageUrl"`
	CanShareData    bool   `json:"CanShareData"`
	
}

type HoldSetup struct {
	ID                      int         `json:"Id"`
	Description             string      `json:"Description"`
	Setby                   interface{} `json:"Setby"`
	DateInserted            interface{} `json:"DateInserted"`
	DateUpdated             interface{} `json:"DateUpdated"`
	DateDeleted             interface{} `json:"DateDeleted"`
	IsLocked                bool        `json:"IsLocked"`
	Holdsets                interface{} `json:"Holdsets"`
	MoonBoardConfigurations interface{} `json:"MoonBoardConfigurations"`
	HoldLayoutID            int         `json:"HoldLayoutId"`
	AllowClimbMethods       bool        `json:"AllowClimbMethods"`
}

type Problem struct {
		Method                 string      `json:"Method"`
		Name                   string      `json:"Name"`
		Grade                  string      `json:"Grade"`
		UserGrade              interface{} `json:"UserGrade"`
		MoonBoardConfiguration MoonBoardConfiguration `json:"MoonBoardConfiguration"`
		MoonBoardConfigurationID int `json:"MoonBoardConfigurationId"`
		Setter	Setter `json:"Setter"`
		FirstAscender bool `json:"FirstAscender"`
		Rating        int  `json:"Rating"`
		UserRating    int  `json:"UserRating"`
		Repeats       int  `json:"Repeats"`
		Attempts      int  `json:"Attempts"`
		Holdsetup     HoldSetup `json:"Holdsetup"`
		IsBenchmark bool `json:"IsBenchmark"`
		Moves       []struct {
			ID          int    `json:"Id"`
			Description string `json:"Description"`
			IsStart     bool   `json:"IsStart"`
			IsEnd       bool   `json:"IsEnd"`
		} `json:"Moves"`
		Holdsets  interface{} `json:"Holdsets"`
		Locations []struct {
			ID              int         `json:"Id"`
			Holdset         interface{} `json:"Holdset"`
			Description     interface{} `json:"Description"`
			X               int         `json:"X"`
			Y               int         `json:"Y"`
			Color           string      `json:"Color"`
			Rotation        int         `json:"Rotation"`
			Type            int         `json:"Type"`
			HoldNumber      interface{} `json:"HoldNumber"`
			Direction       int         `json:"Direction"`
			DirectionString string      `json:"DirectionString"`
		} `json:"Locations"`
		RepeatText     string      `json:"RepeatText"`
		NumberOfTries  interface{} `json:"NumberOfTries"`
		NameForURL     string      `json:"NameForUrl"`
		ID             int         `json:"Id"`
		APIID          int         `json:"ApiId"`
		DateInserted   string      `json:"DateInserted"`
		DateUpdated    interface{} `json:"DateUpdated"`
		DateDeleted    interface{} `json:"DateDeleted"`
		DateTimeString string      `json:"DateTimeString"`
	}