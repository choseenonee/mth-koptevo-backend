package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type routeRepo struct {
	db *sqlx.DB
}

func InitRouteRepo(db *sqlx.DB) Routes {
	return routeRepo{db: db}
}

func (r routeRepo) Create(ctx context.Context, route models.RouteCreate) (int, error) {
	var createdID int

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	routeQuery := `INSERT INTO routes(city_id, price, properties, name) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRowxContext(ctx, routeQuery, route.CityID, route.Price, route.Properties, route.Name).Scan(&createdID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, nil
}

func (r routeRepo) GetByID(ctx context.Context, routeID int) (models.Route, error) {
	var route models.Route

	// TODO: JOIN hotels, обернуть в массив [[u.name, rr.properties], ...], иначе не ебу как это уменьшить (никаки, это пзидец)
	routeQuery := `SELECT r.name, r.price, r.properties, rp.position, c.name, d.name, d.properties, p.properties, array_agg(t.name) AS tags
				   FROM routes r
				   INNER JOIN routes_places rp ON r.id = rp.route_id
				   INNER JOIN places p ON rp.place_id = p.id
				   INNER JOIN district d ON p.district_id = d.id
				   INNER JOIN city c ON p.city_id = c.id
				   LEFT JOIN routes_tags rt ON r.id = rt.route_id
				   LEFT JOIN tags t ON rt.tag_id = t.id
				   LEFT JOIN route_reviews rr ON r.id = rr.route_id
				   INNER JOIN users u ON rr.author_id = u.id
				   WHERE r.id=1
				   GROUP BY r.id, rp.position, c.id, d.id, p.id
				   ORDER BY rp.position;`

	return route, nil
}
