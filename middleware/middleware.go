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
	"github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	KratosSessionKey = "ory_kratos_session"
	TraitsKey        = "kratos_traits"
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
			Last  string `json:"last"`
		} `json:"name"`
		Email string `json:"email"`
	}
	Identity struct {
		ID     string `json:"id"`
		Traits Trait  `json:"traits"`
	}
	KratosSession struct {
		Active                      bool                         `json:"active"`
		AuthenticatedAt             time.Time                    `json:"authenticated_at"`
		AuthenticationMethods       []KratosAuthenticationMethod `json:"authentication_methods"`
		AuthenticatorAssuranceLevel string                       `json:"authenticator_assurance_level"`
		ExpiresAt                   string                       `json:"expires_at"`
		Identity                    Identity                     `json:"identity"`
		ID                          string                       `json:"id"`
	}
)

func (s *KratosSession) ToProtobuf() *todolist.User {
	return &todolist.User{
		Id:        s.ID,
		FirstName: s.Identity.Traits.Name.First,
		LastName:  s.Identity.Traits.Name.Last,
		Email:     s.Identity.Traits.Email,
	}
}
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
func GatewayMetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	// otherwise pass no extra metadata along
	spew.Dump("GatewayMetadataAnnotator")
	user, ok := ctx.Value(TraitsKey).(*todolist.User)
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
		ctx = context.WithValue(ctx, TraitsKey, session.ToProtobuf())

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
	req.AddCookie(cookie)
	resp, err := k.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("wrong status code")
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
