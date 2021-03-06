/*
Package goa provides the runtime support for goa web services. See also http://goa.design.

package design: https://godoc.org/github.com/raphael/goa/design

package dsl: https://godoc.org/github.com/raphael/goa/design/dsl

Code Generation

goa service development begins with writing the *design* of a service. The design is described using
the goa language implemented by the github.com/raphael/goa/design/dsl package. The goagen tool
consumes the metadata produced from executing the design language to generate service specific code
that glues the underlying HTTP server with action specific code and data structures.

The goa package contains supporting functionality for the generated code including basic request
and response state management through the Context data structure, error handling via the
service and controller ErrorHandler field, middleware support via the Middleware data structure as
well as input (and output) format validation algorithms.

Request Context

The Context data structure provides access to both the request and response state. It implements
the golang.org/x/net/Context interface so that deadlines and cancelation signals may also be
implemented with it.

The request state is accessible through the Get, GetMany and Payload methods which return the values
of the request parameters, query strings and request body. Action specific contexts wrap Context and
expose properly typed fields corresponding to the request parameters and body data structure
descriptions appearing in the design.

The response state can be accessed through the ResponseStatus, ResponseLength and Header methods.
The Context type implements the http.ResponseWriter interface and thus action contexts can be used
in places http.ResponseWriter can. Action contexts provide action specific helper methods that write
the responses as described in the design optionally taking an instance of the media type for
responses that contain a body.

Here is an example showing an "update" action corresponding to following design (extract):

	Resource("bottle", func() {
		DefaultMedia(Bottle)
		Action("update", func() {
			Params(func() {
				Param("bottleID", Integer)
			})
			Payload(UpdateBottlePayload)
			Response(OK)
			Response(NotFound)
		})
	})

The action signature generated by goagen is:

	type BottleController interface {
		goa.Controller
		Update(*UpdateBottleContext) error
	}

where UpdateBottleContext is:

	type UpdateBottleContext struct {
        	*goa.Context
        	BottleID  int
        	Payload   *UpdateBottlePayload
	}

and implements:

	func (ctx *UpdateBottleContext) OK(resp *Bottle) error
	func (ctx *UpdateBottleContext) NotFound() error

The definitions of the Bottle and UpdateBottlePayload data structures are ommitted for brievity.

Controllers

There is one controller interface generated per resource defined via the design language. The
interface exposes the controller actions as well as methods to set controller specific middleware
and error handlers (see below). User code must provide data structures that implement these
interfaces when mounting a controller onto a service. The controller data structure should include
an anonymous field of type *goa.ApplicationController which takes care of implementing the
middleware and error handler handling.

Error Handling

The controller action methods generated by goagen such as the Update method of the BottleController
interface shown above all return an error value. The controller or service-wide error handler (if no
controller specific error handler) function is invoked whenever the value returned by a controller
action is not nil. The handler gets both the request context and the error as argument.

The default handler implementation returns a response with status code 500 containing the error
message in the body. A different error handler can be specificied using the SetErrorHandler
function on either a controller or service wide. goa comes with an alternative error handler - the
TerseErrorHandler - which also returns a response with status 500 but does not write the error
message to the body of the response.

Middleware

A goa middleware is a function that takes and returns a Handler. A Handler is a the low level
function which handles incoming HTTP requests. goagen generates the handlers code so each handler
creates the action specific context and calls the controller action with it.

Middleware can be added to a goa service or a specific controller using the Service type Use method.
goa comes with a few stock middleware that handle common needs such as logging, panic recovery or
using the RequestID header to trace requests across multiple services.

Validation

The goa design language documented in the dsl package makes it possible to attach validations to
data structure definitions. One specific type of validation consists of defining the format that a
data structure string field must follow. Example of formats include email, data time, hostnames etc.
The ValidateFormat function provides the implementation for the format validation invoked from the
code generated by goagen.
*/
package goa
