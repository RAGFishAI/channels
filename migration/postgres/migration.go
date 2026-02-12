package postgres

import (
	"time"

	"gorm.io/gorm"
)

type TblChannel struct {
	Id                 int       `gorm:"primaryKey;auto_increment;type:serial"`
	ChannelName        string    `gorm:"type:character varying"`
	ChannelDescription string    `gorm:"type:character varying"`
	SlugName           string    `gorm:"type:character varying"`
	FieldGroupId       int       `gorm:"type:integer"`
	IsActive           int       `gorm:"type:integer"`
	IsDeleted          int       `gorm:"type:integer"`
	CreatedOn          time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy          int       `gorm:"type:integer"`
	ModifiedOn         time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"DEFAULT:NULL"`
	TenantId           int       `gorm:"type:integer"`
}

// MigrateTable creates this package related tables in your database
func MigrationTables(db *gorm.DB) {

	if err := db.AutoMigrate(

		&TblChannel{},
	); err != nil {

		panic(err)
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS email_unique
    ON public.tbl_members USING btree
    (email COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE is_deleted = 0;`)

	db.Exec(`CREATE INDEX IF NOT EXISTS mobile_no_unique
    ON public.tbl_members USING btree
    (mobile_no COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE is_deleted = 0;`)

	//create default channel
	db.Exec(`INSERT INTO public.tbl_channels(id, channel_name, slug_name, field_group_id, is_active, is_deleted, created_on, created_by, channel_description) VALUES (1, 'Default_Channel', 'default_channel', 0, 1, 0, '2024-03-04 10:49:17', '1', 'default description');`)

	db.Exec(`INSERT INTO public.tbl_channel_categories(id, channel_id, category_id, created_at, created_on) VALUES (1, 1, '1,2', 1, '2024-03-04 10:49:17');`)

	//Channel default fields
	db.Exec(`INSERT INTO public.tbl_field_types(id, type_name, type_slug, is_active, is_deleted, created_by, created_on) VALUES (1, 'Label', 'label', 1,  0, 1, '2023-03-14 11:09:12'), (2, 'Text', 'text', 1,  0, 1, '2023-03-14 11:09:12'),(3, 'Link', 'link', 1,  0, 1, '2023-03-14 11:09:12'),(4, 'Date & Time', 'date&time', 1,  0, 1, '2023-03-14 11:09:12'), (5, 'Select', 'select', 1,  0, 1, '2023-03-14 11:09:12'),(6, 'Date', 'date', 1,  0, 1, '2023-03-14 11:09:12'),(7, 'TextBox', 'textbox', 1,  0, 1, '2023-03-14 11:09:12'),(8, 'TextArea', 'textarea', 1, 0, 1, '2023-03-14 11:09:12'), (9, 'Radio Button', 'radiobutton', 1, 0, 1, '2023-03-14 11:09:12'),(10, 'CheckBox', 'checkbox', 1, 0, 1, '2023-03-14 11:09:12'),(11, 'Text Editor', 'texteditor', 1, 0, 1, '2023-03-14 11:09:12'),(12, 'Section', 'section', 1, 0, 1, '2023-03-14 11:09:12'),(13, 'Section Break', 'sectionbreak', 1, 0, 1, '2023-03-14 11:09:12'),(14, 'Members', 'member', 1,  0, 1, '2023-03-14 11:09:12');
	`)
}
