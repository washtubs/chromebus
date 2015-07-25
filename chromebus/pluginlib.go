package chromebus

type Plugin struct {
	Init    func(input chan *ChromebusRecord)
	Cleanup func()
}

func Focus(id string) {

}

func NewTab() {

}

func CloseTab(id string) {

}

func Navigate(id string, url string) {

}
