package web

import "company/finance/server"

func Register() {

	e := server.EP

	RoutGroup := e.Group("api/")

	RoutGroup.POST("login", login)
	RoutGroup.POST("adduser", addUser)
	RoutGroup.POST("addaccount", addAccount)
	RoutGroup.POST("addbudget", addBudget)
	RoutGroup.POST("addexpense", addExpense)

	RoutGroup.POST("budgethistory", budgetHistory)
	RoutGroup.POST("expensehistory", expenseHistory)
	RoutGroup.POST("accounthistory", accountHistory)

}
