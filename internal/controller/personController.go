package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivanjabrony/personApi/internal/model"
	"github.com/ivanjabrony/personApi/internal/model/dto"
	"github.com/ivanjabrony/personApi/internal/service"
)

type PersonCotroller struct {
	personService service.PersonService
}

func NewPersonController(personService service.PersonService) *PersonCotroller {
	return &PersonCotroller{personService: personService}
}

// GetPerson godoc
// @Summary      Get person by ID
// @Description  returning person
// @Tags         person
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of person"
// @Success      200 {object} dto.PersonDto
// @Failure      400 {object} dto.BadResponseDto
// @Router       /persons/{id} [get]
func (pc *PersonCotroller) GetPerson(c *gin.Context) {
	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No iD provided"})
		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ID"})
		return
	}

	person, err := pc.personService.GetPersonById(c.Request.Context(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person info"})
		return
	}

	c.JSON(http.StatusOK, person)
}

// GetPerson godoc
// @Summary      Get all persons with pagination
// @Description  returning persons with pagination
// @Tags         person
// @Accept       json
// @Produce      json
// @Param page query int false "Page number (starting from 1)" default(1)
// @Param page_size query int false "Amount of items on the page" default(10) minimum(1) maximum(100)
// @Success      200 {object} dto.PaginatedPersonsDto
// @Failure      400 {object} dto.BadResponseDto
// @Router       /persons [get]
func (pc *PersonCotroller) GetAllPersons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	persons, err := pc.personService.GetAllPersons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person info"})
		return
	}

	total := len(persons)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	response := dto.PaginatedPersonsDto{
		Data:       persons[offset:min(len(persons), offset+pageSize)],
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	c.JSON(http.StatusOK, response)
}

// GetPerson godoc
// @Summary      Get all persons with filter and pagination
// @Description  returning filtered persons with pagination
// @Tags         person
// @Accept       json
// @Produce      json
// @Param 		 name query string false "Name to match" example("Ivan")
// @Param 		 surname query string false "Surname to match" example("Zabrodin")
// @Param 		 patronymic query string false "Patronymic to match" example("Vladimirovich")
// @Param 	     genders query string false "Collection of genders to match" example("male,female")
// @Param 	     nationalities query string false "Collection of nationalities to match" example("RU,KZ")
// @Param 		 name_like query string false "Name pattern to match" example("Iv**")
// @Param 		 surname_like query string false "Surname patter to match to match" example("Za%")
// @Param 		 patronymic_like query string false "Patronymic pattern to match" example("Vl%")
// @Param 		 age_min query int false "Min wanted age" minimum(0)
// @Param 		 age_max query int false "Max wanted age" maximum(110)
// @Param 		 page query int false "Page number (starting from 1)" default(1)
// @Param 		 page_size query int false "Amount of items on the page" default(10) minimum(1) maximum(100)
// @Success      200 {object} dto.PaginatedPersonsDto
// @Failure      400 {object} dto.BadResponseDto
// @Router       /persons/filtered [get]
func (pc *PersonCotroller) GetFilteredPesons(c *gin.Context) {
	var filter model.PersonFilter

	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}
	if surname := c.Query("surname"); surname != "" {
		filter.Surname = &surname
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		filter.Patronymic = &patronymic
	}
	if nationalities := c.Query("nationalities"); nationalities != "" {
		filter.Nationalities = strings.Split(nationalities, ",")
	}
	if genders := c.Query("genders"); genders != "" {
		filter.Genders = strings.Split(genders, ",")
	}

	if name_like := c.Query("name_like"); name_like != "" {
		filter.NameLike = &name_like
	}
	if surname_like := c.Query("surname_like"); surname_like != "" {
		filter.SurnameLike = &surname_like
	}
	if patronymic_like := c.Query("patronymic_like"); patronymic_like != "" {
		filter.PatronymicLike = &patronymic_like
	}
	if ageMin := c.Query("age_min"); ageMin != "" {
		if val, err := strconv.Atoi(ageMin); err == nil {
			filter.AgeMin = &val
		}
	}
	if ageMax := c.Query("age_max"); ageMax != "" {
		if val, err := strconv.Atoi(ageMax); err == nil {
			filter.AgeMax = &val
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	persons, err := pc.personService.GetPersonsFiltered(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve persons info"})
		return
	}

	total := len(persons)
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	response := dto.PaginatedPersonsDto{
		Data:       persons[offset:min(len(persons), offset+pageSize)],
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	c.JSON(http.StatusOK, response)
}

// CreatePerson godoc
// @Summary     Create person
// @Description Creates new person
// @Tags        person
// @Accept      json
// @Produce     json
// @Param       request body dto.NewPersonDto true "Person data"
// @Success     204 "Creating Success"
// @Failure     400 {object} dto.BadResponseDto
// @Router      /persons [post]
func (pc *PersonCotroller) CreatePerson(c *gin.Context) {
	var createDto dto.NewPersonDto

	err := c.BindJSON(&createDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data for creating"})
		return
	}

	id, err := pc.personService.CreatePerson(c.Request.Context(), &createDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}

	c.JSON(http.StatusOK, id)
}

// UpdatePerson godoc
// @Summary      Update user
// @Description  Updates existing user
// @Tags         person
// @Accept       json
// @Produce      json
// @Param        request body dto.UpdatePersonDto true "Updated data"
// @Success      204 "Update success"
// @Failure      400 {object} dto.BadResponseDto
// @Router       /persons [put]
func (pc *PersonCotroller) UpdatePerson(c *gin.Context) {
	var updateDto dto.UpdatePersonDto

	err := c.BindJSON(&updateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data for updating"})
		return
	}

	err = pc.personService.UpdatePersonById(c.Request.Context(), &updateDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person info"})
		return
	}

	c.JSON(http.StatusOK, updateDto.Id)
}

// DeletePerson godoc
// @Summary      Delete person
// @Description  Deletes person by ID
// @Tags         person
// @Accept       json
// @Produce      json
// @Param        id path int true "Person ID"
// @Success      204 "Delete success"
// @Failure      400 {object} dto.BadResponseDto
// @Router       /persons/{id} [delete]
func (pc *PersonCotroller) DeletePersonById(c *gin.Context) {
	id, exists := c.Params.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No iD provided"})
		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ID"})
		return
	}

	err = pc.personService.DeletePersonById(c.Request.Context(), parsedId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person data"})
		return
	}

	c.JSON(http.StatusOK, id)
}
