// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"github.com/gocql/gocql"
)

type caManager struct {
	clients []*gocql.Session
}

var camanager = newCaManager()

func newCaManager() *caManager {
	return &caManager{}
}

func (cm *caManager) init() error {
	return nil
}

type cdTimMessage struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	Mid        int64      `cql:"mid"`
	ChatID     []byte     `cql:"chatid"`
	FID        int32      `cql:"fid"`
	Stanza     []byte     `cql:"stanza"`
	Timeseries int64      `cql:"timeseries"`
}

func newCdTimMessage() *cdTimMessage {
	return &cdTimMessage{}
}

func (t *cdTimMessage) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timmessage (id, mid, chatid, fid, stanza, timeseries) VALUES (?, ?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Mid, t.ChatID, t.FID, t.Stanza, t.Timeseries).Exec()
}

func (t *cdTimMessage) GetById(id gocql.UUID) (*cdTimMessage, error) {
	var msg cdTimMessage
	query := `SELECT id, mid, chatid, fid, stanza, timeseries FROM timmessage WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&msg.ID, &msg.Mid, &msg.ChatID, &msg.FID, &msg.Stanza, &msg.Timeseries)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (t *cdTimMessage) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timmessage SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimMessage) Delete(id gocql.UUID) error {
	query := `DELETE FROM timmessage WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimMessage) List() ([]*cdTimMessage, error) {
	var messages []*cdTimMessage
	query := `SELECT id, mid, chatid, fid, stanza, timeseries FROM timmessage`
	iter := t.session.Query(query).Iter()
	for {
		var msg cdTimMessage
		if !iter.Scan(&msg.ID, &msg.Mid, &msg.ChatID, &msg.FID, &msg.Stanza, &msg.Timeseries) {
			break
		}
		messages = append(messages, &msg)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

type cdTimUser struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	UUID       int64      `cql:"uuid"`
	Pwd        string     `cql:"pwd"`
	Createtime int64      `cql:"createtime"`
	Ubean      []byte     `cql:"ubean"`
	Timeseries int64      `cql:"timeseries"`
}

func newCdTimUser() *cdTimUser {
	return &cdTimUser{}
}

func (t *cdTimUser) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timuser (id, uuid, pwd, createtime, ubean, timeseries) VALUES (?, ?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.UUID, t.Pwd, t.Createtime, t.Ubean, t.Timeseries).Exec()
}

func (t *cdTimUser) GetById(id gocql.UUID) (*cdTimUser, error) {
	var user cdTimUser
	query := `SELECT id, uuid, pwd, createtime, ubean, timeseries FROM timuser WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&user.ID, &user.UUID, &user.Pwd, &user.Createtime, &user.Ubean, &user.Timeseries)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (t *cdTimUser) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timuser SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimUser) Delete(id gocql.UUID) error {
	query := `DELETE FROM timuser WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimUser) List() ([]*cdTimUser, error) {
	var users []*cdTimUser
	query := `SELECT id, uuid, pwd, createtime, ubean, timeseries FROM timuser`
	iter := t.session.Query(query).Iter()
	for {
		var user cdTimUser
		if !iter.Scan(&user.ID, &user.UUID, &user.Pwd, &user.Createtime, &user.Ubean, &user.Timeseries) {
			break
		}
		users = append(users, &user)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return users, nil
}

type cdTimOffline struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	UUID       int64      `cql:"uuid"`
	Chatid     []byte     `cql:"chatid"`
	Stanza     []byte     `cql:"stanza"`
	Mid        int64      `cql:"mid"`
	Timeseries int64      `cql:"timeseries"`
}

func newTimOffline() *cdTimOffline {
	return &cdTimOffline{}
}

func (t *cdTimOffline) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timoffline (id, uuid, chatid, stanza, mid, timeseries) VALUES (?, ?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.UUID, t.Chatid, t.Stanza, t.Mid, t.Timeseries).Exec()
}

func (t *cdTimOffline) GetById(id gocql.UUID) (*cdTimOffline, error) {
	var offline cdTimOffline
	query := `SELECT id, uuid, chatid, stanza, mid, timeseries FROM timoffline WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&offline.ID, &offline.UUID, &offline.Chatid, &offline.Stanza, &offline.Mid, &offline.Timeseries)
	if err != nil {
		return nil, err
	}
	return &offline, nil
}

func (t *cdTimOffline) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timoffline SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimOffline) Delete(id gocql.UUID) error {
	query := `DELETE FROM timoffline WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimOffline) List() ([]*cdTimOffline, error) {
	var offlines []*cdTimOffline
	query := `SELECT id, uuid, chatid, stanza, mid, timeseries FROM timoffline`
	iter := t.session.Query(query).Iter()
	for {
		var offline cdTimOffline
		if !iter.Scan(&offline.ID, &offline.UUID, &offline.Chatid, &offline.Stanza, &offline.Mid, &offline.Timeseries) {
			break
		}
		offlines = append(offlines, &offline)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return offlines, nil
}

type cdTimRelate struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	UUID       []byte     `cql:"uuid"`
	Status     int32      `cql:"status"`
	Timeseries int64      `cql:"timeseries"`
}

func newCdTimRelate() *cdTimRelate {
	return &cdTimRelate{}
}

func (t *cdTimRelate) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timrelate (id, uuid, status, timeseries) VALUES (?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.UUID, t.Status, t.Timeseries).Exec()
}

func (t *cdTimRelate) GetById(id gocql.UUID) (*cdTimRelate, error) {
	var relate cdTimRelate
	query := `SELECT id, uuid, status, timeseries FROM timrelate WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&relate.ID, &relate.UUID, &relate.Status, &relate.Timeseries)
	if err != nil {
		return nil, err
	}
	return &relate, nil
}

func (t *cdTimRelate) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timrelate SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimRelate) Delete(id gocql.UUID) error {
	query := `DELETE FROM timrelate WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimRelate) List() ([]*cdTimRelate, error) {
	var relates []*cdTimRelate
	query := `SELECT id, uuid, status, timeseries FROM timrelate`
	iter := t.session.Query(query).Iter()
	for {
		var relate cdTimRelate
		if !iter.Scan(&relate.ID, &relate.UUID, &relate.Status, &relate.Timeseries) {
			break
		}
		relates = append(relates, &relate)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return relates, nil
}

type cdTimRoster struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	Unikid     []byte     `cql:"unikid"`
	UUID       int64      `cql:"uuid"`
	Tuuid      int64      `cql:"tuuid"`
	Timeseries int64      `cql:"timeseries"`
}

func newTimRoster() *cdTimRoster {
	return &cdTimRoster{}
}

func (t *cdTimRoster) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timroster (id, unikid, uuid, tuuid, timeseries) VALUES (?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Unikid, t.UUID, t.Tuuid, t.Timeseries).Exec()
}

func (t *cdTimRoster) GetById(id gocql.UUID) (*cdTimRoster, error) {
	var roster cdTimRoster
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timroster WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&roster.ID, &roster.Unikid, &roster.UUID, &roster.Tuuid, &roster.Timeseries)
	if err != nil {
		return nil, err
	}
	return &roster, nil
}

func (t *cdTimRoster) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timroster SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimRoster) Delete(id gocql.UUID) error {
	query := `DELETE FROM timroster WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimRoster) List() ([]*cdTimRoster, error) {
	var rosters []*cdTimRoster
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timroster`
	iter := t.session.Query(query).Iter()
	for {
		var roster cdTimRoster
		if !iter.Scan(&roster.ID, &roster.Unikid, &roster.UUID, &roster.Tuuid, &roster.Timeseries) {
			break
		}
		rosters = append(rosters, &roster)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return rosters, nil
}

type cdTimMucroster struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	Unikid     []byte     `cql:"unikid"`
	UUID       int64      `cql:"uuid"`
	Tuuid      int64      `cql:"tuuid"`
	Timeseries int64      `cql:"timeseries"`
}

func newTimMucroster() *cdTimMucroster {
	return &cdTimMucroster{}
}

func (t *cdTimMucroster) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timmucroster (id, unikid, uuid, tuuid, timeseries) VALUES (?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Unikid, t.UUID, t.Tuuid, t.Timeseries).Exec()
}

func (t *cdTimMucroster) GetById(id gocql.UUID) (*cdTimMucroster, error) {
	var mucroster cdTimMucroster
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timmucroster WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&mucroster.ID, &mucroster.Unikid, &mucroster.UUID, &mucroster.Tuuid, &mucroster.Timeseries)
	if err != nil {
		return nil, err
	}
	return &mucroster, nil
}

func (t *cdTimMucroster) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timmucroster SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimMucroster) Delete(id gocql.UUID) error {
	query := `DELETE FROM timmucroster WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimMucroster) List() ([]*cdTimMucroster, error) {
	var mucrosters []*cdTimMucroster
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timmucroster`
	iter := t.session.Query(query).Iter()
	for {
		var mucroster cdTimMucroster
		if !iter.Scan(&mucroster.ID, &mucroster.Unikid, &mucroster.UUID, &mucroster.Tuuid, &mucroster.Timeseries) {
			break
		}
		mucrosters = append(mucrosters, &mucroster)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return mucrosters, nil
}

type cdTimBlock struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	Unikid     []byte     `cql:"unikid"`
	UUID       int64      `cql:"uuid"`
	Tuuid      int64      `cql:"tuuid"`
	Timeseries int64      `cql:"timeseries"`
}

func newTimBlock() *cdTimBlock {
	return &cdTimBlock{}
}

func (t *cdTimBlock) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timblock (id, unikid, uuid, tuuid, timeseries) VALUES (?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Unikid, t.UUID, t.Tuuid, t.Timeseries).Exec()
}

func (t *cdTimBlock) GetById(id gocql.UUID) (*cdTimBlock, error) {
	var block cdTimBlock
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timblock WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&block.ID, &block.Unikid, &block.UUID, &block.Tuuid, &block.Timeseries)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (t *cdTimBlock) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timblock SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimBlock) Delete(id gocql.UUID) error {
	query := `DELETE FROM timblock WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimBlock) List() ([]*cdTimBlock, error) {
	var blocks []*cdTimBlock
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timblock`
	iter := t.session.Query(query).Iter()
	for {
		var block cdTimBlock
		if !iter.Scan(&block.ID, &block.Unikid, &block.UUID, &block.Tuuid, &block.Timeseries) {
			break
		}
		blocks = append(blocks, &block)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return blocks, nil
}

type cdTimBlockroom struct {
	session    *gocql.Session
	ID         gocql.UUID `cql:"id"`
	Unikid     []byte     `cql:"unikid"`
	UUID       int64      `cql:"uuid"`
	Tuuid      int64      `cql:"tuuid"`
	Timeseries int64      `cql:"timeseries"`
}

func newCdTimBlockroom() *cdTimBlockroom {
	return &cdTimBlockroom{}
}

func (t *cdTimBlockroom) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timblockroom (id, unikid, uuid, tuuid, timeseries) VALUES (?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Unikid, t.UUID, t.Tuuid, t.Timeseries).Exec()
}

func (t *cdTimBlockroom) GetById(id gocql.UUID) (*cdTimBlockroom, error) {
	var blockroom cdTimBlockroom
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timblockroom WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&blockroom.ID, &blockroom.Unikid, &blockroom.UUID, &blockroom.Tuuid, &blockroom.Timeseries)
	if err != nil {
		return nil, err
	}
	return &blockroom, nil
}

func (t *cdTimBlockroom) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timblockroom SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimBlockroom) Delete(id gocql.UUID) error {
	query := `DELETE FROM timblockroom WHERE id = ?`
	return t.session.Query(query, id).Exec()
}

func (t *cdTimBlockroom) List() ([]*cdTimBlockroom, error) {
	var blockrooms []*cdTimBlockroom
	query := `SELECT id, unikid, uuid, tuuid, timeseries FROM timblockroom`
	iter := t.session.Query(query).Iter()
	for {
		var blockroom cdTimBlockroom
		if !iter.Scan(&blockroom.ID, &blockroom.Unikid, &blockroom.UUID, &blockroom.Tuuid, &blockroom.Timeseries) {
			break
		}
		blockrooms = append(blockrooms, &blockroom)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return blockrooms, nil
}

type cdTimDomain struct {
	session       *gocql.Session
	ID            gocql.UUID `cql:"id"`
	Adminaccount  string     `cql:"adminaccount"`
	Adminpassword string     `cql:"adminpassword"`
	Timdomain     string     `cql:"timdomain"`
	Createtime    int64      `cql:"createtime"`
	Timeseries    int64      `cql:"timeseries"`
}

func newCdTimDomain() *cdTimDomain {
	return &cdTimDomain{}
}

func (t *cdTimDomain) Create() error {
	t.ID = gocql.TimeUUID()
	query := `INSERT INTO timdomain (id, adminaccount, adminpassword, timdomain, createtime, timeseries) VALUES (?, ?, ?, ?, ?, ?)`
	return t.session.Query(query, t.ID, t.Adminaccount, t.Adminpassword, t.Timdomain, t.Createtime, t.Timeseries).Exec()
}

func (t *cdTimDomain) GetById(id gocql.UUID) (*cdTimDomain, error) {
	var domain cdTimDomain
	query := `SELECT id, adminaccount, adminpassword, timdomain, createtime, timeseries FROM timdomain WHERE id = ? LIMIT 1`
	err := t.session.Query(query, id).Consistency(gocql.One).Scan(&domain.ID, &domain.Adminaccount, &domain.Adminpassword, &domain.Timdomain, &domain.Createtime, &domain.Timeseries)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (t *cdTimDomain) Update(id gocql.UUID, updates map[string]interface{}) error {
	query := `UPDATE timdomain SET `
	params := []interface{}{}
	i := 0
	for key, value := range updates {
		if i > 0 {
			query += ", "
		}
		query += key + " = ?"
		params = append(params, value)
		i++
	}
	query += " WHERE id = ?"
	params = append(params, id)
	return t.session.Query(query, params...).Exec()
}

func (t *cdTimDomain) Delete(id gocql.UUID) error {
	query := `DELETE FROM timdomain WHERE id = ?`
	return t.session.Query(query, id).Exec()
}
