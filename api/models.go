package api

// Account
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR RUP"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type UpdateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required,min=0"`
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Transfer
type CreateTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

type GetTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListTransfersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type DeleteTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Entries
type CreateEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required"`
}

type GetEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListEntriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type DeleteEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Users
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
