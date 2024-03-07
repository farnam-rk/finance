package acsslog

import "company/finance/acsslog/dailyrotate"

// CreateLogFile open log file (create if not exist)
func CreateLogFile(path string) (*dailyrotate.File, error) {

	file, err := dailyrotate.NewFile(path, func(path string, didRotate bool) { path = "log/"; didRotate = true })
	if err != nil {
		return nil, err
	}
	return file, nil
}
