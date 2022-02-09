package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	KratosSessionKey = "ory_kratos_session"
)

type (
	AuthFunc         func(ctx context.Context) (context.Context, error)
	KratosMiddleware struct {
		APIURL string
		UIURL  string
		Client *http.Client
	}
	KratosAuthenticationMethod struct {
		CompletedAt time.Time `json:"completed_at"`
		Method      string    `json:"method"`
	}
	Trait struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"Last"`
		} `json:"name"`
		Email string `json:"email"`
	}
	KratosSession struct {
		Active                      bool                         `json:"active"`
		AuthenticatedAt             time.Time                    `json:"authenticated_at"`
		AuthenticationMethods       []KratosAuthenticationMethod `json:"authentication_methods"`
		AuthenticatorAssuranceLevel string                       `json:"authenticator_assurance_level"`
		ExpiresAt                   string                       `json:"expires_at"`
		Traits                      Trait                        `json:"traits"`
		ID                          string                       `json:"id"`
	}
)

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

func GatewayResponseModifier(ctx context.Context, response http.ResponseWriter, _ proto.Message) error {
	spew.Dump("GatewayResponseMidifier")
	return nil
}

// look up session and pass userId in to context if it exists
func GatewayMetadataAnnotator(_ context.Context, r *http.Request) metadata.MD {
	// otherwise pass no extra metadata along
	spew.Dump("GatewayMetadataAnnotator")
	return metadata.Pairs()
}

func (k *KratosMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		session, err := k.validateSession(ctx, r)
		if err != nil {
			http.Redirect(w, r, k.UIURL, http.StatusFound)
			return
		}
		spew.Dump(session)
		ctx = context.WithValue(ctx, 0, r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (k *KratosMiddleware) validateSession(ctx context.Context, r *http.Request) (*KratosSession, error) {
	var session KratosSession
	cookie, err := r.Cookie(KratosSessionKey)
	if cookie == nil {
		return nil, errors.New("no session in cookies")
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sessions/whoami", k.APIURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(KratosSessionKey, cookie.Value)
	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil

}
