package models

import (
	// "time"
	//	"github.com/huawei-openlab/newdb/orm"
	//	"github.com/ngaut/log"
)

type Project struct {
	Id          int64     `json:"id"`
	WorkspaceId int64     `json:"workspaceId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Actived     bool      `json:"actived"`
	Created     int64 `json:"created"`
	Updated     int64 `json:"updated"`
	Memo        string    `json:"memo"`
}

func (pj *Project) Create(wid int64, name, description string) (int64, error) {
	/*
	           o := orm.NewOrm()
	   	p := Project{WorkspaceId: wid, Name: name, Description: description, Actived: true}

	   	if err := o.Begin(); err != nil {
	   		log.Errorf("Transcation error: %s", err.Error())

	   		return 0, err
	   	} else {
	   		if id, e := o.Insert(&p); e != nil {
	   			log.Errorf("Create project error: %s", e.Error())

	   			o.Rollback()
	   			return 0, err
	   		} else {
	   			log.Errorf("Create project successfully, id is: %d", id)

	   			o.Commit()
	   			return id, nil
	   		}
	   	}
	*/
	return 0, nil
}

func (pj *Project) Put(id int64, name, description string) error {
	/*	o := orm.NewOrm()
		p := Project{Id: id, Actived: true}

		if err := o.Read(&p, "Id", "Actived"); err != nil {
			log.Errorf("Get project %d error: %s", id, err.Error())

			return err
		} else {
			if err := o.Begin(); err != nil {
				log.Errorf("Transcation error: %s", err.Error())

				o.Rollback()
				return err
			} else {
				p.Name = name
				p.Description = description

				if _, err := o.Update(&p, "Name", "Description"); err != nil {
					log.Errorf("Put project %d error: %s", id, err.Error())

					o.Rollback()
					return err
				} else {
					log.Infof("Put project %d successfully")

					o.Commit()
					return nil
				}
			}
		}
	*/
	return nil
}

func (pj *Project) Get(id int64) (*Project, error) {
	/*
		o := orm.NewOrm()
		p := Project{Id: id, Actived: true}

		if err := o.Read(&p, "Id", "Actived"); err != nil {
			log.Errorf("Get project %d error: %s", id, err.Error())

			return p, err
		} else {
			return p, nil
		}
	*/
	return nil, nil
}

func (pj *Project) Delete(id int64) error {
	/*
		o := orm.NewOrm()
		p := Project{Id: id}

		if err := o.Read(&p, "Id"); err != nil {
			log.Errorf("Get project %d error: %s", id, err.Error())

			o.Rollback()
			return err
		} else {
			if err := o.Begin(); err != nil {
				log.Errorf("Transcation error: %s", err.Error())

				o.Rollback()
				return err
			} else {
				if _, err := o.Delete(&p); err != nil {
					log.Errorf("Delete project %d error: %s", id, err.Error())

					o.Rollback()
					return err
				} else {
					log.Infof("Delete project %d successfully", id)

					//TODO Delete relate objetcts.

					o.Commit()
					return nil
				}
			}
		}
	*/
	return nil

}
