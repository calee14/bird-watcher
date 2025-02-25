package db

import (
	"database/sql"
	"errors"
)

func CreateSubscriber(subscriber *Subscriber) error {
	query := `
	insert into subscribers (email, created_at)
	values (?, ?)
	`

	result, err := DB.Exec(query, subscriber.Email, subscriber.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	subscriber.ID = id
	return nil
}

func GetSubscriberByID(id int64) (*Subscriber, error) {
	query := `select id, email, created_at from subscribers where id = ?`
	var subscriber Subscriber
	err := DB.QueryRow(query, id).Scan(&subscriber.ID, &subscriber.Email, &subscriber.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &subscriber, err
}

func GetAllSubscribers() ([]*Subscriber, error) {
	query := `select id, email, created_at from subscribers order by created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []*Subscriber
	for rows.Next() {
		var subscriber Subscriber
		if err := rows.Scan(&subscriber.ID, &subscriber.Email, &subscriber.CreatedAt); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, &subscriber)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}

func UpdateSubscriber(subscriber *Subscriber) error {
	query := `
	update todos
	set email = ?
	where id = ?
	`

	_, err := DB.Exec(query, subscriber.Email, subscriber.ID)
	return err
}

func DeleteSubscriber(email string) error {
	query := `delete from subscribers where email = ?`

	_, err := DB.Exec(query, email)
	return err
}
