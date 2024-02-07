package http

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

func (r *routes) commentCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	postId := r.getIdFromPath(req, 3)
	if postId == "" {
		r.notFound(w)
		return
	}

	form := req.PostForm
	content := form.Get("commentContent")

	comment := &entity.CommentCreateForm{Content: content}

	err := r.service.Comment.SaveComment(comment, postId)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostId):
			r.notFound(w)
			return
		case errors.Is(err, entity.ErrInvalidFormData):
			data := r.newTemplateData(req)
			data.Form = comment
			r.render(w, req, http.StatusUnprocessableEntity, "view.html", data)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	r.sesm.Put(req.Context(), "flash", "Successfully added your comment!")

	// Change this so it redirects either back to post's page or somewhere else
	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postId), http.StatusSeeOther)
}
