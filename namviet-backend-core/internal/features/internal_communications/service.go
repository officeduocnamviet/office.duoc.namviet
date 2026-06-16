package internal_communications

// Internal Channels
func GetAllInternalChannelsService() ([]InternalChannel, error) {
	return GetAllInternalChannels()
}

func GetInternalChannelByIDService(id int64) (*InternalChannel, error) {
	return GetInternalChannelByID(id)
}

func CreateInternalChannelService(req CreateInternalChannelRequest) (*InternalChannel, error) {
	channelType := "group"
	if req.Type != nil {
		channelType = *req.Type
	}

	channel := &InternalChannel{
		Name: req.Name,
		Type: channelType,
	}

	if err := CreateInternalChannel(channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func UpdateInternalChannelService(id int64, req UpdateInternalChannelRequest) (*InternalChannel, error) {
	channel, err := GetInternalChannelByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		channel.Name = *req.Name
	}
	if req.Type != nil {
		channel.Type = *req.Type
	}

	if err := UpdateInternalChannel(channel); err != nil {
		return nil, err
	}
	return channel, nil
}

func DeleteInternalChannelService(id int64) error {
	return DeleteInternalChannel(id)
}

// Internal Messages
func GetMessagesByChannelIDService(channelID int64) ([]InternalMessage, error) {
	return GetMessagesByChannelID(channelID)
}

func CreateInternalMessageService(req CreateInternalMessageRequest) (*InternalMessage, error) {
	message := &InternalMessage{
		ChannelID: req.ChannelID,
		SenderID:  req.SenderID,
		Content:   req.Content,
	}

	if err := CreateInternalMessage(message); err != nil {
		return nil, err
	}
	return message, nil
}

func DeleteInternalMessageService(id int64) error {
	return DeleteInternalMessage(id)
}
