package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fmt"
	"strconv"
)

type TiersDAL interface {
	GetTiersByUsername(username string) ([]models.TierModel, error)
	GetOneTierByID(tierID string) (models.EditTierModel, error)
	CreateTier(models.TierModel) error
	UpdateTier(tier models.TierModel) error
	ToggleTier(tierID string) error
}

type tiersDAL struct {
	DB *sql.DB
}

var _ TiersDAL = (*tiersDAL)(nil)

func NewTiersDAL(db *sql.DB) TiersDAL {
	return &tiersDAL{db}
}

func (dal *tiersDAL) ToggleTier(tierID string) error {
	query := "UPDATE tbl_subscriptions SET active = 1-active WHERE id= ?"
	_, err := dal.DB.Exec(query, tierID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}
func (dal *tiersDAL) CreateTier(tier models.TierModel) error {
	fmt.Printf("models: %v\n", tier.Child)

	query := "INSERT INTO tbl_subscriptions (username, tier_name, tier_color, tier_benefit, tier_child, base_price, month_two, month_three, month_six) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	base_price, _ := strconv.ParseFloat(tier.Baseprice, 64)
	month2, _ := strconv.ParseInt(tier.Month2, 10, 64)
	month3, _ := strconv.ParseInt(tier.Month3, 10, 64)
	month6, _ := strconv.ParseInt(tier.Month6, 10, 64)
	_, err := dal.DB.Exec(query, tier.Username, tier.Title, tier.Color, tier.Benefit, tier.Child, base_price, month2, month3, month6)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		return err
	}
	return nil
}

func (dal *tiersDAL) UpdateTier(tier models.TierModel) error {
	query := "UPDATE tbl_subscriptions SET tier_name = ?, tier_color = ?, base_price = ?, tier_child = ?, tier_benefit = ?, month_two = ?, month_three = ?, month_six = ?, active = ? WHERE id = ? AND username = ?"
	base_price, _ := strconv.ParseFloat(tier.Baseprice, 64)
	month2, _ := strconv.ParseInt(tier.Month2, 10, 64)
	month3, _ := strconv.ParseInt(tier.Month3, 10, 64)
	month6, _ := strconv.ParseInt(tier.Month6, 10, 64)
	active, _ := strconv.ParseInt(tier.Active, 10, 64)
	_, err := dal.DB.Exec(query, tier.Title, tier.Color, base_price, tier.Child, tier.Benefit, month2, month3, month6, active, tier.ID, tier.Username)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		return err
	}
	return nil
}

func (dal *tiersDAL) GetOneTierByID(tierID string) (models.EditTierModel, error) {
	query := "SELECT tier_name, tier_color, base_price, tier_child, tier_benefit, month_two, month_three, month_six, active FROM tbl_subscriptions WHERE id = ?"
	returnModel := models.EditTierModel{}
	err := dal.DB.QueryRow(query, tierID).Scan(
		&returnModel.Title,
		&returnModel.Color,
		&returnModel.Baseprice,
		&returnModel.Children,
		&returnModel.Benefits,
		&returnModel.Month_two,
		&returnModel.Month_three,
		&returnModel.Month_six,
		&returnModel.Active,
	)
	if err != nil {
		fmt.Printf("Could not get tier: %v", err)
		return returnModel, err
	}

	// Add RelativeTiers Field to returnModel
	var username string
	query = "SELECT username FROM tbl_subscriptions WHERE id = ?"
	err = dal.DB.QueryRow(query, tierID).Scan(&username)
	if err != nil {
		fmt.Printf("error4: %v\n", err)
		return returnModel, err
	}
	query = "SELECT id, tier_name FROM tbl_subscriptions WHERE username = ? AND id != ?"
	rows, err := dal.DB.Query(query, username, tierID)
	if err != nil {
		fmt.Printf("error5: %v\n", err)
		return returnModel, err
	}
	defer rows.Close()
	relativelist := make([]models.RelativeTier, 0)
	for rows.Next() {
		item := models.RelativeTier{}
		err := rows.Scan(
			&item.ID,
			&item.Title,
		)
		if err != nil {
			fmt.Printf("error6: %v\n", err)
			return returnModel, err
		}
		relativelist = append(relativelist, item)
	}
	returnModel.RelativeTiers = relativelist
	return returnModel, nil
}

func (dal *tiersDAL) GetTiersByUsername(username string) ([]models.TierModel, error) {
	query := "SELECT id, tier_name, tier_color, base_price, month_two, month_three, month_six, active FROM tbl_subscriptions WHERE username = ?"
	rows, err := dal.DB.Query(query, username)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	tierslist := make([]models.TierModel, 0)
	for rows.Next() {
		tier := models.TierModel{}
		err := rows.Scan(
			&tier.ID,
			&tier.Title,
			&tier.Color,
			&tier.Baseprice,
			&tier.Month2,
			&tier.Month3,
			&tier.Month6,
			&tier.Active,
		)
		if err != nil {
			fmt.Printf("error3: %v\n", err)
			return nil, err
		}
		tierslist = append(tierslist, tier)
	}
	return tierslist, nil
}
