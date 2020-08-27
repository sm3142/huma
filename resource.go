package huma

import (
	"net/http"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/go-chi/chi"
)

// Resource represents an API resource attached to a router at a specific path
// (URI template). Resources can have operations or subresources attached to
// them.
type Resource struct {
	path   string
	mux    chi.Router
	router *Router

	subResources []*Resource
	operations   []*Operation

	tags []string
}

func (r *Resource) toOpenAPI() *gabs.Container {
	doc := gabs.New()

	for _, sub := range r.subResources {
		doc.Merge(sub.toOpenAPI())
	}

	for _, op := range r.operations {
		doc.Set(op.toOpenAPI(), r.path, strings.ToLower(op.method))
	}

	return doc
}

// Operation creates a new HTTP operation with the given method at this resource.
func (r *Resource) Operation(method, operationID, description string, responses ...Response) *Operation {
	op := newOperation(r, method, operationID, description, responses)
	r.operations = append(r.operations, op)

	return op
}

// Post creates a new HTTP POST operation at this resource.
func (r *Resource) Post(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodPost, operationID, description, responses...)
}

// Head creates a new HTTP HEAD operation at this resource.
func (r *Resource) Head(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodHead, operationID, description, responses...)
}

// Get creates a new HTTP GET operation at this resource.
func (r *Resource) Get(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodGet, operationID, description, responses...)
}

// Put creates a new HTTP PUT operation at this resource.
func (r *Resource) Put(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodPut, operationID, description, responses...)
}

// Patch creates a new HTTP PATCH operation at this resource.
func (r *Resource) Patch(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodPatch, operationID, description, responses...)
}

// Delete creates a new HTTP DELETE operation at this resource.
func (r *Resource) Delete(operationID, description string, responses ...Response) *Operation {
	return r.Operation(http.MethodDelete, operationID, description, responses...)
}

// AddMiddleware adds a new standard middleware to this resource, so it will
// apply to requests at the resource's path (including any subresources).
// Middleware can also be applied at the router level to apply to all requests.
func (r *Resource) AddMiddleware(middlewares ...func(next http.Handler) http.Handler) *Resource {
	r.mux.Use(middlewares...)
	return r
}

// SubResource creates a new resource attached to this resource. Any passed
// path parts and params are attached to the existing resource path.
func (r *Resource) SubResource(parts ...string) *Resource {
	uriTemplate := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if part[0] == '/' {
			// This is a path component
			uriTemplate += part
		} else {
			// This is a parameter component
			uriTemplate += "/{" + part + "}"
		}
	}

	sub := &Resource{
		path:         r.path + uriTemplate,
		mux:          r.mux.Route(uriTemplate, nil),
		subResources: []*Resource{},
		operations:   []*Operation{},
		tags:         append([]string{}, r.tags...),
	}

	r.subResources = append(r.subResources, sub)

	return sub
}

// AddTags appends to the list of tags, used for documentation.
func (r *Resource) AddTags(names ...string) *Resource {
	r.tags = append(r.tags, names...)
	return r
}
