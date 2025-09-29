package handler

import (
	"net/http"
	"server/internal/model"
	"server/internal/service"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	addrSvc *service.AddressService
}

func NewAddressHandler(addrSvc *service.AddressService) *AddressHandler{
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
  r.Zip   = strings.TrimSpace(r.Zip)
  r.City    = strings.TrimSpace(r.City)
  r.Country = strings.TrimSpace(r.Country) // better a dropdown
}

// CreateAddress handles POST /api/v1/users/addr/add
func (h *AddressHandler) CreateAddress(c echo.Context) error {
	// Bind into the request type
	req := new(address)
	bindErr := c.Bind(&req)
	if  bindErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request payload"})
	}
	
	// Normalize & validate
	validateErr := c.Validate(req)
	if validateErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validateErr.Error())
	}
	
	// Extract the user ID from the JWT token
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)                  
  userID := int(claims["user_id"].(float64))
	
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
	if  createErr != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": createErr.Error()})
	}
	
	// Return the newly-created address
	return c.JSON(http.StatusCreated, addr)
}

// GetAddress handles GET api/v1/users/addr/:id
func (h *AddressHandler) GetAddress(c echo.Context) error {
	// parse and validate address id from url params
	addrID, addrIdErr := strconv.Atoi(c.Param("id"))
	if addrIdErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid Address ID"})
	}

	// extract user_id from JWT
  claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
  userID := int(claims["user_id"].(float64))
	
	addr, addrErr := h.addrSvc.GetAddress(userID, addrID)
	if addrErr != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": addrErr.Error()})
	}
	
	if addr.UId != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you do not have access to this address"})
  }
	
	return c.JSON(http.StatusOK, addr)
}


// DeleteAddress handles DELETE api/v1/users/addr/:id
func (h *AddressHandler) DeleteAddress(c echo.Context) error {
	// parse and validate address id from url params
	addrID, addrIdErr := strconv.Atoi(c.Param("id"))
	if addrIdErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid Address ID"})
	}
	
	// extract user_id from JWT
  claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
  userID := int(claims["user_id"].(float64))
	
	addr, addrErr := h.addrSvc.GetAddress(userID, addrID)
	if addrErr != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": addrErr.Error()})
	}
	
	if addr.UId != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you do not have access to this address"})
  }
	
	delErr := h.addrSvc.DeleteAddress(userID, addrID)
	if delErr != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": delErr.Error()})
	}
	
	return c.NoContent(http.StatusNoContent)
}

func (h *AddressHandler) UpdateAddress(c echo.Context) error {
	// parse & validate id from url params
	id, paramErr := strconv.Atoi(c.Param("id"))
	if paramErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid address ID"})
	}

	// bind & normalize & validate body
	req := new(address)
	bindErr := c.Bind(req)

	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request payload"})
	}

	validateErr := c.Validate(req);
	if  validateErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validateErr.Error())
	}

	// extract userID from JWT
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

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
	updateErr := h.addrSvc.UpdateAddress(userID, addr);
	if  updateErr != nil {
		switch updateErr {
		case service.ErrForbidden:
			return c.JSON(http.StatusForbidden, echo.Map{"error": updateErr.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": updateErr.Error()})
		}
	}

	return c.JSON(http.StatusOK, addr)
}