package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type (
	// AuthFunc is the pluggable function that performs authentication.
	//
	// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
	// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
	//
	// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
	// please make sure that the `Context` returned is a child `Context` of the one passed in.
	//
	// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
	// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
	// (authed, but lacking perms) appropriately.
	AuthFunc func(ctx context.Context) (context.Context, error)
	// KratosMiddleware is a simple authentication middleware
	// for Ory Kratos used by gRPC-gateway
	//
	// The idea of middleware is simple
	//
	// 1. Get ory_kratos_session cookie
	// 2. Check this cookie
	// 3. Get identity and verify it
	// 4. If anything goes wrong redirect to Ory Kratos UI
	KratosMiddleware struct {
		APIURL string
		UIURL  string
		Client *http.Client
	}
	// Identity represents identity sent from Kratos
	Identity struct {
		ID     string `json:"id"`
		Traits struct {
			Name struct {
				First string `json:"first"`
				Last  string `json:"last"`
			} `json:"name"`
			Email string `json:"email"`
		} `json:"traits"`
	}
	// KratosSession represents Kratos session returned
	// from /session/whoami from Ory Kratos
	KratosSession struct {
		Active                bool      `json:"active"`
		AuthenticatedAt       time.Time `json:"authenticated_at"`
		AuthenticationMethods []struct {
			CompletedAt time.Time `json:"completed_at"`
			Method      string    `json:"method"`
		} `json:"authentication_methods"`
		AuthenticatorAssuranceLevel string   `json:"authenticator_assurance_level"`
		ExpiresAt                   string   `json:"expires_at"`
		Identity                    Identity `json:"identity"`
		ID                          string   `json:"id"`
	}
)

// ToProtobuf converts Kratos identity to user
func (s *KratosSession) ToProtobuf() *todolist.User {
	return &todolist.User{
		Id:        s.ID,
		FirstName: s.Identity.Traits.Name.First,
		LastName:  s.Identity.Traits.Name.Last,
		Email:     s.Identity.Traits.Email,
	}
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		newCtx, err = authFunc(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// GatewayResponseModifier modifies response
func GatewayResponseModifier(ctx context.Context, response http.ResponseWriter, _ proto.Message) error {
	return nil
}

// GatewayMetadataAnnotator looks up session and pass userId in to context if it exists
func GatewayMetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	// otherwise pass no extra metadata along
	user, ok := ctx.Value(config.KratosTraitsKey).(*todolist.User)
	if !ok {
		return metadata.Pairs()
	}
	if user != nil {
		md := metadata.Pairs("user_id", user.Id)
		md.Append("first_name", user.FirstName)
		md.Append("last_name", user.LastName)
		return md
	}
	return metadata.Pairs()
}

// Middleware implements middleware
func (k *KratosMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		session, err := k.validateSession(r)
		if err != nil {
			http.Redirect(w, r, k.UIURL, http.StatusFound)
			return
		}
		ctx = context.WithValue(ctx, config.KratosTraitsKey, session.ToProtobuf())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (k *KratosMiddleware) validateSession(r *http.Request) (*KratosSession, error) {
	var session KratosSession
	cookie, err := r.Cookie(config.KratosSessionKey)
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, errors.New("no session in cookies")
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sessions/whoami", k.APIURL), http.NoBody)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)
	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("wrong status code")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}
