package http

import (
	"errors"
	"for9may/internal/dto"
	"for9may/internal/service"
	"for9may/pkg/logger"
	"for9may/resources/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type GalleryHandler struct {
	GalleryService *service.GalleryService
}

func NewGalleryHandler(galleryService *service.GalleryService) *GalleryHandler {
	return &GalleryHandler{
		GalleryService: galleryService,
	}
}

// CreatePost
// @Summary create new post
// @Description create new post in gallery
// @Tags Gallery
// @Accept json
// @Produce json
// @Param post body dto.CreateGalleryPostDTO true "post in gallery info"
// @Success 201
// @Failure 422
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /gallery [post]
func (g *GalleryHandler) CreatePost(c *gin.Context) {
	var post dto.CreateGalleryPostDTO
	if err := c.ShouldBindJSON(&post); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	id, err := g.GalleryService.CreatePost(c, &post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

// GetPosts
// @Summary get posts
// @Description get posts
// @Tags Gallery
// @Accept json
// @Produce json
// @Success 200 {object} dto.GalleryPostDTO "posts"
// @Failure 422
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /gallery [get]
func (g *GalleryHandler) GetPosts(c *gin.Context) {
	posts, err := g.GalleryService.GetPosts(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.InternalServerError{})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// DeletePost
// @Summary delete post
// @Description delete post from gallery
// @Tags Gallery
// @Accept json
// @Produce json
// @Param id path string true "post id" format(uuid)
// @Success 204
// @Failure 422
// @Failure 500 {object} web.InternalServerError "Internal server error"
// @Router /gallery/{id} [delete]
func (g *GalleryHandler) DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "id mast be uuid v4")
		return
	}

	if err := g.GalleryService.DeletePost(c, postID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UploadPostFile
// @Summary upload file
// @Description Upload file, use only .jpg and .png
// @Tags Gallery
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file (jpeg/png)"
// @Param id path string true "post id" format(uuid)
// @Success 201
// @Failure 413
// @Failure 500
// @Router /gallery/file/upload/{id} [post]
func (g *GalleryHandler) UploadPostFile(c *gin.Context) {
	localLogger := logger.GetLoggerFromCtx(c)

	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "id mast be uuid v4")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "fail is not load"})
			return
		}

		localLogger.Error(c, "load file error", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.InternalServerError{})
		return
	}
	if file.Size == 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.ValidationError{Message: "file is empty"})
		return
	}
	link, err := g.GalleryService.UploadPostFile(c, file, postID)

	c.JSON(http.StatusCreated, gin.H{"link": link})
}
