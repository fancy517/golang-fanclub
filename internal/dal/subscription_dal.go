package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fmt"
	"strconv"
)

type SubscriptionDAL interface {
	GetBasePrice(tierid string) (float64, error)
	CreateSubscriber(userid string, tierid string, duration string) error
	GetTierCards(userid string, filter string) ([]models.TierCards, error)
	GetTierCounts(userid string) (models.TiersCount, error)
	GetTierOwner(tierid string) (int, error)
}

type subscriptionDAL struct {
	DB *sql.DB
}

var _ SubscriptionDAL = (*subscriptionDAL)(nil)

func NewSubscriptionDAL(db *sql.DB) SubscriptionDAL {
	return &subscriptionDAL{db}
}

func (dal *subscriptionDAL) GetTierOwner(tierid string) (int, error) {
	query := `SELECT users.id FROM tbl_subscriptions LEFT JOIN users ON users.name = tbl_subscriptions.username WHERE tbl_subscriptions.id = ?`
	ownerid := 0
	err := dal.DB.QueryRow(query, tierid).Scan(&ownerid)
	if err != nil {
		fmt.Printf("Error Scanning %v\n", err)
		return 0, err
	}
	return ownerid, nil
}

func (dal *subscriptionDAL) GetTierCounts(userid string) (models.TiersCount, error) {
	res := models.TiersCount{
		Active:  "0",
		Expired: "0",
	}
	query := `SELECT COUNT(*) AS cnt FROM (SELECT CASE WHEN end_date < NOW() THEN "expired" ELSE "active" END AS STATUS FROM tbl_subscribers WHERE tbl_subscribers.userid = ?) AS A GROUP BY A.status ORDER BY A.status`
	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("Error querying %v\n", err)
		return res, err
	}
	defer rows.Close()

	idx := 0
	for rows.Next() {
		count := 0
		err := rows.Scan(&count)
		if err != nil {
			fmt.Printf("Error scanning %v\n", err)
			return res, err
		}
		if idx == 0 {
			res.Active = strconv.Itoa(count)
		} else {
			res.Expired = strconv.Itoa(count)
		}
		idx++
	}

	return res, nil
}

func (dal *subscriptionDAL) GetTierCards(userid string, filter string) ([]models.TierCards, error) {
	query := `SELECT
		users.name AS username,
		tbl_userdata.displayname,
		users.avatar,
		tbl_userdata.userBanner AS banner,
		tbl_subscriptions.tier_name AS tiername,
		tbl_subscribers.end_date AS expiration,
		tbl_subscriptions.base_price AS baseprice,
		CASE
		WHEN tbl_subscribers.end_date < NOW()
		THEN "expired"
		ELSE "active"
		END AS astatus
	FROM
		tbl_subscribers
		LEFT JOIN tbl_subscriptions
		ON tbl_subscribers.tier_id = tbl_subscriptions.id
		LEFT JOIN users
		ON tbl_subscriptions.username = users.name
		LEFT JOIN tbl_userdata
		ON users.id = tbl_userdata.userid
	WHERE tbl_subscribers.userid = ? `

	if filter == "0" {
		query = query + ` AND tbl_subscribers.end_date >= NOW()`
	} else if filter == "1" {
		query = query + ` AND tbl_subscribers.end_date < NOW()`
	}

	rows, err := dal.DB.Query(query, userid)
	if err != nil {
		fmt.Printf("Error querying :%v\n", err)
		return nil, err
	}
	defer rows.Close()

	res := make([]models.TierCards, 0)
	for rows.Next() {
		item := models.TierCards{}
		err := rows.Scan(
			&item.Creator.Username,
			&item.Creator.Displayname,
			&item.Creator.Avatar,
			&item.Creator.Banner,
			&item.Tiername,
			&item.Expiration,
			&item.Baseprice,
			&item.Status,
		)
		if err != nil {
			fmt.Printf("Error creating %v\n", err)
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (dal *subscriptionDAL) CreateSubscriber(userid string, tierid string, duration string) error {
	// Find the history of subscribers that are currently subscribed
	query := "SELECT count(*) as count FROM tbl_subscribers where userid = ? AND tier_id = ?"
	var count int
	dal.DB.QueryRow(query, userid, tierid).Scan(&count)
	// Insert the subscriber
	if count == 0 {
		query = `INSERT INTO tbl_subscribers (userid, tier_id, tier_owner, end_date)
		SELECT ?, ?, (SELECT username FROM tbl_subscriptions WHERE id = ?), DATE_ADD(NOW(), INTERVAL ? MONTH)`
		_, err := dal.DB.Exec(query, userid, tierid, tierid, duration)
		if err != nil {
			fmt.Printf("Error inserting subscribers: %v\n", err)
			return err
		}
	} else { // Update the subscriber
		query = `UPDATE
		tbl_subscribers
	  SET
		end_date =
		CASE
		  WHEN end_date < NOW()
		  THEN DATE_ADD(NOW(), INTERVAL ? MONTH)
		  ELSE DATE_ADD(end_date, INTERVAL ? MONTH)
		END
	  WHERE userid = ? AND tier_id = ?`
		_, err := dal.DB.Exec(query, duration, duration, userid, tierid)
		if err != nil {
			fmt.Printf("Error updating subscribers: %v\n", err)
			return err
		}
	}
	return nil
}

func (dal *subscriptionDAL) GetBasePrice(tierid string) (float64, error) {
	query := "SELECT base_price FROM tbl_subscriptions WHERE id = ?"
	var base_price float64
	err := dal.DB.QueryRow(query, tierid).Scan(&base_price)
	if err != nil {
		fmt.Printf("Error Get Base Price : %v\n", err)
		return 0, err
	}
	return base_price, nil
}
