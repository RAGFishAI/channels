package channels

import (
	"fmt"
	"time"

	"github.com/RAGFishAI/team"
	"gorm.io/gorm"
)

type Filter struct {
	Keyword    string
	Status     bool
	CreateOnly bool
}

type Channels struct {
	Id           int
	Slug         string
	Limit        int
	Offset       int
	Keyword      string
	IsActive     bool
	SortBy       string
	SortingOrder int
	TenantId     string
	EntriesCount bool
	AuthorDetail bool
	CreateOnly   bool
	Count        bool
	ChannelFile  bool
}
type TblFiles struct {
	Id             int       `gorm:"primaryKey;auto_increment;type:serial"`
	FileName       string    `gorm:"type:character varying"`
	UniqueFileName string    `gorm:"type:character varying"`
	FilePath       string    `gorm:"type:character varying"`
	FileId         string    `gorm:"type:character varying"`
	FolderName     string    `gorm:"type:character varying"`
	ChannelId      int       `gorm:"type:integer"`
	CreatedBy      int       `gorm:"type:integer"`
	CreatedOn      time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	IsActive       int       `gorm:"type:integer"`
	IsDeleted      int       `gorm:"type:integer;DEFAULT:0"`
	DeletedOn      time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy      int       `gorm:"type:integer;DEFAULT:NULL"`
	DateString     string    `gorm:"-"`
	TenantId       string    `gorm:"type:character varying"`
	FolderId       int       `gorm:"type:integer"`
}

type Tblchannel struct {
	Id                 int       `gorm:"column:id"`
	ChannelName        string    `gorm:"column:channel_name"`
	ChannelUniqueId    string    `gorm:"column:channel_unique_id"`
	ChannelDescription string    `gorm:"column:channel_description"`
	SlugName           string    `gorm:"column:slug_name"`
	FieldGroupId       int       `gorm:"column:field_group_id"`
	IsActive           int       `gorm:"column:is_active"`
	IsDeleted          int       `gorm:"column:is_deleted"`
	CreatedOn          time.Time `gorm:"column:created_on"`
	CreatedBy          int       `gorm:"column:created_by"`
	ModifiedOn         time.Time `gorm:"column:modified_on;DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"column:modified_by;DEFAULT:NULL"`
	DateString         string    `gorm:"-"`
	EntriesCount       int       `gorm:"-"`

	ProfileImagePath string       `gorm:"<-:false"`
	AuthorDetails    team.TblUser `gorm:"foreignKey:Id;references:CreatedBy"`
	ChannelType      string       `gorm:"column:channel_type"`
	CollectionCount  int          `gorm:"column:collection_count"`
	CloneCount       int          `gorm:"column:clone_count"`
	TenantId         string       `gorm:"column:tenant_id"`
	Username         string       `gorm:"<-:false"`
	FirstName        string       `gorm:"<-:false"`
	LastName         string       `gorm:"<-:false"`
	NameString       string       `gorm:"<-:false"`
	ImagePath        string       `gorm:"column:image_path"`
	SeoTitle         string       `gorm:"column:seo_title"`
	SeoDescription   string       `gorm:"column:seo_description"`
	SeoKeyword       string       `gorm:"column:seo_keyword"`
	FileCount        int64        `gorm:"<-:false"`
	FirstFolderId    int          `gorm:"<-:false"`
	FilesData        []TblFiles   `gorm:"foreignKey:ChannelId;references:Id"`
	FolderCount      int          `gorm:"<-:false"`
}

type TblChannel struct {
	Id                 int       `gorm:"primaryKey;auto_increment;type:serial"`
	ChannelName        string    `gorm:"type:character varying"`
	ChannelUniqueId    string    `gorm:"type:character varying"`
	ChannelDescription string    `gorm:"type:character varying"`
	SlugName           string    `gorm:"type:character varying"`
	IsActive           int       `gorm:"type:integer"`
	IsDeleted          int       `gorm:"type:integer"`
	CreatedOn          time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy          int       `gorm:"type:integer"`
	ModifiedOn         time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"DEFAULT:NULL"`
	TenantId           string    `gorm:"type:character varying"`
	ChannelType        string    `gorm:"type:character varying"`
}

type ChannelCreate struct {
	ChannelName        string
	ChannelUniqueId    string
	ChannelDescription string
	CategoryIds        []string
	CreatedBy          int
	CollectionCount    int
	ImagePath          string
	SeoTitle           string
	SeoDescription     string
	SeoKeyword         string
	SlugName           string
}

type ChannelModel struct {
	Userid     int
	Dataaccess int
}

var CH ChannelModel

// soft delete check
func IsDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = 0")
}

/*channel list*/
func (Ch ChannelModel) Channellist(DB *gorm.DB, channel *Channel, inputs Channels, channels *[]Tblchannel, count *int64) (err error) {

	query := DB.Table("tbl_channels").
		Select(`
        tbl_channels.*, 
        COUNT(tbl_files.id) as file_count,
        (SELECT COUNT(*) FROM tbl_folders 
         WHERE tbl_folders.channel_id = tbl_channels.id 
         AND tbl_folders.is_deleted = 0) as folder_count,
        (SELECT id FROM tbl_folders 
         WHERE tbl_folders.channel_id = tbl_channels.id 
         AND tbl_folders.is_deleted = 0 
         ORDER BY tbl_folders.id 
         LIMIT 1) as first_folder_id
    `).
		Where("tbl_channels.is_deleted = 0").
		Joins("LEFT JOIN tbl_files ON tbl_channels.id = tbl_files.channel_id AND tbl_files.is_deleted = 0").
		Group("tbl_channels.id")

	if inputs.TenantId != "" {

		query = query.Where("tbl_channels.tenant_id=?", inputs.TenantId)

	} else {

		query = query.Where("tbl_channels.tenant_id is null")
	}

	if inputs.CreateOnly && Ch.Dataaccess == 1 {

		query = query.Where("tbl_channels.created_by = ?", Ch.Userid)
	}

	if inputs.Keyword != "" {

		query = query.Where("LOWER(TRIM(channel_name)) LIKE LOWER(TRIM(?))", "%"+inputs.Keyword+"%")
	}

	if inputs.IsActive {

		query = query.Where("tbl_channels.is_active=1")

	}

	if inputs.Count {

		err = query.Count(count).Error

		if err != nil {

			return err
		}
	}

	if inputs.AuthorDetail {

		query = query.Preload("AuthorDetails", "is_deleted = ?", 0)
	}

	if inputs.SortBy != "" {

		if inputs.SortingOrder == 0 {

			query = query.Order(inputs.SortingOrder)

		} else if inputs.SortingOrder == 1 {

			query = query.Order(inputs.SortBy + " desc")

		}

	} else {

		query = query.Order("id desc")
	}

	if inputs.Limit != 0 {

		query = query.Limit(inputs.Limit)
	}

	if inputs.Offset != -1 {

		query = query.Offset(inputs.Offset)
	}

	err = query.Find(&channels).Error

	if err != nil {

		return err
	}

	return nil
}

/*Craete channel */
func (Ch ChannelModel) CreateChannel(chn *TblChannel, DB *gorm.DB) (TblChannel, error) {

	if err := DB.Debug().Table("tbl_channels").Create(&chn).Error; err != nil {

		return TblChannel{}, err

	}

	return *chn, nil

}

func (Ch ChannelModel) ChannelDetail(DB *gorm.DB, inputs Channels, channelDetail *Tblchannel) error {

	query := DB.Table("tbl_channels").Where("tbl_channels.is_deleted = 0")

	if inputs.Id != 0 {

		query = query.Where("id=?", inputs.Id)
	}

	if inputs.Slug != "" {

		query = query.Where("slug_name=?", inputs.Slug)
	}

	if inputs.TenantId != "" {

		query = query.Where("tbl_channels.tenant_id=?", inputs.TenantId)
	}

	if inputs.CreateOnly && Ch.Dataaccess == 1 {

		query = query.Where("tbl_channels.created_by = ?", Ch.Userid)
	}

	if inputs.Keyword != "" {

		query = query.Where("LOWER(TRIM(channel_name)) LIKE LOWER(TRIM(?))", "%"+inputs.Keyword+"%")
	}

	if inputs.IsActive {

		query = query.Where("tbl_channels.is_active=?", 1)
	}

	if inputs.AuthorDetail {

		query = query.Preload("AuthorDetails", "is_deleted = ?", 0)
	}

	err := query.First(&channelDetail).Error

	if err != nil {

		return err
	}

	return nil
}

func (Ch ChannelModel) GetChannelByChannelName(name string, DB *gorm.DB, tenantid string) (ch Tblchannel, err error) {

	if err := DB.Table("tbl_channels").Where("LOWER(TRIM(channel_name)) = LOWER(TRIM(?))   and tenant_id=? and is_deleted=0", name, tenantid).First(&ch).Error; err != nil {

		return Tblchannel{}, err
	}

	return ch, nil
}

/*Get Channel*/
func (Ch ChannelModel) GetChannelById(id int, DB *gorm.DB, tenantid string) (ch Tblchannel, err error) {

	if tenantid != "" {

		if err := DB.Table("tbl_channels").Where("id=? and tenant_id=?", id, tenantid).First(&ch).Error; err != nil {

			return Tblchannel{}, err
		}
	} else {

		if err := DB.Table("tbl_channels").Where("id=? and tenant_id is null", id).First(&ch).Error; err != nil {

			return Tblchannel{}, err
		}
	}
	return ch, nil
}

/*Delete Channel*/
func (Ch ChannelModel) DeleteChannelById(id int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_channels").Where("id=? and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil
}

/*Isactive channel*/
func (Ch ChannelModel) ChannelIsActive(tblch *TblChannel, id, val int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_channels").Where("id=? and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_active": val, "modified_on": tblch.ModifiedOn, "modified_by": tblch.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*Update Channel Details*/
func (Ch ChannelModel) UpdateChannelDetails(chn *TblChannel, id int, DB *gorm.DB, TenantId string) error {

	if err := DB.Table("tbl_channels").Where("id=? and tenant_id=?", id, TenantId).UpdateColumns(map[string]interface{}{"channel_name": chn.ChannelName, "slug_name": chn.SlugName, "channel_unique_id": chn.ChannelUniqueId, "channel_description": chn.ChannelDescription, "modified_by": chn.ModifiedBy, "modified_on": chn.ModifiedOn}).Error; err != nil {

		return err
	}

	fmt.Println("UpdateChannelDetails:", chn)

	return nil
}

func (ch ChannelModel) GetChannelCount(count *int64, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
		// Joins("inner join tbl_channel_categories on tbl_channel_categories.channel_id = tbl_channels.id").
		Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.status = 1 and tbl_channel_entries.tenant_id=?", tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (ch ChannelModel) GetChannels(channels *[]Tblchannel, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_channels").Where("is_deleted = 0 and is_active = 1 and tenant_id=?", tenantid).Find(&channels).Error; err != nil {

		return err
	}

	return nil
}

// Channel type change
func (ch ChannelModel) ChangeChanelType(Channels Tblchannel, DB *gorm.DB) (Error error) {

	if Channels.CollectionCount != 0 {
		if err := DB.Table("tbl_channels").Where("id=?", Channels.Id).Updates(map[string]interface{}{"collection_count": Channels.CollectionCount}).Error; err != nil {

			return err

		}
	} else {
		if err := DB.Table("tbl_channels").Where("id=?", Channels.Id).Updates(map[string]interface{}{"clone_count": Channels.CloneCount}).Error; err != nil {

			return err

		}

	}

	return nil
}

//Total Channel Count

func (ch ChannelModel) AllChannelCount(DB *gorm.DB, tenantid string) (count int64, err error) {

	if err := DB.Table("tbl_channels").Where("tbl_channels.is_deleted = 0 and  tenant_id = ?", tenantid).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}

// Last 10 days Channel Count
func (ch ChannelModel) NewChannelCount(DB *gorm.DB, tenantid string) (count int64, err error) {

	if err := DB.Table("tbl_channels").Where("created_on >=? and  tenant_id=? and is_deleted = 0", time.Now().AddDate(0, 0, -10), tenantid).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}

func (ch ChannelModel) CheckNameInChannel(channelid int, channelname string, DB *gorm.DB, tenantid string) (channel TblChannel, err error) {

	if channelid == 0 {

		if err := DB.Table("tbl_channels").Where("LOWER(TRIM(channel_name))=LOWER(TRIM(?)) and tenant_id=? and is_deleted=0", channelname, tenantid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}
	} else {

		if err := DB.Table("tbl_channels").Where("LOWER(TRIM(channel_name))=LOWER(TRIM(?)) and id not in (?) and tenant_id=?   and is_deleted=0", channelname, channelid, tenantid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}
	}

	return channel, nil

}

func (ch ChannelModel) GetChannelId(chname string, tenantid string, DB *gorm.DB) (int, error) {

	var Id int // Define the variable to hold the result

	if err := DB.Table("tbl_channels").Where("channel_name = ? and tenant_id=? and is_deleted=0", chname, tenantid).Select("Id").Scan(&Id).Error; err != nil {
		return 0, err
	}

	return Id, nil

}

func (ch ChannelModel) CheckNameInFolder(channelid, folderid int, foldername string, DB *gorm.DB, tenantid string) (channel TblChannel, err error) {

	if folderid == 0 {

		if err := DB.Table("tbl_folders").Where("LOWER(TRIM(folder_name))=LOWER(TRIM(?)) and channel_id=? and tenant_id=? and is_deleted=0", foldername, channelid, tenantid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}
	} else {

		if err := DB.Table("tbl_folders").Where("LOWER(TRIM(folder_name))=LOWER(TRIM(?)) and id not in (?) and channel_id=? and tenant_id=?   and is_deleted=0", foldername, folderid, channelid, tenantid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}
	}

	return channel, nil

}

func (ch ChannelModel) GetFilesByChannelId(channelid int, DB *gorm.DB, tenantid string) (files []TblFiles, err error) {
	if err := DB.Table("tbl_files").Where("channel_id=? and is_deleted=0 and tenant_id=?", channelid, tenantid).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}
