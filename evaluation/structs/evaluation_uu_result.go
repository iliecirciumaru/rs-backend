package structs

type EvaluationUUResult struct {
	Neighbours      uint       `json:"neighbours"`
	GlobalRMSE      float64    `json:"global_rmse"`
	PredictionsUsed int        `json:"predictions_used"`
	TotalTime       float64    `json:"total_time"`
	TimePerUser     float64    `json:"time_per_user"`
	UsersRMSE       []UserRMSE `json:"-"`
}

type UserRMSE struct {
	UserID int64
	RMSE   float64
}
