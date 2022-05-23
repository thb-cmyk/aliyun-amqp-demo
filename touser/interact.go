package touser

type Message struct {
	Id string

	Content []byte // the content is a uri typed link. example: /acount/product/device/property ...
}

type URI struct {
	Acount    string
	Product   string
	Device    string
	Property  string
	Event     string
	StartTiem string
	EndTIme   string
}

func Message_Create(id string, content []byte) *Message {
	message := new(Message)

	message.Id = id
	message.Content = content

	return message
}

func (m *Message) Message_Get_Id() string {
	return m.Id
}

func (m *Message) Message_Get_Content() []byte {
	return m.Content
}

func URI_Property_Assemble(acount string, product string, device string, property string, starttime string, endtime string) *URI {
	uri := new(URI)

	uri.Acount = acount
	uri.Product = product
	uri.Device = device
	uri.Property = property
	uri.StartTiem = starttime
	uri.EndTIme = endtime

	return uri
}

func URI_Event_Assemble(acount string, product string, device string, event string, starttime string, endtime string) *URI {
	uri := new(URI)

	uri.Acount = acount
	uri.Product = product
	uri.Device = device
	uri.Event = event
	uri.StartTiem = starttime
	uri.EndTIme = endtime

	return uri
}
