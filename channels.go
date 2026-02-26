package channels

import (
	"fmt"
	"strings"
	"time"

	"github.com/RAGFishAI/channels/migration"
)

// Channelsetup used to initialie channel configuration
func ChannelSetup(config Config) *Channel {

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Channel{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
	}

}

// get all channel list
func (channel *Channel) ListChannel(inputs Channels) (channelList []Tblchannel, channelcount int, err error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {
		return []Tblchannel{}, 0, autherr
	}

	CH.Userid = channel.Userid
	CH.Dataaccess = channel.DataAccess

	var (
		channellist []Tblchannel
		count       int64
	)

	err = CH.Channellist(channel.DB, channel, inputs, &channellist, &count)

	if err != nil {
		return []Tblchannel{}, 0, err
	}

	if inputs.ChannelFile {

		for i := range channellist {
			files, err := CH.GetFilesByChannelId(channellist[i].Id, channel.DB, inputs.TenantId)
			if err != nil {
				return []Tblchannel{}, 0, err
			}
			channellist[i].FilesData = files
		}
	}
	return channellist, int(count), nil
}

func (channel *Channel) ChannelDetail(inputs Channels) (channelDetails Tblchannel, err error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {
		return Tblchannel{}, autherr
	}

	CH.Userid = channel.Userid
	CH.Dataaccess = channel.DataAccess

	if err = CH.ChannelDetail(channel.DB, inputs, &channelDetails); err != nil {

		return Tblchannel{}, err
	}

	return channelDetails, nil
}

/*create channel*/
func (channel *Channel) CreateChannel(channelcreate ChannelCreate, tenantid string) (TblChannel, error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return TblChannel{}, autherr
	}

	/*create channel*/
	var cchannel TblChannel
	cchannel.ChannelName = channelcreate.ChannelName
	cchannel.ChannelUniqueId = channelcreate.ChannelUniqueId
	cchannel.ChannelDescription = channelcreate.ChannelDescription
	cchannel.SlugName = strings.ToLower(strings.ReplaceAll(channelcreate.SlugName, " ", "-"))
	cchannel.IsActive = 1
	cchannel.CreatedBy = channelcreate.CreatedBy
	cchannel.TenantId = tenantid
	cchannel.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	ch, chanerr := CH.CreateChannel(&cchannel, channel.DB)

	if chanerr != nil {

		fmt.Println(chanerr)
	}

	return ch, nil
}

/*Get channel by name*/
func (channel *Channel) GetchannelByName(channelname string, tenantid string) (channels Tblchannel, err error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return Tblchannel{}, autherr
	}

	channellist, err1 := CH.GetChannelByChannelName(channelname, channel.DB, tenantid)

	if err1 != nil {

		return Tblchannel{}, err1
	}

	return channellist, nil

}

/*Get Channels By Id*/
func (channel *Channel) GetChannelsById(channelid int, tenantid string) (channelList Tblchannel, err error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return Tblchannel{}, autherr
	}

	channellist, err := CH.GetChannelById(channelid, channel.DB, tenantid)

	if err != nil {

		return Tblchannel{}, err

	}

	return channellist, nil
}

/*Delete Channel*/
func (channel *Channel) DeleteChannel(channelid, modifiedby int, tenantid string) error {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return autherr
	}

	if channelid <= 0 {

		return ErrorChannelId
	}

	CH.DeleteChannelById(channelid, channel.DB, tenantid)

	return nil

}

/*Change Channel status*/
// status 0 = inactive
// status 1 = active
func (channel *Channel) ChangeChannelStatus(channelid int, status, modifiedby int, tenantid string) (bool, error) {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return false, autherr
	}

	if channelid <= 0 {

		return false, ErrorChannelId
	}

	var channelstatus TblChannel

	channelstatus.ModifiedBy = modifiedby

	channelstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	CH.ChannelIsActive(&channelstatus, channelid, status, channel.DB, tenantid)

	return true, nil

}

/*Edit channel*/

func (channel *Channel) EditChannel(ChannelName string, channeluniqueid string, channelslugname string, ChannelDescription string, modifiedby int, channelid int, tenantid string) error {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return autherr
	}

	var chn TblChannel

	chn.ChannelName = ChannelName

	chn.ChannelUniqueId = channeluniqueid

	chn.ChannelDescription = ChannelDescription

	chn.SlugName = strings.ReplaceAll(strings.ToLower(channelslugname), " ", "-")

	chn.ModifiedBy = modifiedby

	chn.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	CH.UpdateChannelDetails(&chn, channelid, channel.DB, tenantid)

	return nil
}

// Get channel count
func (channel *Channel) GetChannelCount(tenantid string) (count int, err error) {

	var chcount int64

	err = CH.GetChannelCount(&chcount, channel.DB, tenantid)

	if err != nil {

		return 0, err
	}

	return int(chcount), nil

}

// Channel type change
func (channel *Channel) ChannelType(Channels Tblchannel) error {

	autherr := AuthandPermission(channel)

	if autherr != nil {

		return autherr
	}

	var channeltype Tblchannel

	channeltype.Id = Channels.Id

	channeltype.CollectionCount = Channels.CollectionCount

	channeltype.CloneCount = Channels.CloneCount

	err := CH.ChangeChanelType(channeltype, channel.DB)

	if err != nil {

		return err
	}

	return nil

}

func (channel *Channel) CheckNameInChannel(channelid int, cname string, tenantid string) (bool, error) {

	channeldet, err := CH.CheckNameInChannel(channelid, cname, channel.DB, tenantid)

	if err != nil {
		return false, err
	}
	if channeldet.Id == 0 {

		return false, err
	}

	return true, nil

}
func (channel *Channel) GetChannal(chname string, tenantid string) int {

	channelid, _ := CH.GetChannelId(chname, tenantid, channel.DB)
	return channelid
}

func (channel *Channel) CheckNameInFolder(channelid int, folderid int, slugname string, tenantid string) (bool, error) {

	channeldet, err := CH.CheckNameInFolder(channelid, folderid, slugname, channel.DB, tenantid)

	if err != nil {
		return false, err
	}
	if channeldet.Id == 0 {

		return false, err
	}

	return true, nil

}

func (channel *Channel) GetFolderFilelByChannelid(channelid int, tenantid string) ([]TblFolder, error) {

	FolderndFiledata, err := CH.GetFolderFilelByChannelid(channelid, tenantid, channel.DB)

	if err != nil {

		return []TblFolder{}, nil
	}

	return FolderndFiledata, nil
}
