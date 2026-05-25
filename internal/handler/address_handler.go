package handler

import (
	"net/http"
	"server/internal/model"
	"server/internal/response"
	"server/internal/service"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	addrSvc *service.AddressService
}

func NewAddressHandler(addrSvc *service.AddressService) *AddressHandler {
	return &AddressHandler{addrSvc: addrSvc}
}

// address for sanitation
type address struct {
	Addr1     string `json:"addr_1"    validate:"required,min=5"`
	Addr2     string `json:"addr_2"`
	Zip       string `json:"zip"       validate:"required,min=4"`
	City      string `json:"city"      validate:"required"`
	Country   string `json:"country"   validate:"required"`
	IsDefault bool   `json:"isdefault"`
}

// Normalize implements Normalizable
func (r *address) Normalize() {
	r.Addr1 = strings.TrimSpace(r.Addr1)
	r.Addr2 = strings.TrimSpace(r.Addr2)
	r.Zip = strings.TrimSpace(r.Zip)
	r.City = strings.TrimSpace(r.City)
	r.Country = strings.TrimSpace(r.Country) // better a dropdown
}

// CreateAddress handles POST /api/v1/users/addr/add
func (h *AddressHandler) CreateAddress(c echo.Context) error {
	// Bind into the request type
	req := new(address)
	bindErr := c.Bind(&req)
	if bindErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidPayload, "invalid request payload")
	}

	// Normalize & validate
	validateErr := c.Validate(req)
	if validateErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeValidationError, validateErr.Error())
	}

	// Extract the user ID from the JWT token
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Map the request into model, setting UId from the token
	addr := &model.Address{
		UId:       userID,
		Addr_1:    req.Addr1,
		Addr_2:    req.Addr2,
		Zip:       req.Zip,
		City:      req.City,
		Country:   req.Country,
		IsDefault: req.IsDefault,
	}

	// Call service
	createErr := h.addrSvc.CreateAddress(userID, addr)
	if createErr != nil {
		code, status, msg := MapServiceError(createErr)
		return response.Error(c, status, code, msg)
	}

	// Return the newly-created address
	return c.JSON(http.StatusCreated, addr)
}

// GetAddresses handles GET /api/v1/users/addresses
// Returns all addresses for the authenticated user
func (h *AddressHandler) GetUserAddresses(c echo.Context) error {
	// Extract user ID from JWT
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Get all addresses for user
	addresses, err := h.addrSvc.GetAddresses(userID)
	if err != nil {
		code, status, msg := MapServiceError(err)
		return response.Error(c, status, code, msg)
	}

	// Return empty array if no addresses (not an error)
	if addresses == nil {
		addresses = []model.Address{}
	}

	return c.JSON(http.StatusOK, addresses)
}

// GetAddress handles GET api/v1/users/addr/:id
func (h *AddressHandler) GetAddress(c echo.Context) error {
	// parse and validate address id from url params
	addrID, addrIdErr := strconv.Atoi(c.Param("id"))
	if addrIdErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidAddressID, "invalid address ID")
	}

	// extract user_id from JWT
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	addr, addrErr := h.addrSvc.GetAddress(userID, addrID)
	if addrErr != nil {
		code, status, msg := MapServiceError(addrErr)
		return response.Error(c, status, code, msg)
	}

	return c.JSON(http.StatusOK, addr)
}

// DeleteAddress handles DELETE api/v1/users/addr/:id
func (h *AddressHandler) DeleteAddress(c echo.Context) error {
	// parse and validate address id from url params
	addrID, addrIdErr := strconv.Atoi(c.Param("id"))
	if addrIdErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidAddressID, "invalid address ID")
	}

	// extract user_id from JWT
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	delErr := h.addrSvc.DeleteAddress(userID, addrID)
	if delErr != nil {
		code, status, msg := MapServiceError(delErr)
		return response.Error(c, status, code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AddressHandler) UpdateAddress(c echo.Context) error {
	// parse & validate id from url params
	id, paramErr := strconv.Atoi(c.Param("id"))
	if paramErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidAddressID, "invalid address ID")
	}

	// bind & normalize & validate body
	req := new(address)
	bindErr := c.Bind(req)

	if bindErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidPayload, "invalid request payload")
	}

	validateErr := c.Validate(req)
	if validateErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeValidationError, validateErr.Error())
	}

	// extract userID from JWT
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// map to model.Address
	addr := &model.Address{
		ID:        id,
		UId:       userID,
		Addr_1:    req.Addr1,
		Addr_2:    req.Addr2,
		Zip:       req.Zip,
		City:      req.City,
		Country:   req.Country,
		IsDefault: req.IsDefault,
	}

	// call service
	updateErr := h.addrSvc.UpdateAddress(userID, addr)
	if updateErr != nil {
		code, status, msg := MapServiceError(updateErr)
		return response.Error(c, status, code, msg)
	}

	return c.JSON(http.StatusOK, addr)
}
