package service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
	jwtentity "github.com/daiyanuthsa/grpc-ecom-be/internal/entity/jwt"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/repository"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService *gocache.Cache
}

func (s *authService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Implement your registration logic here
	//Cek Email apakah sudah terdaftar
	user, err := s.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	
	//jika email sudah terdaftar, kembalikan error
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Email is already registered"),
		}, nil
	}
	// cek password dan confirm password
	if request.Password != request.ConfirmPassword {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password and confirm password do not match"),
		}, nil
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Failed to hash password"),
		}, nil
	}

	//simpan data ke database dan kembalikan respon sukses
	newUser := entity.User{
		Id:       uuid.NewString(),
		Email:    request.Email,
		Password: string(hashedPassword),
		FullName: request.FullName,
		RoleCode: entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,

	}
	err = s.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User registered successfully"),
	}, nil
}

func (s *authService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Implement your login logic here
	//Cek Email apakah sudah terdaftar
	user, err := s.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	//jika email tidak terdaftar, kembalikan error
	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("Email is not registered"),
		}, nil
	}

	//cek password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("Invalid email or password"),
		}, nil
	}

	//jika login berhasil, kembalikan respon sukses
	now := time.Now()
	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtentity.JWTClaims{
		FullName: user.FullName,
		Email:    user.Email,
		RoleCode: user.RoleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("Failed to generate access token"),
		}, nil
	}

	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login successful"),
		AccessToken: tokenString,
	}, nil
}

func (s *authService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// Implement your logout logic here (if any)
	// ambil token dari metadata
	jwtToken, err := jwtentity.ParseTokenFromContex(ctx)
	if err != nil {
		return nil, err
	}
	// kembalikan token menjadi entity jwt
	
	tokenClaims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	//masukan token ke memory db .cache
	s.cacheService.Set(jwtToken,"", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout successful"),
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	// cek password dan confirm password
	if request.NewPassword != request.NewConfirmPassword {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("Password and confirm password do not match"),
		}, nil
	}
	// Langkah 1: Ambil identitas pengguna dari token JWT di context
	jwtToken, err := jwtentity.ParseTokenFromContex(ctx)
	if err != nil {
		// Jika token tidak ada atau tidak valid, middleware sudah seharusnya menangani ini,
		// tapi kita kembalikan error untuk keamanan.
		return nil, utils.UnauthenticatedResponse()
	}

	claims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}
	
	log.Println(claims.Email)

	// Langkah 2: Dapatkan data pengguna saat ini dari database
	user, err := s.authRepository.GetUserByEmail(ctx, claims.Email )
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.UnauthenticatedResponse()
	}
		
	// Langkah 3: Verifikasi `OldPassword` dengan hash yang ada di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		// Jika error (tidak cocok), kembalikan respon error yang jelas
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("Incorrect old password"),
		}, nil
	}

	// Langkah 4: Cek apakah password baru dan konfirmasinya cocok
	if request.NewPassword != request.NewConfirmPassword {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("New password and confirmation do not match"),
		}, nil
	}

	// Langkah 5: Hash password baru
	// (Perbaikan bug: menggunakan request.NewPassword, bukan request.Password)
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("Failed to hash password"),
		}, nil
	}

	// Langkah 6: Update password di database
	err = s.authRepository.UpdateUserPassword(ctx, user.Id, string(newHashedPassword), user.FullName)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("Failed to update password"),
		}, nil
	}

	// Langkah 7 (Opsional tapi direkomendasikan):
	// Blacklist token saat ini untuk memaksa pengguna login kembali dengan password baru.
	// Ini adalah praktik keamanan yang baik.
	// remainingDuration := time.Until(claims.ExpiresAt.Time)
	// s.cacheService.Set(jwtToken, "blacklisted", remainingDuration)

	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Password changed successfully"),
	}, nil

}

func (s *authService) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	// Langkah 1: Ambil identitas pengguna dari token JWT di context
	jwtToken, err := jwtentity.ParseTokenFromContex(ctx)
	if err != nil {
		// Jika token tidak ada atau tidak valid, middleware sudah seharusnya menangani ini,
		// tapi kita kembalikan error untuk keamanan.
		return nil, utils.UnauthenticatedResponse()
	}

	claims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	user, err := s.authRepository.GetUserByEmail(ctx, claims.Email )
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.UnauthenticatedResponse()
	}
log.Println(user.CreatedAt)
	return &auth.GetProfileResponse{
		Base: utils.SuccessResponse("Get user profile successful"),
		UserId: user.Id,
		FullName: user.FullName,
		Email: user.Email,
		RoleCode: user.RoleCode,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {

	return &authService{
		authRepository: authRepository,
		cacheService: cacheService,
	}
}