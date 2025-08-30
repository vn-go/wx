

package swaggers

// Contact contains contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Info contains metadata about the API.
type Info struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Version     string   `json:"version"`
	Contact     *Contact `json:"contact,omitempty"`
	// License, TermsOfService,... có thể thêm nếu cần
}

// Operation describes a single API operation on a path.
type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"` // Thay thế cho "in: body" parameters
	Responses   map[string]Response   `json:"responses,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
	// Servers, Deprecated, Callbacks,... có thể bổ sung
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name            string             `json:"name"`
	In              string             `json:"in"` // "query", "header", "path", "cookie"
	Description     string             `json:"description,omitempty"`
	Required        bool               `json:"required,omitempty"`
	Schema          *Schema            `json:"schema,omitempty"`
	Example         interface{}        `json:"example,omitempty"`
	Examples        map[string]Example `json:"examples,omitempty"`
	Deprecated      bool               `json:"deprecated,omitempty"`
	AllowEmptyValue bool               `json:"allowEmptyValue,omitempty"`
}

// RequestBody describes a single request body.
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

// MediaType describes the media type object with schema and examples.
type MediaType struct {
	Schema   *Schema             `json:"schema,omitempty"`
	Example  interface{}         `json:"example,omitempty"`
	Examples map[string]Example  `json:"examples,omitempty"`
	Encoding map[string]Encoding `json:"encoding,omitempty"`
}

// Encoding provides information about how a specific property value will be serialized depending on its type.
type Encoding struct {
	ContentType   string            `json:"contentType,omitempty"`   // Content-Type for encoding a specific property. Default is the media type of the overall request body.
	Headers       map[string]Header `json:"headers,omitempty"`       // A map allowing additional information to be provided as headers, for example Content-Disposition.
	Style         string            `json:"style,omitempty"`         // Describes how a specific property value will be serialized. Default depends on the property type.
	Explode       bool              `json:"explode,omitempty"`       // When style is form, specifies whether the parameter value should allow multiple values.
	AllowReserved bool              `json:"allowReserved,omitempty"` // Determines whether reserved characters are allowed unencoded.
}

// Response describes a single response from an API Operation.
type Response struct {
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Links       map[string]Link      `json:"links,omitempty"`
}

// Schema represents a JSON schema object.
type Schema struct {
	Ref                  string             `json:"$ref,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Format               string             `json:"format,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty"` // can be bool or *Schema
	Items                *Schema            `json:"items,omitempty"`
	Description          string             `json:"description,omitempty"`
	Nullable             bool               `json:"nullable,omitempty"`
	Example              interface{}        `json:"example,omitempty"`
	Enum                 []interface{}      `json:"enum,omitempty"`
	// Có thể mở rộng thêm các thuộc tính JSON Schema
}

// Swagger is the root document object for the OpenAPI 3.x spec.
type Swagger struct {
	OpenAPI      string                 `json:"openapi"` // "3.0.3"
	Info         Info                   `json:"info"`
	Servers      []Server               `json:"servers,omitempty"`
	Paths        map[string]PathItem    `json:"paths"`
	Components   Components             `json:"components,omitempty"`
	Security     []map[string][]string  `json:"security,omitempty"`
	Tags         []Tag                  `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Server describes a server.
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Components holds various schemas, responses, parameters, examples, and security schemes.
type Components struct {
	Schemas         map[string]*Schema        `json:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`
	Examples        map[string]Example        `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]Link           `json:"links,omitempty"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty"`
}

// Các struct bổ sung bạn có thể định nghĩa thêm như Header, Link, Example, Encoding, SecurityScheme, Callback, Tag, ExternalDocumentation...

// PathItem mô tả các phương thức HTTP của một đường dẫn.
type PathItem struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
}

// Example object for OpenAPI 3.x
type Example struct {
	Summary       string      `json:"summary,omitempty"`       // Mô tả ngắn gọn ví dụ
	Description   string      `json:"description,omitempty"`   // Mô tả chi tiết ví dụ
	Value         interface{} `json:"value,omitempty"`         // Giá trị ví dụ thực tế (có thể là bất cứ JSON gì)
	ExternalValue string      `json:"externalValue,omitempty"` // URL tham khảo ví dụ bên ngoài (nếu có)
}

// Header describes a header object in OpenAPI 3.x
type Header struct {
	Description     string               `json:"description,omitempty"`     // Header description
	Required        bool                 `json:"required,omitempty"`        // Is header required
	Deprecated      bool                 `json:"deprecated,omitempty"`      // Is header deprecated
	AllowEmptyValue bool                 `json:"allowEmptyValue,omitempty"` // Allows empty value
	Style           string               `json:"style,omitempty"`           // Serialization style (default: simple)
	Explode         bool                 `json:"explode,omitempty"`         // Explode property for serialization
	AllowReserved   bool                 `json:"allowReserved,omitempty"`   // Allow reserved characters
	Schema          *Schema              `json:"schema,omitempty"`          // Defines the type used for the header
	Example         interface{}          `json:"example,omitempty"`         // Example of header value
	Examples        map[string]Example   `json:"examples,omitempty"`        // Multiple examples
	Content         map[string]MediaType `json:"content,omitempty"`         // Content of header with different media types
}

// SecurityScheme describes a security scheme used by the API.
type SecurityScheme struct {
	Type             string      `json:"type"`                       // "apiKey", "http", "oauth2", "openIdConnect"
	Description      string      `json:"description,omitempty"`      // Optional description
	Name             string      `json:"name,omitempty"`             // Required if type is "apiKey" — name of the header, query or cookie parameter
	In               string      `json:"in,omitempty"`               // Required if type is "apiKey" — "query", "header" or "cookie"
	Scheme           string      `json:"scheme,omitempty"`           // Required if type is "http" — e.g. "basic", "bearer"
	BearerFormat     string      `json:"bearerFormat,omitempty"`     // Optional hint to format of bearer token
	Flows            *OAuthFlows `json:"flows,omitempty"`            // Required if type is "oauth2"
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty"` // Required if type is "openIdConnect"
}

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// Link represents a possible design-time link for a response.
// The presence of a link does not guarantee the caller can successfully invoke it,
// rather it provides a known relationship and traversal mechanism between responses and other operations.
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty"` // A relative or absolute reference to an OAS operation.
	OperationId  string                 `json:"operationId,omitempty"`  // The name of an existing, resolvable OAS operation.
	Parameters   map[string]interface{} `json:"parameters,omitempty"`   // A map representing parameters to pass to the linked operation.
	RequestBody  interface{}            `json:"requestBody,omitempty"`  // A literal value or expression to use as a request body when calling the linked operation.
	Description  string                 `json:"description,omitempty"`  // A description of the link.
	Server       *Server                `json:"server,omitempty"`       // A server object to be used by the target operation.
}

// Callback is a map of possible out-of-band callbacks related to the parent operation.
// Each value in the map is a PathItem that describes a request that may be initiated by the API provider
// and the expected responses.
type Callback map[string]PathItem

// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"` // A short description of the target documentation.
	URL         string `json:"url"`                   // The URL for the external documentation. MUST be in the format of a URL.
}

// Tag adds metadata to a single tag that is used by the Operation Object.
// It is not mandatory to have a Tag Object per tag used there.
type Tag struct {
	Name         string                 `json:"name"`                   // The name of the tag.
	Description  string                 `json:"description,omitempty"`  // A short description for the tag.
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"` // Additional external documentation for this tag.
}
