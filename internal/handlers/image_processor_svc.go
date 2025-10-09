package httphandlers

import (
	"net/http"
	"strings"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) uploadImage(c *ginext.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp{
			Error:   "Необходимо передать файл изображения",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !validateContentType(contentType) {
		c.JSON(http.StatusBadRequest, errResp{
			Error:   "Неподдерживаемый тип файла",
			Code:    http.StatusBadRequest,
			Details: contentType,
		})
		return
	}

	id, err := h.svc.UploadImage(c.Request.Context(), file, header.Filename)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Ошибка при загрузке изображения: %s", header.Filename)
		c.JSON(http.StatusInternalServerError, errResp{
			Error:   "Ошибка при загрузке изображения",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, uploadResp{
		ID:      id,
		Status:  "uploaded",
		Message: "Изображение успешно загружено",
	})
}

func (h *Handler) getImage(c *ginext.Context) {
	id := c.Param("id")
	imageType := c.Query("type")

	if !validateImgType(imageType) {
		c.JSON(http.StatusBadRequest, errResp{
			Error:   "Неподдерживаемый тип изображения",
			Code:    http.StatusBadRequest,
			Details: imageType,
		})
		return
	}

	path, err := h.svc.GetImage(c.Request.Context(), id, imageType)
	if err != nil {
		zlog.Logger.Error().Err(err).Msgf("Ошибка при получении изображения: %s", id)

		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, errResp{
				Error:   "Изображение не найдено",
				Code:    http.StatusNotFound,
				Details: err.Error(),
			})
			return
		} else if strings.Contains(err.Error(), "еще не обработано") {
			c.JSON(http.StatusNotFound, errResp{
				Error:   "Изображение еще не обработано",
				Code:    http.StatusNotFound,
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errResp{
			Error:   "Ошибка при получении изображения",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}

	c.File(path)
}

func (h *Handler) deleteImage(c *ginext.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteImage(c.Request.Context(), id); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Ошибка при удалении изображения: %s", id)
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, errResp{
				Error:   "Изображение не найдено",
				Code:    http.StatusNotFound,
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errResp{
			Error:   "Ошибка при удалении изображения",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
	}

	c.JSON(http.StatusOK, deleteResp{
		Status:  "deleted",
		Message: "Изображение успешно удалено",
	})
}

func (h *Handler) getStatus(c *ginext.Context) {
	id := c.Param("id")

	meta, err := h.svc.GetImgMeta(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, errResp{
				Error:   "Изображение не найдено",
				Code:    http.StatusNotFound,
				Details: err.Error(),
			})
			return
		} else if strings.Contains(err.Error(), "еще не обработано") {
			c.JSON(http.StatusNotFound, errResp{
				Error:   "Изображение еще не обработано",
				Code:    http.StatusNotFound,
				Details: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errResp{
			Error:   "Ошибка при получении метаданных изображения",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, statusResp{
		ID:      meta.ID,
		Status:  string(meta.Status),
		Message: meta.ErrorMessage,
	})
}
